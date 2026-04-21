package domain

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// JobStatus represents the possible statuses for a Job.
type JobStatus string

const (
	JobStatusPending    JobStatus = "PENDING"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusFailed     JobStatus = "FAILED"
)

// allowedJobStatuses holds every valid status value.
var allowedJobStatuses = []JobStatus{
	JobStatusPending,
	JobStatusProcessing,
	JobStatusCompleted,
	JobStatusFailed,
}

type Job struct {
	ID               string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"notnull"`
	VideoID          string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Video            *Video    `json:"video" valid:"-"`
	Error            string    `json:"error" valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func isValidJobStatus(status string) bool {
	for _, s := range allowedJobStatuses {
		if string(s) == status {
			return true
		}
	}
	return false
}

func NewJob(output string, status JobStatus, video *Video) (*Job, error) {
	job := Job{
		OutputBucketPath: output,
		Status:           string(status),
		Video:            video,
		// VideoID:          video.ID,
	}
	job.prepare()
	err := job.Validate()
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (job *Job) prepare() {
	job.ID = uuid.NewV4().String()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	job.Status = string(JobStatusPending)
}

func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)
	if err != nil {
		return err
	}

	if !isValidJobStatus(job.Status) {
		return errors.New("invalid job status: " + job.Status)
	}

	return nil
}
