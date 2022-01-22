package types

import "time"

const DefaultEmailSendExp = 5 * time.Minute

type JobType int

const (
	JobSend JobType = 1
)

type NewJobArgs struct {
	Type     JobType
	UniqueID string
	After    time.Duration
	Payload  MailMessage
}

type MailMessage struct {
	Name    string
	ToEmail string
	Subject string
	Content string
}

type CancelJobArgs struct {
	JobID string
}

type GetJobIDArgs struct {
	UniqueID string
}
