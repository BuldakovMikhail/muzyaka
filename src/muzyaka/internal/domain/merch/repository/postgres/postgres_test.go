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
	repository := NewMerchRepository(db)

	err = db.Exec("insert into musicians (name, description) values ('test', 'test')").Error
	if err != nil {
		log.Fatal(err)
	}

	merch := models.Merch{
		Id:          0,
		Name:        "test",
		PhotoFiles:  [][]byte{[]byte("test1"), []byte("test2"), []byte("test3")},
		Description: "test",
		OrderUrl:    "test.com",
	}

	id, err := repository.AddMerch(&merch, 1)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	getM, err := repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.PhotoFiles, getM.PhotoFiles)
	assert.Subset(t, getM.PhotoFiles, merch.PhotoFiles)

	merch.PhotoFiles = [][]byte{[]byte("test1"), []byte("test2"), []byte("test3"), []byte("test4")}
	err = repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.PhotoFiles, getM.PhotoFiles)
	assert.Subset(t, getM.PhotoFiles, merch.PhotoFiles)

	merch.PhotoFiles = [][]byte{[]byte("test1"), []byte("test3"), []byte("test5")}
	err = repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.PhotoFiles, getM.PhotoFiles)
	assert.Subset(t, getM.PhotoFiles, merch.PhotoFiles)

	merch.PhotoFiles = [][]byte{[]byte("test6")}
	err = repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.PhotoFiles, getM.PhotoFiles)
	assert.Subset(t, getM.PhotoFiles, merch.PhotoFiles)
}
