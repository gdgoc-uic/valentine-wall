package main

import (
	"bytes"
	"html/template"
	"log"
	"net/mail"

	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type TemplatedMailSender struct {
	subjectTemplate            *template.Template
	hasSubjectTemplateCompiled bool
	template                   *template.Template
	emailName                  string
	subject                    string
	data                       any
}

func (t *TemplatedMailSender) Message(meta settings.MetaConfig, toRecipientEmail string) (*mailer.Message, error) {
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

	return &mailer.Message{
		From: mail.Address{
			Address: meta.SenderAddress,
			Name:    meta.SenderName,
		},
		To:      mail.Address{Address: toRecipientEmail},
		Subject: subjectBuf.String(),
		HTML:    contentBuf.String(),
	}, nil
}

func (t *TemplatedMailSender) With(data interface{}) *TemplatedMailSender {
	return &TemplatedMailSender{
		template:                   t.template,
		emailName:                  t.emailName,
		subject:                    t.subject,
		subjectTemplate:            t.subjectTemplate,
		hasSubjectTemplateCompiled: t.hasSubjectTemplateCompiled,
		data:                       data,
	}
}

func newTemplatedMailSender(tmpl *template.Template, emailName, subject string) *TemplatedMailSender {
	return &TemplatedMailSender{
		template:  tmpl,
		emailName: emailName,
		subject:   subject,
	}
}

type emailTemplatesList struct {
	reply   *TemplatedMailSender
	message *TemplatedMailSender
	welcome *TemplatedMailSender
}

var emailTemplates emailTemplatesList

func init() {
	log.Println("loading email templates...")
	rawEmailTemplates := template.Must(template.ParseGlob("./templates/mail/*.txt.tpl"))
	log.Printf("%d email templates have been loaded\n", len(rawEmailTemplates.Templates()))
	emailTemplates = emailTemplatesList{
		reply:   newTemplatedMailSender(rawEmailTemplates.Lookup("reply.txt.tpl"), "Mr. Kupido", "Your message has received a reply!"),
		message: newTemplatedMailSender(rawEmailTemplates.Lookup("message.txt.tpl"), "Mr. Kupido", "You received a new message!"),
		welcome: newTemplatedMailSender(rawEmailTemplates.Lookup("welcome.txt.tpl"), "UIC Valentine Wall", "Welcome to UIC Valentine Wall 2023!"),
	}
}
