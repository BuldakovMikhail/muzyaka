package minio

import (
	"context"
	minio2 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"log"
	"src/internal/lib/testhelpers"
	"src/internal/models"
	"testing"
)

func TestRepo_TrackStorageAdd(t *testing.T) {
	ctx := context.Background()

	minioContainer, err := testhelpers.Start(ctx, testhelpers.Options{
		ImageTag:     "RELEASE.2024-01-16T16-07-38Z",
		RootUser:     "3846587325",
		RootPassword: "te782tcb7tr3va7brkwev7awst",
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer minioContainer.Terminate(ctx)

	minioURI := minioContainer.ConnectionURI()
	client, err := minio2.New(minioURI, &minio2.Options{
		Creds:  credentials.NewStaticV4(minioContainer.RootUser, minioContainer.RootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	storage := NewTrackStorage(client)

	err = client.MakeBucket(ctx, TrackBucket, minio2.MakeBucketOptions{})
	if err != nil {
		log.Fatalf("failed to create bucket: %s", err)
	}

	track := models.TrackObject{
		TrackMeta: models.TrackMeta{
			Id:     0,
			Source: "aboba",
			Name:   "aboba",
			Genre:  "aboba",
		},
		Payload: []byte{1, 2, 3},
	}

	err = storage.UploadObject(&track)
	assert.NoError(t, err)

	trackLoaded, err := storage.LoadObject(track.ExtractMeta())
	assert.NoError(t, err)

	assert.Equal(t, trackLoaded.Payload, track.Payload)
}

func TestRepo_TrackStorageUpdate(t *testing.T) {
	ctx := context.Background()

	minioContainer, err := testhelpers.Start(ctx, testhelpers.Options{
		ImageTag:     "RELEASE.2024-01-16T16-07-38Z",
		RootUser:     "3846587325",
		RootPassword: "te782tcb7tr3va7brkwev7awst",
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer minioContainer.Terminate(ctx)

	minioURI := minioContainer.ConnectionURI()
	client, err := minio2.New(minioURI, &minio2.Options{
		Creds:  credentials.NewStaticV4(minioContainer.RootUser, minioContainer.RootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	storage := NewTrackStorage(client)

	err = client.MakeBucket(ctx, TrackBucket, minio2.MakeBucketOptions{})
	if err != nil {
		log.Fatalf("failed to create bucket: %s", err)
	}

	track := models.TrackObject{
		TrackMeta: models.TrackMeta{
			Id:     0,
			Source: "aboba",
			Name:   "aboba",
			Genre:  "aboba",
		},
		Payload: []byte{1, 2, 3},
	}

	err = storage.UploadObject(&track)
	assert.NoError(t, err)

	trackLoaded, err := storage.LoadObject(track.ExtractMeta())
	assert.NoError(t, err)

	assert.Equal(t, trackLoaded.Payload, track.Payload)

	track.Payload = []byte{4, 5, 6}
	err = storage.UploadObject(&track)
	assert.NoError(t, err)

	trackLoaded, err = storage.LoadObject(track.ExtractMeta())
	assert.NoError(t, err)

	assert.Equal(t, trackLoaded.Payload, track.Payload)

}

func TestRepo_TrackStorageDelete(t *testing.T) {
	ctx := context.Background()

	minioContainer, err := testhelpers.Start(ctx, testhelpers.Options{
		ImageTag:     "RELEASE.2024-01-16T16-07-38Z",
		RootUser:     "3846587325",
		RootPassword: "te782tcb7tr3va7brkwev7awst",
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer minioContainer.Terminate(ctx)

	minioURI := minioContainer.ConnectionURI()
	client, err := minio2.New(minioURI, &minio2.Options{
		Creds:  credentials.NewStaticV4(minioContainer.RootUser, minioContainer.RootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	storage := NewTrackStorage(client)

	err = client.MakeBucket(ctx, TrackBucket, minio2.MakeBucketOptions{})
	if err != nil {
		log.Fatalf("failed to create bucket: %s", err)
	}

	track := models.TrackObject{
		TrackMeta: models.TrackMeta{
			Id:     0,
			Source: "aboba",
			Name:   "aboba",
			Genre:  "aboba",
		},
		Payload: []byte{1, 2, 3},
	}

	err = storage.UploadObject(&track)
	assert.NoError(t, err)

	trackLoaded, err := storage.LoadObject(track.ExtractMeta())
	assert.NoError(t, err)

	assert.Equal(t, trackLoaded.Payload, track.Payload)

	err = storage.DeleteObject(track.ExtractMeta())
	assert.NoError(t, err)

	trackLoaded, err = storage.LoadObject(track.ExtractMeta())
	assert.Error(t, err)
}
