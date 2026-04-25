package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//golang não é necessário dizer qual interface tua classe implementa

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDb struct {
	Db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDb {
	return &VideoRepositoryDb{Db: db}
}

func (repo VideoRepositoryDb) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}
	err := repo.Db.Create(video).Error

	if err != nil {
		return nil, err
	}
	return video, err
}

func (repo VideoRepositoryDb) Find(id string) (*domain.Video, error) {
	var video domain.Video
	repo.Db.First(&video, "id = ?", id) // fill the video found to the video variable for reference
	if video.ID == "" {
		return nil, fmt.Errorf("Video does not exist")
	}

	return &video, nil

}
