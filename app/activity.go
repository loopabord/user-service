package app

import (
	"context"
	"log"
	"userservice/database"
	"userservice/entity"

	"go.uber.org/zap"
)

func withDBAndLogger(ctx context.Context, fn func(ctx context.Context, logger *zap.SugaredLogger) (interface{}, error)) (interface{}, error) {
	// Establish a session with the database cluster
	err := database.Connect()
	if err != nil {
		return nil, err // Return an error if connection fails
	}
	defer database.Close() // Defer closing the session to ensure it's closed after this function returns
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	defer logger.Sync()

	return fn(ctx, sugar)
}

func CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	result, err := withDBAndLogger(ctx, func(ctx context.Context, logger *zap.SugaredLogger) (interface{}, error) {
		err := database.CreateUser(ctx, &user, logger)
		if err != nil {
			return entity.User{}, err // Return an error if insertion fails
		}
		return user, nil
	})
	if err != nil {
		return entity.User{}, err
	}
	var value entity.User
	value = result.(entity.User)
	log.Println(value.CreatedAt)
	return value, nil
}

func UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	result, err := withDBAndLogger(ctx, func(ctx context.Context, logger *zap.SugaredLogger) (interface{}, error) {
		err := database.UpdateUser(ctx, &user, logger)
		if err != nil {
			return entity.User{}, err // Return an error if update fails
		}
		return user, nil
	})
	if err != nil {
		return entity.User{}, err
	}
	return result.(entity.User), nil
}

func ReadUser(ctx context.Context, id string) (entity.User, error) {
	result, err := withDBAndLogger(ctx, func(ctx context.Context, logger *zap.SugaredLogger) (interface{}, error) {
		user, err := database.ReadUserById(ctx, id, logger)
		if err != nil {
			return entity.User{}, err // Return an error if retrieval fails
		}
		return user, nil
	})
	if err != nil {
		return entity.User{}, err
	}
	return *result.(*entity.User), nil
}

func ReadAllUsers(ctx context.Context, authorId string) ([]entity.User, error) {
	result, err := withDBAndLogger(ctx, func(ctx context.Context, logger *zap.SugaredLogger) (interface{}, error) {
		users, err := database.ReadAll(ctx, authorId, logger)
		if err != nil {
			return nil, err // Return an error if reading fails
		}
		return users, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]entity.User), nil
}

func DeleteUser(ctx context.Context, id string) (string, error) {
	result, err := withDBAndLogger(ctx, func(ctx context.Context, logger *zap.SugaredLogger) (interface{}, error) {
		err := database.DeleteUser(ctx, id, logger)
		if err != nil {
			return "no good (bad)", err // Return an error if deletion fails
		}
		return "good yes", nil
	})
	if err != nil {
		return "no good (bad)", err
	}
	return result.(string), nil
}
