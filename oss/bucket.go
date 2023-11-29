package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/succko/hera/global"
	"io/ioutil"
	"os"
	"strings"
)

type bucket struct {
}

var Bucket = new(bucket)

func (bucket *bucket) PutObject(objectKey string, content string) {
	// 指定Object存储类型为低频访问。
	storageType := oss.ObjectStorageClass(oss.StorageStandard)
	// 指定Object访问权限为私有。
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)
	// 将字符串"Hello OSS"上传至exampledir目录下的exampleobject.txt文件。
	err := global.App.Oss.PutObject(objectKey, strings.NewReader(content), storageType, objectAcl)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}

func (bucket *bucket) SelectObject(objectKey string) (string, error) {
	selReq := oss.SelectRequest{}
	// 使用SELECT语句查询文件中的数据。
	selReq.Expression = `select * from ossobject`
	body, err := global.App.Oss.SelectObject(objectKey, selReq)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 读取内容。
	fc, err := ioutil.ReadAll(body)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	defer body.Close()
	fmt.Println(string(fc))
	return string(fc), nil
}
