package usecases_test

import (
	"encoder/application/repositories"
	usecases "encoder/application/use_cases"
	"encoder/domain"
	"encoder/framework/database"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func setupEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	return video, repo
}

func TestVideoDownload(t *testing.T) {
	setupEnv()
	video, repo := prepare()

	videoUseCase := usecases.NewVideoUseCase()

	videoUseCase.Video = video
	videoUseCase.VideoRepository = repo

	err := videoUseCase.Download(os.Getenv("GOOGLE_STORAGE_BUCKET_NAME"))

	require.Nil(t, err)

	err = videoUseCase.Fragment()

	require.Nil(t, err)

	err = videoUseCase.Encode()

	require.Nil(t, err)

	err = videoUseCase.Finish()

	require.Nil(t, err)
}
