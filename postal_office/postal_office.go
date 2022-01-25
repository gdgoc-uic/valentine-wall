package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mailgun/mailgun-go/v4"
	goNanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"
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

func dataDir() string {
	baseDir := "."
	if execPath, err := os.Executable(); err != nil {
		log.Println(err)
	} else {
		baseDir = filepath.Dir(execPath)
	}

	return filepath.Join(baseDir, "_data")
}

func databasePath() string {
	return filepath.Join(dataDir(), fmt.Sprintf("postal_office-%s.db", targetEnv))
}

func needOpenDB() bool {
	if _, err := os.Stat(dataDir()); os.IsNotExist(err) {
		return false
	} else if _, err := os.Stat(databasePath()); os.IsNotExist(err) {
		return false
	}

	return true
}

func loadDatabase() (*sqlx.DB, error) {
	if err := os.MkdirAll(dataDir(), os.ModePerm); err != nil {
		return nil, err
	}

	databasePath := databasePath()
	log.Printf("loading database from %s\n", databasePath)
	db, err := sqlx.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}

	schema := `
CREATE TABLE IF NOT EXISTS pending_jobs (
id TEXT PRIMARY KEY,
unique_id TEXT NOT NULL,
payload TEXT NOT NULL,
updated_at TIMESTAMP NOT NULL
)`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}

type PostalOffice struct {
	sema              *semaphore.Weighted
	ctx               context.Context
	cancelCh          chan string
	mg                *mailgun.MailgunImpl
	mappedPendingJobs *sync.Map
	pendingJobsCount  int
}

func deleteJob(po *PostalOffice, uniqueId string) {
	log.Printf("deleting job with unique id %s\n", uniqueId)
	po.mappedPendingJobs.Delete(uniqueId)
	po.sema.Release(1)
	po.pendingJobsCount--
}

func (po *PostalOffice) NewJob(args *types.NewJobArgs, jobId *string) error {
	if len(args.UniqueID) == 0 {
		return fmt.Errorf("provide a unique id")
	}

	if err := po.sema.Acquire(po.ctx, 1); err != nil {
		return err
	}

	gotJobId := args.ID
	if len(gotJobId) == 0 {
		var err error
		gotJobId, err = goNanoid.New()
		if err != nil {
			return err
		}
	}

	timeNow := time.Now()
	go func(po *PostalOffice, msg *types.MailMessage, sendDuration time.Duration, assignedJobId string, assignedUniqueId string) {
		defer deleteJob(po, assignedUniqueId)

		for {
			select {
			case jobIdToCancel := <-po.cancelCh:
				if jobIdToCancel == assignedJobId || jobIdToCancel == "*" {
					log.Printf("[%s] job cancelled\n", assignedJobId)
					return
				}
			case <-time.After(sendDuration):
				mgMsg := po.mg.NewMessage(
					fmt.Sprintf("%s <mailgun@%s>", msg.Name, mailgunDomain),
					msg.Subject,
					msg.Content,
					msg.ToEmail,
				)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				if _, _, err := po.mg.Send(ctx, mgMsg); err != nil {
					log.Printf("[%s] job error: %s\n", assignedJobId, err.Error())
				} else {
					log.Printf("[%s] job success\n", assignedJobId)
				}

				return
			}
		}
	}(po, args.Payload, args.After, gotJobId, args.UniqueID)

	*jobId = gotJobId
	pendingJob := &types.PendingJob{
		ID:       gotJobId,
		UniqueID: args.UniqueID,
		Payload: &types.JobPayload{
			Type:      args.Type,
			SendAfter: args.After,
			Message:   args.Payload,
		},
		UpdatedAt: timeNow,
	}

	pendingJob.Payload.ParentJob = pendingJob
	po.mappedPendingJobs.Store(args.UniqueID, pendingJob)
	po.pendingJobsCount++

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
	rawPendingJob, exists := po.mappedPendingJobs.Load(args.UniqueID)
	if !exists {
		return fmt.Errorf("no job id for specific unique id found")
	}
	*jobId = rawPendingJob.(*types.PendingJob).ID
	return nil
}

func loadPendingJobs(po *PostalOffice) error {
	if !needOpenDB() {
		log.Println("skipping loading pending jobs")
		return nil
	}

	db, err := loadDatabase()
	if err != nil {
		return err
	}

	defer db.Close()
	pendingJobs := []*types.PendingJob{}
	if err := db.Select(&pendingJobs, "SELECT * FROM pending_jobs"); err != nil {
		return err
	}

	success := 0
	for _, pj := range pendingJobs {
		var jobId string
		err := po.NewJob(&types.NewJobArgs{
			Type:     pj.Payload.Type,
			UniqueID: pj.UniqueID,
			ID:       pj.ID,
			After:    pj.Payload.RemainingSendAfter,
			Payload:  pj.Payload.Message,
		}, &jobId)
		if err != nil {
			log.Println(err)
		} else {
			success++
		}
	}

	db.MustExec("DELETE FROM pending_jobs")
	log.Printf("%d jobs were loaded successfully\n", success)
	return nil
}

func savePendingJobs(po *PostalOffice) error {
	log.Println("saving pending jobs...")
	if po.pendingJobsCount == 0 {
		return nil
	}

	db, err := loadDatabase()
	if err != nil {
		return err
	}

	defer db.Close()
	tx := db.MustBegin()
	po.mappedPendingJobs.Range(func(key, value interface{}) bool {
		pj := value.(*types.PendingJob)
		timeDiff := time.Now().Sub(pj.UpdatedAt)
		pj.Payload.RemainingSendAfter = pj.Payload.SendAfter - timeDiff
		if _, err := tx.NamedExec("INSERT INTO pending_jobs (id, unique_id, payload, updated_at) VALUES (:id, :unique_id, :payload, :updated_at)", pj); err != nil {
			log.Println(err)
			return false
		}
		return true
	})
	return tx.Commit()
}

func gracefulShutdown(po *PostalOffice) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		// clean up here
		if err := savePendingJobs(po); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Shutting down gracefully.")
		os.Exit(0)
	}()
}

func main() {
	postalOffice := &PostalOffice{
		sema:              semaphore.NewWeighted(int64(maxJobs)),
		ctx:               context.Background(),
		cancelCh:          make(chan string, maxJobs*2),
		mg:                mailgun.NewMailgun(mailgunDomain, mailgunApiKey),
		mappedPendingJobs: &sync.Map{},
	}

	if err := loadPendingJobs(postalOffice); err != nil {
		log.Println(err)
	}

	rpc.Register(postalOffice)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		log.Fatal("listen error:", e)
	}

	go http.Serve(l, nil)
	go gracefulShutdown(postalOffice)
	forever := make(chan int)
	<-forever
}
