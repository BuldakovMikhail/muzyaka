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

func (u userRepository) IsTrackLiked(userId uint64, trackId uint64) (bool, error) {
	tx := u.db.
		Where("user_id = ? AND track_id = ?", userId, trackId).
		First(&dao.UserTrack{})

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if tx.Error != nil {
		return false, errors.Wrap(tx.Error, "database error (table user_track)")
	}

	return true, nil
}

func (u userRepository) GetAllLikedTracks(userId uint64) ([]uint64, error) {
	var rels []*dao.UserTrack
	tx := u.db.Where("user_id = ?", userId).Limit(dao.MaxLimit).Find(&rels)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table user_track)")
	}

	var ans []uint64

	for _, v := range rels {
		ans = append(ans, v.TrackId)
	}

	return ans, nil
}

func (u userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user dao.User
	tx := u.db.Where("email = ?", email).Take(&user)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users_musicians)")
	}

	return dao.ToModelUser(&user), nil
}

func (u userRepository) LikeTrack(userId uint64, trackId uint64) error {
	tx := u.db.Create(dao.UserTrack{
		TrackId: trackId,
		UserId:  userId,
	})

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table users_track)")
	}

	return nil
}

func (u userRepository) DislikeTrack(userId uint64, trackId uint64) error {
	tx := u.db.Delete(dao.UserTrack{}, "user_id = ? AND track_id = ?", userId, trackId)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table users_track)")
	}

	if tx.RowsAffected == 0 {
		return models.ErrNothingToDelete
	}

	return nil
}

func (u userRepository) GetUser(id uint64) (*models.User, error) {
	var user dao.User

	tx := u.db.Where("id = ?", id).Take(&user)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table user)")
	}

	return dao.ToModelUser(&user), nil
}

func (u userRepository) AddUserWithMusician(musician *models.Musician, user *models.User) (uint64, error) {
	pgMusician := dao.ToPostgresMusician(musician)
	pgUser := dao.ToPostgresUser(user)

	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pgMusician).Error; err != nil {
			return err
		}

		temp := musician
		temp.Id = pgMusician.ID
		pgMusicianPhotos := dao.ToPostgresMusicianPhotos(temp)
		if len(pgMusicianPhotos) != 0 {
			if err := tx.Create(&pgMusicianPhotos).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(&pgUser).Error; err != nil {
			return err
		}

		pgRelation := dao.UserMusician{
			UserId:     pgUser.ID,
			MusicianId: pgMusician.ID,
		}

		if err := tx.Create(&pgRelation).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "database error (table musician)")
	}
	musician.Id = pgMusician.ID
	user.Id = pgUser.ID
	return pgUser.ID, nil
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

	if tx.RowsAffected == 0 {
		return models.ErrNothingToDelete
	}

	return nil
}
