package client

import (
	"fmt"
	"log"
	"net/rpc"
	"reflect"
	"sync"

	"github.com/nedpals/valentine-wall/postal_office/types"
)

type Client struct {
	sync.Mutex
	rpcClient *rpc.Client
	address   string
}

func (cl *Client) Call(proc string, args interface{}, reply interface{}) error {
	if err := cl.checkClient(); err != nil {
		return err
	}

	return cl.rpcClient.Call(proc, args, reply)
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

func (cl *Client) BeginReconnect(rpcError error) {
	newRpcError := rpcError
	for newRpcError == nil || cl.rpcClient == nil {
		if cl.rpcClient == nil {
			log.Println("Restarting RPC Connection due to nil connection")
			if cl.rpcClient, newRpcError = rpc.DialHTTP("tcp", cl.address); newRpcError != nil {
				log.Println(newRpcError)
			}
		} else if rpcError == rpc.ErrShutdown || reflect.TypeOf(rpcError) == reflect.TypeOf((*rpc.ServerError)(nil)).Elem() {
			log.Println("Restarting RPC Connection due to error")
			cl.Lock()
			if cl.rpcClient, newRpcError = rpc.DialHTTP("tcp", cl.address); newRpcError != nil {
				log.Println(newRpcError)
			}
			cl.Unlock()
		} else {
			log.Println(newRpcError)
			return
		}
	}
}

func DialHTTP(address string) (*Client, error) {
	var err error
	connection := &Client{address: address}
	connection.rpcClient, err = rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
