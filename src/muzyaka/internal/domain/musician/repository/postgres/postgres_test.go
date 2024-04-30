package postgres

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"src/internal/lib/testhelpers"
	"src/internal/models"
	"testing"
)

func TestRepo_PhotosChange(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("error terminating postgres container: %s", err)
		}
	}()

	db, err := gorm.Open(postgres.Open(pgContainer.ConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	repository := NewMusicianRepository(db)

	musician := models.Musician{
		Id:          0,
		Name:        "test",
		PhotoFiles:  [][]byte{[]byte("test1"), []byte("test2"), []byte("test3")},
		Description: "test",
	}

	id, err := repository.AddMusician(&musician)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	getM, err := repository.GetMusician(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &musician)
	assert.Subset(t, musician.PhotoFiles, getM.PhotoFiles)
	assert.Subset(t, getM.PhotoFiles, musician.PhotoFiles)

	musician.PhotoFiles = [][]byte{[]byte("test1"), []byte("test2"), []byte("test3"), []byte("test4")}
	err = repository.UpdateMusician(&musician)
	assert.NoError(t, err)

	getM, err = repository.GetMusician(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &musician)
	assert.Subset(t, musician.PhotoFiles, getM.PhotoFiles)
	assert.Subset(t, getM.PhotoFiles, musician.PhotoFiles)

	musician.PhotoFiles = [][]byte{[]byte("test1"), []byte("test3"), []byte("test5")}
	err = repository.UpdateMusician(&musician)
	assert.NoError(t, err)

	getM, err = repository.GetMusician(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &musician)
	assert.Subset(t, musician.PhotoFiles, getM.PhotoFiles)
	assert.Subset(t, getM.PhotoFiles, musician.PhotoFiles)

	musician.PhotoFiles = [][]byte{[]byte("test6")}
	err = repository.UpdateMusician(&musician)
	assert.NoError(t, err)

	getM, err = repository.GetMusician(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &musician)
	assert.Subset(t, musician.PhotoFiles, getM.PhotoFiles)
	assert.Subset(t, getM.PhotoFiles, musician.PhotoFiles)
}
