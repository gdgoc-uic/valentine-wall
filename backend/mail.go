package main

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/nedpals/valentine-wall/postal_office/types"
)

type EmailSender interface {
	Message(toRecipientEmail string) (*types.MailMessage, error)
	SendAfter() time.Duration
}

func newEmailSendJob(cl *rpc.Client, ms EmailSender, toRecipient string, uid string) (string, error) {
	if cl == nil {
		return "", fmt.Errorf("postal client is disconnected")
	}

	gotMsgPayload, err := ms.Message(toRecipient)
	if err != nil {
		return "", err
	}

	var receivedJobId string
	mailJobArgs := &types.NewJobArgs{
		Type:     types.JobSend,
		After:    ms.SendAfter(),
		Payload:  gotMsgPayload,
		UniqueID: uid,
	}
	if err := cl.Call("PostalOffice.NewJob", mailJobArgs, &receivedJobId); err != nil {
		return "", err
	}
	return receivedJobId, nil
}

type EmailSenderFunc func(string) (string, error)

func (esf EmailSenderFunc) Message(toRecipientEmail string) (*types.MailMessage, error) {
	gotContent, err := esf(toRecipientEmail)
	if err != nil {
		return nil, err
	}

	return &types.MailMessage{
		Name:    "Mr. Kupido",
		ToEmail: toRecipientEmail,
		Subject: "Valentine Wall - Welcome!",
		Content: gotContent,
	}, nil
}

func (esf EmailSenderFunc) SendAfter() time.Duration {
	return 20 * time.Second
}
