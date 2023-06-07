package handle_volc

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos/enum"
)

func SetPicPublic(ctx context.Context, picUrls []string) bool {

	var (
		ak = "AKLTMDA2MDBhYmY4ZWEzNGY1MWIxNjlhMDE3MDVkMDI3NzA"
		sk = "T1RneFlUSTRPVE5sTUdNd05EUmtNemc0T1dWak1ERTJNakZqWVdObE1ESQ=="

		endpoint = "tos-cn-beijing.volces.com"
		region   = "cn-beijing"

		bucket_name = "ttb6cfc07229f306bb01-env-0inwysq77q"
	)
	credential := tos.NewStaticCredentials(ak, sk)
	client, err := tos.NewClientV2(endpoint, tos.WithCredentials(credential), tos.WithRegion(region))
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	for _, v := range picUrls {
		objectKey := v
		// 设置对象 ACL
		input := &tos.PutObjectACLInput{
			Bucket: bucket_name,
			Key:    objectKey,
			// 如果桶开启的多版本，通过设置 VersionID 来设置指定版本
			ACL: enum.ACLPublicRead,
		}
		output, err := client.PutObjectACL(ctx, input)

		if err != nil || output == nil {
			msg := fmt.Sprintf("PutObjectACL Request ID:", output.RequestID)
			log.Error(msg)
			return false
		}

	}

	// 使用结束后，关闭 client
	client.Close()

	return true
}
