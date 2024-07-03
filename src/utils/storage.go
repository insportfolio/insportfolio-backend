package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

func UploadFile(file *multipart.FileHeader) (string, error) {
	region := os.Getenv("AWS_S3_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	bucket := os.Getenv("AWS_S3_BUCKET")

	nameSplitted := strings.Split(file.Filename, ".")
	ext := nameSplitted[len(nameSplitted)-1]
	name := strings.ReplaceAll(nameSplitted[0], " ", "")
	key := "/portfolios/" + name + "_" + uuid.New().String() + "." + ext

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return "", errors.New("error creating session")
	}
	svc := s3.New(sess)
	f, err := file.Open()
	if err != nil {
		return "", errors.New("error opening file")
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("error opening file")
	}
	return key, nil
}

func GetImageFullPath(path string) string {
	bucket := os.Getenv("AWS_S3_BUCKET")
	region := os.Getenv("AWS_S3_REGION")
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, path)

}
