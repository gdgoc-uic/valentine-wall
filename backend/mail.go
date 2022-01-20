package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var messageSendWindowDuration = 10 * time.Minute
var pendingEmailMessages = &sync.Map{}

func addEmailSendEntry(typ string, mg *mailgun.MailgunImpl, es EmailSender, recipientEmail string) {
	id := fmt.Sprintf("send_%s", es.PendingMessageID())
	pendingEmailMessages.Store(id, time.AfterFunc(messageSendWindowDuration, func() {
		defer pendingEmailMessages.Delete(id)
		if _, _, err := sendEmail(mg, es, recipientEmail); err != nil {
			log.Println(err)
		}
	}))
}

type EmailSender interface {
	Message(mg *mailgun.MailgunImpl, toRecipientEmail string) *mailgun.Message
	PendingMessageID() string
}

func sendEmail(mg *mailgun.MailgunImpl, ms EmailSender, toRecipient string) (string, string, error) {
	msg := ms.Message(mg, toRecipient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mg.Send(ctx, msg)
}
