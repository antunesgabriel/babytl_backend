package utils

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	AccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
	SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	Region          = os.Getenv("AWS_REGION")
)

const ACL = "public-read"
const BUCKET_NAME = "babytimeline"

type S3Handler struct {
	Session *session.Session
	Bucket  string
	Region  string
}

func NewS3Handler() (*S3Handler, error) {
	s, err := session.NewSession(&aws.Config{Region: aws.String(Region)})
	sess := session.Must(s, err)

	handler := S3Handler{
		Session: sess,
		Bucket:  BUCKET_NAME,
		Region:  string(Region),
	}

	return &handler, err
}

func (h *S3Handler) UploadFile(key string, filename string, folder string) (string, error) {

	uploader := s3manager.NewUploader(h.Session)

	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("os.Open - filename: %s, err: %v", filename, err)

		return "", err
	}

	defer file.Close()

	filePath := path.Join(folder, path.Base(filename))

	_, errUpload := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.Bucket),
		Key:    aws.String(filePath),
		ACL:    aws.String(ACL),
		Body:   file,
	})

	fileUrl := fmt.Sprintf("https://%s.s3.sa-east-1.amazonaws.com/%s", h.Bucket, filePath)

	return fileUrl, errUpload
}

func (h *S3Handler) ReadFile(key string) {
	//
}
