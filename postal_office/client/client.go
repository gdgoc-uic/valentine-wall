package client

import (
	"fmt"
	"log"
	"net/rpc"
	"reflect"
	"sync"

	"github.com/nedpals/valentine-wall/postal_office/types"
	"github.com/oleiade/lane"
)

type PendingJob struct {
	job  string
	args interface{}
}

type Client struct {
	sync.Mutex
	rpcClient *rpc.Client
	address   string
	// pendingJobs is used when the server is down.
	// all the pending jobs are then cleared when the server
	// has already been reconnected.
	pendingJobs *lane.Queue
}

func (cl *Client) Call(proc string, args interface{}, reply interface{}) error {
	if err := cl.checkClient(); err != nil {
		return err
	} else if err := cl.rpcClient.Call(proc, args, reply); err != nil {
		go cl.BeginReconnect(err)
		cl.pendingJobs.Enqueue(&PendingJob{proc, args})
		return err
	}
	return nil
}

func (cl *Client) checkClient() error {
	if cl == nil {
		return fmt.Errorf("postal client is disconnected")
	}
	return nil
}

func (cl *Client) NewJob(args *types.NewJobArgs) (string, error) {
	var receivedJobId string
	if err := cl.Call("PostalOffice.NewJob", args, &receivedJobId); err != nil {
		return "", err
	}
	return receivedJobId, nil
}

func (cl *Client) GetJobID(uid string) (string, error) {
	var receivedJobId string
	if err := cl.Call("PostalOffice.NewJob", &types.GetJobIDArgs{UniqueID: uid}, &receivedJobId); err != nil {
		return "", err
	}
	return receivedJobId, nil
}

func (cl *Client) CancelJob(jobId string) error {
	var ok bool
	if err := cl.Call("PostalOffice.NewJob", &types.CancelJobArgs{JobID: jobId}, &ok); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("job not cancelled")
	}
	return nil
}

func (cl *Client) CancelJobByUID(uid string) error {
	receivedJobId, err := cl.GetJobID(uid)
	if err != nil {
		return err
	}
	return cl.CancelJob(receivedJobId)
}

func (cl *Client) resendJobs() {
	wg := sync.WaitGroup{}
	success := 0
	for cl.pendingJobs.Head() != nil {
		j := cl.pendingJobs.Dequeue().(*PendingJob)
		wg.Add(1)
		go func(j *PendingJob) {
			if err := cl.Call(j.job, j.args, nil); err != nil {
				log.Printf("job error: %s\n", err)
			} else {
				success++
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
	log.Printf("%d jobs have been successfully sent\n", success)
}

func (cl *Client) reconnect() error {
	cl.Lock()
	defer cl.Unlock()
	var newRpcError error
	if cl.rpcClient, newRpcError = rpc.DialHTTP("tcp", cl.address); newRpcError != nil {
		return newRpcError
	}
	cl.resendJobs()
	return nil
}

func (cl *Client) BeginReconnect(rpcError error) {
	newRpcError := rpcError
	for newRpcError == nil || cl.rpcClient == nil {
		if cl.rpcClient == nil {
			log.Println("Restarting RPC Connection due to nil connection")
			if newRpcError = cl.reconnect(); newRpcError != nil {
				log.Println(newRpcError)
			}
		} else if rpcError == rpc.ErrShutdown || reflect.TypeOf(rpcError) == reflect.TypeOf((*rpc.ServerError)(nil)).Elem() {
			log.Println("Restarting RPC Connection due to error")
			if newRpcError = cl.reconnect(); newRpcError != nil {
				log.Println(newRpcError)
			}
		} else {
			log.Println(newRpcError)
			return
		}
	}
}

func DialHTTP(address string) (*Client, error) {
	var err error
	connection := &Client{address: address, pendingJobs: lane.NewQueue()}
	connection.rpcClient, err = rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
