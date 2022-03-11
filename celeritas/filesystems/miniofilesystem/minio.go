package miniofilesystem

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	Endpoint string
	Key      string
	Secret   string
	UseSSL   bool
	Region   string
	Bucket   string
}

func (m *Minio) getCredentials() *minio.Client {
	client, err := minio.New(m.Endpoint, *minio.Options{
		Creds:  credentials.NewStaticV4(m.Key, m.Secret, ""),
		Secure: m.UseSSL,
	})
	if err != nil {
		log.Println(err)
	}

	return client
}

func (m *Minio) Put(fileName, folder string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// .Base() split the path on "/" and return the last element
	objectName := path.Base(fileName)
	// client in order to make the connection to the endpoint
	client := m.getCredentials()
	uploadInfo, err := client.FPutObject(ctx, m.Bucket, fmt.Sprintf("%s/%s", folder, objectName), fileName, minio.PutObjectOptions{})
	if err != nil {
		log.Println("failed with FPutObject")
		log.Println(err)
		log.Println("uploadInfo: ", uploadInfo)
		return err
	}

	return nil
}
