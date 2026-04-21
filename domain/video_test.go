package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()
	require.Error(t, err)

}

func TestVideoIdIsNotUuid(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "invalid-id"
	video.ResourceID = "invalid-resource-id"
	video.FilePath = "invalid-file-path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Error(t, err)

	// video.ID = ""
	// err = video.Validate()
	// require.Error(t, err)
}
