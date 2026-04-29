package usecases

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"io"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

type VideoUseCase struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoUseCase() VideoUseCase {
	return VideoUseCase{}
}

func (v *VideoUseCase) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	//stablish the connection with the bucket
	bucket := client.Bucket(bucketName)

	//pass the videoId to be downloaded
	obj := bucket.Object(v.Video.FilePath)

	//download the video
	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("Video %v was successfully stored at %v", v.Video.ID, v.Video.FilePath)

	return nil
}

func (v *VideoUseCase) Fragment() error {
	err := os.MkdirAll(os.Getenv("LOCAL_STORAGE_PATH")+ "/" + v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4"

	target := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	log.Printf("Video %v was successfully fragmented", v.Video.ID)

	printOutput(output)

	return nil
}

func (v *VideoUseCase) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "--o")
	cmdArgs = append(cmdArgs, os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mpd")
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "--/opt/bento4/bin")

	cmd := exec.Command("mp4dash", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	log.Printf("Video %v was successfully encoded", v.Video.ID)

	printOutput(output)

	return nil
}

func (v *VideoUseCase) Finish() error{
	err := os.Remove(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		log.Printf("Error deleting video %v: %v", v.Video.ID, err)
		return err
	}
	err = os.Remove(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".frag")
	if err != nil {
		log.Printf("Error deleting video %v: %v", v.Video.ID, err)
		return err
	}

	err = os.RemoveAll(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID)
	if err != nil {
		log.Printf("Error deleting video %v: %v", v.Video.ID, err)
		return err
	}

	log.Printf("Video %v files were successfully deleted", v.Video.ID)

	return nil
}


func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("Output: %s", string(out))
	}
}
