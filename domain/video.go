package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	ResourceID string    `json:"resource_id" valid:"uuid" gorm:"type:uuid;notnull"`
	FilePath   string    `json:"file_path" valid:"notnull" gorm:"column:file_path;type:varchar(255)"`
	CreatedAt  time.Time `json:"-" valid:"-" gorm:"auto_create_time"`
	Jobs       []Job     `json:"-" valid:"-" gorm:"ForeignKey:VideoID"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// toda vez que for criar um video, irá ser alterado a referencia
func NewVideo() *Video {
	return &Video{}
}

// independente do método executado, será alterado na struct principal
func (video *Video) Validate() error {
	_, err := govalidator.ValidateStruct(video)
	return err
}
