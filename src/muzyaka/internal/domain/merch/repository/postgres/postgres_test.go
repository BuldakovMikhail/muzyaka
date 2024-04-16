package postgres

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	repository2 "src/internal/domain/merch/repository"
	"src/internal/lib/testhelpers"
	"src/internal/models"
	"testing"
)

type MerchRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  repository2.MerchRepository
	ctx         context.Context
	db          *gorm.DB
}

func (suite *MerchRepoTestSuite) SetupSuite() {
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

	repository := NewMerchRepository(db)
	suite.repository = repository

	suite.db = db
}

func (suite *MerchRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *MerchRepoTestSuite) TestPhotosChange() {
	t := suite.T()

	merch := models.Merch{
		Id:          0,
		Name:        "test",
		Photos:      []string{"test1", "test2", "test3"},
		Description: "test",
		OrderUrl:    "test.com",
	}

	id, err := suite.repository.AddMerch(&merch)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	getM, err := suite.repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.Photos, getM.Photos)
	assert.Subset(t, getM.Photos, merch.Photos)

	merch.Photos = []string{"test1", "test2", "test3", "test4"}
	err = suite.repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = suite.repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.Photos, getM.Photos)
	assert.Subset(t, getM.Photos, merch.Photos)

	merch.Photos = []string{"test1", "test3", "test5"}
	err = suite.repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = suite.repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.Photos, getM.Photos)
	assert.Subset(t, getM.Photos, merch.Photos)

	merch.Photos = []string{"test6"}
	err = suite.repository.UpdateMerch(&merch)
	assert.NoError(t, err)

	getM, err = suite.repository.GetMerch(id)
	assert.NoError(t, err)
	assert.NotNil(t, getM)
	assert.Equal(t, getM, &merch)
	assert.Subset(t, merch.Photos, getM.Photos)
	assert.Subset(t, getM.Photos, merch.Photos)
}

func TestAlbumRepoTestSuite(t *testing.T) {
	suite.Run(t, new(MerchRepoTestSuite))
}
