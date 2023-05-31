package tos

import (
	"bytes"
	"context"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"strings"
)

const REGION = "cn-beijing"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetObject(endpoint string, accessKey string, secretKey string, bucketName string, objectKey string) (interface{}, error) {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(REGION), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)
	res, err := client.GetObjectV2(context.Background(), &tos.GetObjectV2Input{
		Bucket: bucketName,
		Key:    objectKey,
	})
	checkErr(err)
	return res, nil
}

func PutObject(endpoint string, accessKey string, secretKey string, bucketName string, objectKey string) (interface{}, error) {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(REGION), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)
	res, err := client.PutObjectV2(context.Background(), &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: bucketName,
			Key:    objectKey,
		},
		Content: strings.NewReader("object content"),
	})
	checkErr(err)
	return res, nil

}
func UploadPart(endpoint string, accessKey string, secretKey string, bucketName string, objectKey string) (interface{}, error) {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(REGION), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)
	upload, err := client.CreateMultipartUploadV2(context.Background(), &tos.CreateMultipartUploadV2Input{
		Bucket: bucketName,
		Key:    objectKey,
	})
	checkErr(err)
	buf := make([]byte, 5<<20)
	part1, err := client.UploadPartV2(context.Background(), &tos.UploadPartV2Input{
		UploadPartBasicInput: tos.UploadPartBasicInput{
			Bucket:     bucketName,
			Key:        objectKey,
			UploadID:   upload.UploadID,
			PartNumber: 1,
		},
		Content: bytes.NewReader(buf),
	})
	checkErr(err)
	part2, err := client.UploadPartV2(context.Background(), &tos.UploadPartV2Input{
		UploadPartBasicInput: tos.UploadPartBasicInput{
			Bucket:     bucketName,
			Key:        objectKey,
			UploadID:   upload.UploadID,
			PartNumber: 2,
		},
		Content: bytes.NewReader(buf),
	})
	checkErr(err)
	res, err := client.CompleteMultipartUploadV2(context.Background(), &tos.CompleteMultipartUploadV2Input{
		Bucket:   bucketName,
		Key:      objectKey,
		UploadID: upload.UploadID,
		Parts: []tos.UploadedPartV2{
			{
				PartNumber: part1.PartNumber,
				ETag:       part2.ETag,
			},
			{
				PartNumber: part2.PartNumber,
				ETag:       part2.ETag,
			},
		},
	})
	checkErr(err)
	return res, nil
}
