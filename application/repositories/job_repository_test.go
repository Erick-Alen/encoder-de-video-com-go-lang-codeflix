package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRepositoryInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)
	job, err := domain.NewJob("output_path", domain.JobStatusPending, video)
	require.Nil(t, err)

	jobRepo := repositories.JobRepositoryDb{Db: db}
	jobRepo.Insert(job)

	j, err := jobRepo.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, job.ID, j.ID)
	// require.Equal(t, video.ID, j.Video.ID)
	// require.Equal(t, "output/raw", j.Video.FilePath)
	// require.Equal(t, domain.JobStatusPending, j.Status)
}

// func TestJobRepositoryFind(t *testing.T) {
// 	db := database.NewDbTest()
// 	defer db.Close()

// 	video := domain.NewVideo()
// 	video.ID = uuid.NewV4().String()
// 	video.FilePath = "path/to/video"
// 	video.CreatedAt = time.Now()

// 	repo := repositories.VideoRepositoryDb{Db: db}
// 	repo.Insert(video)
// 	job, _:= domain.NewJob("output/raw", domain.JobStatusPending, video)
// 	repo.Insert(video)

// 	foundJob, err := repo.Find(job.ID)

// 	require.Nil(t, err)
// 	require.NotEmpty(t, foundJob.ID)
// 	require.Equal(t, job.ID, foundJob.ID)
// 	require.Equal(t, video.ID, foundJob.VideoID)
// }

func TestJobRepositoryUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)
	job, err := domain.NewJob("output_path", domain.JobStatusPending, video)
	require.Nil(t, err)

	jobRepo := repositories.JobRepositoryDb{Db: db}
	jobRepo.Insert(job)

	job.Status = domain.JobStatusCompleted
	jobRepo.Update(job)

	updatedJob, err := jobRepo.Find(job.ID)

	require.NotEmpty(t, updatedJob.ID)
	require.Nil(t, err)
	require.Equal(t, job.Status, updatedJob.Status)
}
