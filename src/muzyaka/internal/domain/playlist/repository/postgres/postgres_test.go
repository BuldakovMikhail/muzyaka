package postgres

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	postgres2 "src/internal/domain/album/repository/postgres"
	"src/internal/lib/testhelpers"
	"src/internal/models"
	"testing"
)

func TestRepo_CreatePlaylist(t *testing.T) {
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

	err = db.Exec("insert into users (name, email, password)\nvalues ('Sasha', 'test3@gmail.test', 'aaaaaa')").Error
	if err != nil {
		log.Fatal(err)
	}

	repository := postgres2.NewAlbumRepository(db)

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
	require.NoError(t, err)
	require.NotNil(t, id)

	repositoryPlaylist := NewPlaylistRepository(db)

	playlist := models.Playlist{
		Id:          0,
		Name:        "testp",
		CoverFile:   []byte{1, 2, 3},
		Description: "testp",
	}

	id, err = repositoryPlaylist.AddPlaylist(&playlist, 1)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	pgPlaylist, err := repositoryPlaylist.GetPlaylist(id)
	assert.Equal(t, pgPlaylist, &playlist)
	assert.NoError(t, err)
}
