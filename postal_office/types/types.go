package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

const DefaultEmailSendExp = 5 * time.Minute

type JobType int

const (
	JobSend JobType = 1
)

type NewJobArgs struct {
	ID         string
	Type       JobType
	UniqueID   string
	TimeToSend time.Duration
	Payload    *MailMessage
}

type MailMessage struct {
	Name    string `json:"name"`
	ToEmail string `json:"to_email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type CancelJobArgs struct {
	JobID string
}

type GetJobIDArgs struct {
	UniqueID string
}

type JobPayload struct {
	ParentJob  *PendingJob   `json:"-"`
	Type       JobType       `json:"type"`
	TimeToSend time.Duration `json:"time_to_send"`
	Message    *MailMessage  `json:"message"`
}

func (jp *JobPayload) Scan(value interface{}) error {
	strJob, isString := value.(string)
	if value == nil || !isString {
		return nil
	}

	sr := strings.NewReader(strJob)
	decodedJp := JobPayload{}
	if err := json.NewDecoder(sr).Decode(&decodedJp); err != nil {
		return err
	}

	*jp = decodedJp
	return nil
}

func (jp *JobPayload) Value() (driver.Value, error) {
	encodedJp, err := json.Marshal(jp)
	if err != nil {
		return nil, err
	}
	return string(encodedJp), nil
}

type PendingJob struct {
	ID        string      `db:"id"`
	UniqueID  string      `db:"unique_id"`
	Payload   *JobPayload `db:"payload"`
	UpdatedAt time.Time   `db:"updated_at"`
}
