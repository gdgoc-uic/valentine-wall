package main

import (
	"bytes"
	"fmt"
	"net/rpc"
	"text/template"
	"time"

	"github.com/nedpals/valentine-wall/postal_office/types"
)

type EmailSender interface {
	Message(toRecipientEmail string) (*types.MailMessage, error)
	TimeToSend() time.Duration
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
		Type:       types.JobSend,
		TimeToSend: ms.TimeToSend(),
		Payload:    gotMsgPayload,
		UniqueID:   uid,
	}
	if err := cl.Call("PostalOffice.NewJob", mailJobArgs, &receivedJobId); err != nil {
		return "", err
	}
	return receivedJobId, nil
}

type TemplatedMailSender struct {
	subjectTemplate            *template.Template
	hasSubjectTemplateCompiled bool
	template                   *template.Template
	emailName                  string
	subject                    string
	timeToSend                 time.Duration
	data                       interface{}
}

func newTemplatedMailSender(tmpl *template.Template, emailName, subject string, timeToSend time.Duration) *TemplatedMailSender {
	return &TemplatedMailSender{
		template:   tmpl,
		emailName:  emailName,
		subject:    subject,
		timeToSend: timeToSend,
	}
}

func (t *TemplatedMailSender) Message(toRecipientEmail string) (*types.MailMessage, error) {
	if !t.hasSubjectTemplateCompiled {
		var err error
		t.subjectTemplate, err = template.New("subject").Parse(t.subject)
		if err != nil {
			return nil, err
		}
		t.hasSubjectTemplateCompiled = true
	}

	subjectBuf := &bytes.Buffer{}
	contentBuf := &bytes.Buffer{}
	if err := t.subjectTemplate.Execute(subjectBuf, t.data); err != nil {
		return nil, err
	} else if err := t.template.Execute(contentBuf, t.data); err != nil {
		return nil, err
	}

	return &types.MailMessage{
		Name:    t.emailName,
		ToEmail: toRecipientEmail,
		Subject: subjectBuf.String(),
		Content: contentBuf.String(),
	}, nil
}

func (t *TemplatedMailSender) TimeToSend() time.Duration {
	return t.timeToSend
}

func (t *TemplatedMailSender) With(data interface{}) *TemplatedMailSender {
	return &TemplatedMailSender{
		template:                   t.template,
		emailName:                  t.emailName,
		subject:                    t.subject,
		timeToSend:                 t.timeToSend,
		subjectTemplate:            t.subjectTemplate,
		hasSubjectTemplateCompiled: t.hasSubjectTemplateCompiled,
		data:                       data,
	}
}
