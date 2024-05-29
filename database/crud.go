package database

import (
	"context"
	"userservice/entity"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func CreateUser(ctx context.Context, user *entity.User, logger *zap.SugaredLogger) error {
	logger.Infof("createdat %s", user.CreatedAt)
	var id string
	_, err := db.NewInsert().Model(user).ExcludeColumn("created_at").Returning("id", "created_at").Exec(ctx, &id)
	if err != nil {
		logger.Errorw("Failed to insert user", "error", err)
		return err
	}
	logger.Infof("Inserted %s", id)

	if err != nil {
		logger.Errorw("Failed to retrieve ID of inserted user", "error", err)
		return err
	}

	// Update the user object with the generated ID
	user.Id = id
	return nil
}

func ReadAll(ctx context.Context, authorId string, logger *zap.SugaredLogger) ([]entity.User, error) {
	logger.Info("Retrieving users by author ID")

	var users []entity.User
	err := db.NewSelect().Model(&users).Order("created_at DESC").Scan(ctx)
	if err != nil {
		logger.Warnw("Failed to retrieve users by author ID", "error", err)
		return nil, err
	}

	return users, nil
}

func ReadUserById(ctx context.Context, id uuid.UUID, logger *zap.SugaredLogger) (*entity.User, error) {
	logger.Infof("Retrieving user with ID: %s", id)
	user := new(entity.User)
	err := db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		logger.Errorw("Failed to retrieve user by ID", "error", err)
		return nil, err
	}
	return user, nil
}

func UpdateUser(ctx context.Context, user *entity.User, logger *zap.SugaredLogger) error {
	logger.Infof("Updating user with ID: %s", user.Id)
	_, err := db.NewUpdate().Model(user).Where("id = ?", user.Id).Exec(ctx)
	if err != nil {
		logger.Errorw("Failed to update user", "error", err)
		return err
	}
	return nil
}

func DeleteUser(ctx context.Context, id uuid.UUID, logger *zap.SugaredLogger) error {
	logger.Infof("Deleting %s", id)
	_, err := db.NewDelete().Model(&entity.User{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		logger.Errorw("Failed to delete user", "error", err)
		return err
	}
	return nil
}
