package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	goNanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nedpals/valentine-wall/postal_office/types"
	"golang.org/x/sync/semaphore"
)

var mailgunApiKey string
var mailgunDomain string
var maxJobs = 1250
var port = 3350
var targetEnv = "development"

func init() {
	if gotMailgunDomain, exists := os.LookupEnv("MAILGUN_DOMAIN"); exists {
		mailgunDomain = gotMailgunDomain
	}

	if gotMailgunApiKey, exists := os.LookupEnv("MAILGUN_API_KEY"); exists {
		mailgunApiKey = gotMailgunApiKey
	}

	if gotMaxJobs, exists := os.LookupEnv("MAX_JOBS"); exists {
		if convertedMaxJobs, err := strconv.Atoi(gotMaxJobs); err != nil {
			log.Fatalln(err)
		} else {
			maxJobs = convertedMaxJobs
		}
	}

	if gotPort, exists := os.LookupEnv("PORT"); exists {
		if convertPort, err := strconv.Atoi(gotPort); err != nil {
			log.Fatalln(err)
		} else {
			port = convertPort
		}
	}

	if gotEnv, exists := os.LookupEnv("ENV"); exists {
		targetEnv = gotEnv
	}
}
}

type PostalOffice struct {
	sema              *semaphore.Weighted
	ctx               context.Context
	cancelCh          chan string
	mg                *mailgun.MailgunImpl
	mappedPendingJobs *sync.Map
}

func (po *PostalOffice) NewJob(args *types.NewJobArgs, jobId *string) error {
	if len(args.UniqueID) == 0 {
		return fmt.Errorf("provide a unique id")
	}

	if err := po.sema.Acquire(po.ctx, 1); err != nil {
		return err
	}

	gotJobId, err := goNanoid.New()
	if err != nil {
		return err
	}

	go func(mg *mailgun.MailgunImpl, msg types.MailMessage, sendDuration time.Duration, assignedJobId string, assignedUniqueId string) {
		for {
			select {
			case jobIdToCancel := <-po.cancelCh:
				if jobIdToCancel == assignedJobId {
					log.Printf("[%s] job cancelled", assignedJobId)
					return
				}
			case <-time.After(sendDuration):
				mgMsg := mg.NewMessage(
					fmt.Sprintf("%s <mailgun@%s>", msg.Name, mailgunDomain),
					msg.Subject,
					msg.Content,
					msg.ToEmail,
				)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				defer po.sema.Release(1)
				defer po.mappedPendingJobs.Delete(assignedUniqueId)

				if _, _, err := mg.Send(ctx, mgMsg); err != nil {
					log.Printf("[%s] job error: %s", assignedJobId, err.Error())
				} else {
					log.Printf("[%s] job success", assignedJobId)
				}

				return
			}
		}
	}(po.mg, args.Payload, args.After, gotJobId, args.UniqueID)

	*jobId = gotJobId
	po.mappedPendingJobs.Store(args.UniqueID, gotJobId)
	log.Printf("new job received: %s\n", gotJobId)
	return nil
}

func (po *PostalOffice) CancelJob(args *types.CancelJobArgs, ok *bool) error {
	po.cancelCh <- args.JobID
	*ok = true
	log.Printf("[%s] cancelling job\n", args.JobID)
	return nil
}

func (po *PostalOffice) GetJobID(args *types.GetJobIDArgs, jobId *string) error {
	gotJobId, exists := po.mappedPendingJobs.Load(args.UniqueID)
	if !exists {
		return fmt.Errorf("no job id for specific unique id found")
	}
	*jobId = gotJobId.(string)
	return nil
}

func main() {
	postalOffice := &PostalOffice{
		sema:              semaphore.NewWeighted(int64(maxJobs)),
		ctx:               context.Background(),
		cancelCh:          make(chan string, maxJobs*2),
		mg:                mailgun.NewMailgun(mailgunDomain, mailgunApiKey),
		mappedPendingJobs: &sync.Map{},
	}

	rpc.Register(postalOffice)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		log.Fatal("listen error:", e)
	}

	http.Serve(l, nil)
}
