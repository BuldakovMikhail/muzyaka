package postgres

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"src/internal/domain/album/repository"
	"src/internal/lib/testhelpers"
	"src/internal/models"
	"testing"
)

type AlbumRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  repository.AlbumRepository
	ctx         context.Context
}

func (suite *AlbumRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	db, err := gorm.Open(postgres.Open(suite.pgContainer.ConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Exec("insert into genres (name) values ('test')").Error
	if err != nil {
		log.Fatal(err)
	}

	repository := NewAlbumRepository(db)
	suite.repository = repository
}

func (suite *AlbumRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *AlbumRepoTestSuite) TestAddAlbumWithTracks() {
	t := suite.T()

	album := &models.Album{
		Id:    0,
		Name:  "TestName",
		Cover: "TestCover",
		Type:  "LP",
	}

	tracks := []*models.Track{
		&models.Track{
			Id:     0,
			Source: "TestSrc1",
			Name:   "TestName1",
			Genre:  "test",
		},
		&models.Track{
			Id:     0,
			Source: "TestSrc2",
			Name:   "TestName2",
			Genre:  "test",
		},
		&models.Track{
			Id:     0,
			Source: "TestSrc3",
			Name:   "TestName3",
			Genre:  "test",
		},
	}

	id, err := suite.repository.AddAlbumWithTracks(album, tracks)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	getAl, err := suite.repository.GetAlbum(id)
	assert.NoError(t, err)
	assert.NotNil(t, getAl)

	assert.Equal(t, getAl, album)
}

//	func (suite *AlbumRepoTestSuite) TestGetCustomerByEmail() {
//		t := suite.T()
//
//		customer, err := suite.repository.GetCustomerByEmail(suite.ctx, "john@gmail.com")
//		assert.NoError(t, err)
//		assert.NotNil(t, customer)
//		assert.Equal(t, "John", customer.Name)
//		assert.Equal(t, "john@gmail.com", customer.Email)
//	}
func TestAlbumRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumRepoTestSuite))
}
