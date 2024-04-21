package minio

import (
	"context"
	minio2 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	repository2 "src/internal/domain/track/repository"
	"src/internal/lib/testhelpers"
	"src/internal/models"
	"testing"
)

type TrackStorageTestSuite struct {
	suite.Suite
	minioContainer testhelpers.Container
	storage        repository2.TrackStorage
	ctx            context.Context
	client         *minio2.Client
}

func (suite *TrackStorageTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	container, err := testhelpers.Start(suite.ctx, testhelpers.Options{
		ImageTag:     "RELEASE.2024-01-16T16-07-38Z",
		RootUser:     "3846587325",
		RootPassword: "te782tcb7tr3va7brkwev7awst",
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	suite.minioContainer = container

	minioURI := container.ConnectionURI()
	minioClient, err := minio2.New(minioURI, &minio2.Options{
		Creds:  credentials.NewStaticV4(container.RootUser, container.RootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	suite.client = minioClient

	storage := NewTrackStorage(minioClient)

	suite.storage = storage

	err = suite.client.MakeBucket(suite.ctx, TrackBucket, minio2.MakeBucketOptions{})
	if err != nil {
		log.Fatalf("failed to create bucket: %s", err)
	}
}

func (suite *TrackStorageTestSuite) TearDownSuite() {
	suite.minioContainer.Terminate(suite.ctx)
}

func (suite *TrackStorageTestSuite) TestTrackStorage() {
	t := suite.T()

	track := models.TrackObject{
		TrackMeta: models.TrackMeta{
			Id:     0,
			Source: "aboba",
			Name:   "aboba",
			Genre:  "aboba",
		},
		Payload:     []byte{1, 2, 3},
		PayloadSize: 3,
	}

	err := suite.storage.UploadObject(&track)
	assert.NoError(t, err)

	trackLoaded, err := suite.storage.LoadObject(track.ExtractMeta())
	assert.NoError(t, err)

	assert.Equal(t, trackLoaded.PayloadSize, track.PayloadSize)
	assert.Equal(t, trackLoaded.Payload, track.Payload)

	track.Payload = []byte{4, 5, 6}
	err = suite.storage.UploadObject(&track)
	assert.NoError(t, err)

	trackLoaded, err = suite.storage.LoadObject(track.ExtractMeta())
	assert.NoError(t, err)

	assert.Equal(t, trackLoaded.PayloadSize, track.PayloadSize)
	assert.Equal(t, trackLoaded.Payload, track.Payload)

	err = suite.storage.DeleteObject(track.ExtractMeta())
	assert.NoError(t, err)

	trackLoaded, err = suite.storage.LoadObject(track.ExtractMeta())
	assert.Error(t, err)
}

func TestTrackStorageTestSuite(t *testing.T) {
	suite.Run(t, new(TrackStorageTestSuite))
}
