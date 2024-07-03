package minio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"io"
	"src/internal/domain/track/repository"
	"src/internal/models"
)

const TrackBucket = "track-bucket"

type trackStorage struct {
	client *minio.Client
}

func NewTrackStorage(client *minio.Client) repository.TrackStorage {
	return trackStorage{client: client}
}

func (t trackStorage) UploadObject(track *models.TrackObject) error {
	ctx := context.TODO()

	_, err := t.client.PutObject(ctx,
		TrackBucket,
		track.Source,
		bytes.NewReader(track.Payload),
		int64(len(track.Payload)),
		minio.PutObjectOptions{ContentType: "audio"})

	if err != nil {
		return errors.Wrap(err, "album.minio failed to put")
	}
	return nil
}

func (t trackStorage) LoadObject(track *models.TrackMeta) (*models.TrackObject, error) {
	ctx := context.TODO()

	obj, err := t.client.GetObject(ctx, TrackBucket, track.Source, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "album.minio failed to get")
	}
	defer obj.Close()

	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "album.minio failed to get")
	}
	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "album.minio failed to get")
	}
	ret := models.TrackObject{
		TrackMeta: models.TrackMeta{
			Id:     track.Id,
			Source: track.Source,
			Name:   track.Name,
			Genre:  track.Genre,
		},
		Payload: buffer,
	}

	return &ret, nil
}

func (t trackStorage) DeleteObject(track *models.TrackMeta) error {
	ctx := context.TODO()

	err := t.client.RemoveObject(ctx, TrackBucket, track.Source, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "album.minio failed to delete")
	}

	return nil
}
