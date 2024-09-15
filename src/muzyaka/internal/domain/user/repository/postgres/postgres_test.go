package postgres

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	repository2 "src/internal/domain/user/repository"
	"src/internal/lib/testhelpers"
	"src/internal/models"
	"testing"
)

type UserRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  repository2.UserRepository
	ctx         context.Context
	db          *gorm.DB
}

func (suite *UserRepoTestSuite) SetupSuite() {
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

	repository := NewUserRepository(db)
	suite.repository = repository

	suite.db = db
}

func (suite *UserRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *UserRepoTestSuite) TestUserRepo() {
	t := suite.T()

	user := models.User{
		Id:       0,
		Name:     "Test",
		Password: "Test",
		Role:     "musician",
		Email:    "test@test.ru",
	}

	id, err := suite.repository.AddUser(&user)
	assert.NoError(t, err)
	assert.Equal(t, id, uint64(1))

	pgUser, err := suite.repository.GetUser(id)
	assert.NoError(t, err)
	assert.NotNil(t, pgUser)
	assert.Equal(t, pgUser, &user)

	user.Name = "changedTest"
	err = suite.repository.UpdateUser(&user)
	assert.NoError(t, err)

	pgUser, err = suite.repository.GetUser(id)
	assert.NoError(t, err)
	assert.NotNil(t, pgUser)
	assert.Equal(t, pgUser, &user)

	err = suite.repository.DeleteUser(id)
	assert.NoError(t, err)

	pgUser, err = suite.repository.GetUser(id)
	assert.Error(t, err)
	assert.Nil(t, pgUser)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
