package s3

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

func LoadConfig() bool {
	viper.SetConfigFile("config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("viper read configuration failed, %s", err)
		return false
	}

	return true
}

var sess *session.Session

func MakeS3Session() error {
	if !LoadConfig() {
		return errors.New("s3 connection get failed")
	}

	accessKeyId := viper.GetString("AccessKeyID")
	secretKey := viper.GetString("SecretAccessKey")
	rigion := viper.GetString("MyRegion")
	var err error
	sess, err = session.NewSession(
		&aws.Config{
			Region: aws.String(rigion),
			Credentials: credentials.NewStaticCredentials(
				accessKeyId,
				secretKey,
				"",
			),
		},
	)

	if err != nil {
		return errors.New("aws connection get failed")
	}

	return nil
}

func UploadFileToS3(filePath string, keyName, fileContent string) error {
	if sess == nil {
		if err := MakeS3Session(); err != nil {
			return err
		}
	}
	svc := s3.New(sess)
	var bucket string
	if fileContent == "video/mp4" {
		// check for video file
		bucket = viper.GetString("MyVideoBucket")

	} else {
		// check for thumbnail as well
		if strings.Contains(keyName, "thumbnail") {
			bucket = viper.GetString("MyThumbnail")
		} else {
			bucket = viper.GetString("MyBucket")
		}
	}

	if svc != nil && bucket != "" {

		tempFile, err := os.Open(filePath)

		if err != nil {
			return err
		}
		defer tempFile.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, tempFile); err != nil {
			return err
		}

		result, err := svc.PutObject(&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(keyName),
			Body:        bytes.NewReader(buf.Bytes()),
			ContentType: &fileContent,
		})

		if err != nil {
			return err
		}

		log.Println("Upload File Location : ", result)
		log.Println("Image Has Been Uploaded Successfully")
		return nil
	}

	return errors.New("something went wrong")
}

func GetTheObject(fileKey string) (string, error) {
	if sess == nil {
		if err := MakeS3Session(); err != nil {
			return "", err
		}
	}
	svc := s3.New(sess)

	bucketName := viper.GetString("MyBucket")
	expiration := time.Now().Add(24 * time.Hour) // URL will expire in 24 hours
	duration := time.Until(expiration)
	// Generate a pre-signed URL for the object
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})

	url, err := req.Presign(duration)
	if err != nil {
		fmt.Println("Error generating pre-signed URL:", err)
		return "", err
	}

	return url, nil
}

func GetVideoObjectUrl(fileKey string) (string, error) {
	if sess == nil {
		if err := MakeS3Session(); err != nil {
			return "", err
		}
	}
	svc := s3.New(sess)

	bucketName := viper.GetString("MyVideoBucket")
	expiration := time.Now().Add(24 * time.Hour) // URL will expire in 24 hours
	duration := time.Until(expiration)
	// Generate a pre-signed URL for the object
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})

	url, err := req.Presign(duration)
	if err != nil {
		fmt.Println("Error generating pre-signed URL:", err)
		return "", err
	}

	return url, nil
}

func GetVideoObjectInput(fileKey string) *s3.GetObjectOutput {
	if sess == nil {
		if err := MakeS3Session(); err != nil {
			return nil
		}
	}
	svc := s3.New(sess)

	bucketName := viper.GetString("MyVideoBucket")

	// Generate a pre-signed URL for the object
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func GetVideoThumbnailObjectUrl(fileKey string) (string, error) {
	if sess == nil {
		if err := MakeS3Session(); err != nil {
			return "", err
		}
	}
	svc := s3.New(sess)

	bucketName := viper.GetString("MyThumbnail")
	expiration := time.Now().Add(24 * time.Hour) // URL will expire in 24 hours
	duration := time.Until(expiration)
	// Generate a pre-signed URL for the object
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})

	url, err := req.Presign(duration)
	if err != nil {
		fmt.Println("Error generating pre-signed URL:", err)
		return "", err
	}

	return url, nil
}
