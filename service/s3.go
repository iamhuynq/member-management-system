package service

import (
	"os"
	"path/filepath"

	"github.com/tribalmedia/vista/setting"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ConnectToS3 is ...
func ConnectToS3() *s3.S3 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	creds := ec2rolecreds.NewCredentials(sess)
	_, err := creds.Get()
	if err != nil {
		Logger.Error(err.Error())
	}

	cfg := aws.NewConfig().WithCredentials(creds).WithRegion(setting.Region)
	svc := s3.New(sess, cfg)

	return svc
}

// UploadImageToS3 is ...
func UploadImageToS3(pathToFile string, folderForSaving int) {
	file, err := os.Open(pathToFile)
	if err != nil {
		Logger.Fatal(err.Error())
	}
	defer os.Remove(pathToFile)

	svc := ConnectToS3()
	var pathOnS3 string
	if folderForSaving == setting.MemberFolderType {
		pathOnS3 = setting.S3MemberFolder + filepath.Base(pathToFile)
	} else if folderForSaving == setting.TeamFolderType {
		pathOnS3 = setting.S3TeamFolder + filepath.Base(pathToFile)
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(setting.BucketName),
		Key:         aws.String(pathOnS3),
		Body:        file,
		ContentType: aws.String("image/jpeg"),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		Logger.Error(err.Error())
	}
}

// DeleteImageFromS3 is ...
func DeleteImageFromS3(fileName string, folderForSaving int) {
	svc := ConnectToS3()
	var pathOnS3 string
	if folderForSaving == setting.MemberFolderType {
		pathOnS3 = setting.S3MemberFolder + fileName
	} else if folderForSaving == setting.TeamFolderType {
		pathOnS3 = setting.S3TeamFolder + fileName
	}

	// Delete the item
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(setting.BucketName),
		Key:    aws.String(pathOnS3),
	})

	if err != nil {
		Logger.Error(err.Error())
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(setting.BucketName),
		Key:    aws.String(pathOnS3),
	})

	if err != nil {
		Logger.Error(err.Error())
	}
}
