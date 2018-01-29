package main

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"time"
	"fmt"
	"os"
	"log"
)

var Keep = &cobra.Command{
	Use:   "keep [amount] [region] [bucket]",
	Short: "Keep amount of files in S3",
	Long: `Keep amount of files in S3.
If you want to keep 5 files, will remove file older No.6`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		amount, _ := strconv.Atoi(args[0])
		region := args[1]
		bucket := args[2]


		if !hasAwsAccess() {
			log.Fatalln("Environment variable is empty.", "AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY")
		}

		c := NewS3Client(region, bucket)

		s3Objects, err := c.fetchObjects()

		if err != nil {
			panic(err)
		}

		err = c.deleteOldObjects(s3Objects, amount)

		if err != nil {
			panic(err)
		}
	},
}

type S3Client struct {
	region string
	bucket string
}

type S3Object struct {
	key       string
	createdAt time.Time
}

func hasAwsAccess() bool {
	awsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	return awsAccessKeyId != "" && awsSecretAccessKey != ""
}

func NewS3Client(region string, bucket string) *S3Client {
	return &S3Client{region, bucket}
}

func NewService(region *string) *s3.S3 {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*region),
	}))

	return s3.New(sess)
}

// S3からすべてのオブジェクトを取得する
func (s *S3Client) fetchObjects() ([]*S3Object, error) {
	svc := NewService(&s.region)

	response, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: &s.bucket})

	if err != nil {
		return nil, err
	}

	var objects []*S3Object

	for _, object := range response.Contents {
		s3Object := &S3Object{*object.Key, *object.LastModified}
		objects = append(objects, s3Object)
	}

	return objects, nil
}

// 残さないオブジェクトを削除する
func (s *S3Client) deleteOldObjects(s3Objects []*S3Object, amount int) error {
	deleteObjectList := deleteObjectList(s3Objects, amount)
	if len(deleteObjectList) < 1 {
		fmt.Println("No delete list.")
		return nil
	}

	var objects []*s3.ObjectIdentifier
	for _, o := range deleteObjectList {
		objects = append(objects, &s3.ObjectIdentifier{Key: &o.key})
		fmt.Printf("Delete: %v %v\n", o.createdAt, o.key)
	}

	deleteObjectsInput := &s3.DeleteObjectsInput{
		Bucket: &s.bucket,
		Delete: &s3.Delete{
			Objects: objects,
		},
	}

	svc := NewService(&s.region)
	_, err := svc.DeleteObjects(deleteObjectsInput)

	if err != nil {
		panic(err)
	}

	return nil
}

// 削除するオブジェクトを配列で返す
func deleteObjectList(s3Objects []*S3Object, amount int) []*S3Object {
	if len(s3Objects) <= amount {
		return make([]*S3Object, 0, 0)
	}

	deleteObjectList := s3Objects
	for i := 0; i < amount; i++ {
		deleteObjectList = remove(deleteObjectList, latestObject(deleteObjectList))
	}

	return deleteObjectList
}

// 最新のオブジェクトを返す
func latestObject(s3Objects []*S3Object) *S3Object {
	var result *S3Object

	for _, o := range s3Objects {
		if result == nil || o.createdAt.After(result.createdAt) {
			result = o
		}
	}
	fmt.Printf("%v %v\n", result.createdAt, result.key)

	return result
}

// 配列から削除するUtil的存在
func remove(s3Objects []*S3Object, search *S3Object) []*S3Object {
	var result []*S3Object
	for _, o := range s3Objects {
		if o != search {
			result = append(result, o)
		}
	}
	return result
}
