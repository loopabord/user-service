package database

import (
	"context"
	"userservice/entity"

	"go.uber.org/zap"
)

func CreateUser(ctx context.Context, user *entity.User, logger *zap.SugaredLogger) error {
	logger.Infof("Creating user with email: %s", user.Email)

	// Encrypt email and name before storing
	encryptedEmail, err := Encrypt(user.Email)
	if err != nil {
		logger.Errorw("Failed to encrypt email", "error", err)
		return err
	}
	encryptedName, err := Encrypt(user.Name)
	if err != nil {
		logger.Errorw("Failed to encrypt name", "error", err)
		return err
	}

	// Replace user's email and name with encrypted versions
	user.Email = encryptedEmail
	user.Name = encryptedName

	var id string
	_, err = db.NewInsert().Model(user).ExcludeColumn("created_at").Returning("id", "created_at").Exec(ctx, &id)
	if err != nil {
		logger.Errorw("Failed to insert user", "error", err)
		return err
	}

	// Update the user object with the generated ID
	user.Id = id
	return nil
}

func ReadAll(ctx context.Context, authorId string, logger *zap.SugaredLogger) ([]entity.User, error) {
	logger.Info("Retrieving all users")

	var users []entity.User
	err := db.NewSelect().Model(&users).Order("created_at DESC").Scan(ctx)
	if err != nil {
		logger.Errorw("Failed to retrieve users", "error", err)
		return nil, err
	}

	// Decrypt email and name for all retrieved users
	for i := range users {
		decryptedEmail, err := Decrypt(users[i].Email)
		if err != nil {
			logger.Errorw("Failed to decrypt email", "error", err)
			return nil, err
		}
		decryptedName, err := Decrypt(users[i].Name)
		if err != nil {
			logger.Errorw("Failed to decrypt name", "error", err)
			return nil, err
		}
		users[i].Email = decryptedEmail
		users[i].Name = decryptedName
	}

	return users, nil
}

func ReadUserById(ctx context.Context, id string, logger *zap.SugaredLogger) (*entity.User, error) {
	logger.Infof("Retrieving user with ID: %s", id)
	user := new(entity.User)
	err := db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		logger.Errorw("Failed to retrieve user by ID", "error", err)
		return nil, err
	}

	// Decrypt email and name after retrieving
	decryptedEmail, err := Decrypt(user.Email)
	if err != nil {
		logger.Errorw("Failed to decrypt email", "error", err)
		return nil, err
	}
	decryptedName, err := Decrypt(user.Name)
	if err != nil {
		logger.Errorw("Failed to decrypt name", "error", err)
		return nil, err
	}

	// Replace encrypted email and name with decrypted versions
	user.Email = decryptedEmail
	user.Name = decryptedName

	return user, nil
}

func UpdateUser(ctx context.Context, user *entity.User, logger *zap.SugaredLogger) error {
	logger.Infof("Updating user with ID: %s", user.Id)

	// Encrypt email and name before updating
	encryptedEmail, err := Encrypt(user.Email)
	if err != nil {
		logger.Errorw("Failed to encrypt email", "error", err)
		return err
	}
	encryptedName, err := Encrypt(user.Name)
	if err != nil {
		logger.Errorw("Failed to encrypt name", "error", err)
		return err
	}

	// Replace user's email and name with encrypted versions
	user.Email = encryptedEmail
	user.Name = encryptedName

	_, err = db.NewUpdate().Model(user).Where("id = ?", user.Id).Exec(ctx)
	if err != nil {
		logger.Errorw("Failed to update user", "error", err)
		return err
	}
	return nil
}

func DeleteUser(ctx context.Context, id string, logger *zap.SugaredLogger) error {
	logger.Infof("Deleting user with ID: %s", id)
	_, err := db.NewDelete().Model(&entity.User{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		logger.Errorw("Failed to delete user", "error", err)
		return err
	}
	return nil
}
