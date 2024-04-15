package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	repository2 "src/internal/domain/user/repository"
	"src/internal/models"
	"src/internal/models/dao"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository2.UserRepository {
	return &userRepository{db: db}
}

func (u userRepository) GetUser(id uint64) (*models.User, error) {
	var user dao.User

	tx := u.db.Where("id = ?", id).Take(&user)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table user)")
	}

	return dao.ToModelUser(&user), nil
}

func (u userRepository) UpdateUser(user *models.User) error {
	pgUser := dao.ToPostgresUser(user)

	tx := u.db.Omit("id").Updates(&pgUser)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table user)")
	}

	return nil
}

func (u userRepository) AddUser(user *models.User) (uint64, error) {
	pgUser := dao.ToPostgresUser(user)

	tx := u.db.Create(&pgUser)
	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "database error (table user)")
	}

	user.Id = pgUser.ID
	return pgUser.ID, nil
}

func (u userRepository) DeleteUser(id uint64) error {
	tx := u.db.Delete(&dao.User{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table user)")
	}

	return nil
}
