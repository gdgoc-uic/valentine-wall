package main

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/nedpals/valentine-wall/postal_office/types"
)

type EmailSender interface {
	Message(toRecipientEmail string) types.MailMessage
	SendAfter() time.Duration
}

func newEmailSendJob(cl *rpc.Client, ms EmailSender, toRecipient string, uid string) (string, error) {
	if cl == nil {
		return "", fmt.Errorf("postal client is disconnected")
	}

	var receivedJobId string
	mailJobArgs := &types.NewJobArgs{
		Type:     types.JobSend,
		After:    ms.SendAfter(),
		Payload:  ms.Message(toRecipient),
		UniqueID: uid,
	}
	if err := cl.Call("PostalOffice.NewJob", mailJobArgs, &receivedJobId); err != nil {
		return "", err
	}
	return receivedJobId, nil
}
