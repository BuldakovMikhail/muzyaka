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

	merch := models.Merch{
		Id:          0,
		Name:        "test",
		Photos:      []string{"test1", "test2", "test3"},
		Description: "test",
		OrderUrl:    "test.com",
	}

	id, err := repository.AddMerch(&merch)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	getM, err := repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.Photos, getM.Photos)
	assert.Subset(t, getM.Photos, merch.Photos)

	merch.Photos = []string{"test1", "test2", "test3", "test4"}
	err = repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.Photos, getM.Photos)
	assert.Subset(t, getM.Photos, merch.Photos)

	merch.Photos = []string{"test1", "test3", "test5"}
	err = repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.Photos, getM.Photos)
	assert.Subset(t, getM.Photos, merch.Photos)

	merch.Photos = []string{"test6"}
	err = repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.Photos, getM.Photos)
	assert.Subset(t, getM.Photos, merch.Photos)
}
