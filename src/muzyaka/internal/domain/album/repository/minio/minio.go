package minio

import "github.com/minio/minio-go/v7"

type trackStorage struct {
	client *minio.Client
}
