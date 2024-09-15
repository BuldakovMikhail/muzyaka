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

func TestRepo_AddAlbumWithTracks(t *testing.T) {
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
	err = db.Exec("insert into genres (name) values ('test')").Error
	if err != nil {
		log.Fatal(err)
	}

	err = db.Exec("insert into musicians (name, description) values ('test', 'test')").Error
	if err != nil {
		log.Fatal(err)
	}

	repository := NewAlbumRepository(db)

	album := &models.Album{
		Id:        0,
		Name:      "TestName",
		CoverFile: []byte("TestCover"),
		Type:      "LP",
	}

	tracks := []*models.TrackMeta{
		{
			Id:     0,
			Source: "TestSrc1",
			Name:   "TestName1",
			Genre:  "test",
		},
		{
			Id:     0,
			Source: "TestSrc2",
			Name:   "TestName2",
			Genre:  "test",
		},
		{
			Id:     0,
			Source: "TestSrc3",
			Name:   "TestName3",
			Genre:  "test",
		},
	}

	id, err := repository.AddAlbumWithTracksOutbox(album, tracks, 1)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	getAl, err := repository.GetAlbum(id)
	assert.NoError(t, err)
	assert.NotNil(t, getAl)

	assert.Equal(t, getAl, album)

	tracksFromPg, err := repository.GetAllTracksForAlbum(id)
	assert.Equal(t, len(tracksFromPg), len(tracks))
	assert.NoError(t, err)

	err = repository.DeleteAlbumOutbox(id)
	assert.NoError(t, err)

	tracksFromPg, err = repository.GetAllTracksForAlbum(id)
	assert.Equal(t, len(tracksFromPg), 0)
	assert.NoError(t, err)
}
