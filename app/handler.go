package app

import (
	"context"
	"encoding/json"
	"log"

	"userservice/entity"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type ActionFunc func(ctx context.Context, data []byte) (interface{}, error)

func handleMessage(action ActionFunc) nats.MsgHandler {
	return func(msg *nats.Msg) {
		result, err := action(context.Background(), msg.Data)
		if err != nil {
			log.Println("Error:", err)
			msg.Respond([]byte("Error processing request"))
			return
		}

		response, err := json.Marshal(result)
		if err != nil {
			log.Println("Error marshalling response:", err)
			msg.Respond([]byte("Error marshalling response"))
			return
		}

		msg.Respond(response)
	}
}

// CreateUserHandler handles the CreateUser messages
func CreateUserHandler() nats.MsgHandler {
	return handleMessage(func(ctx context.Context, data []byte) (interface{}, error) {
		var user entity.User
		if err := json.Unmarshal(data, &user); err != nil {
			return nil, err
		}
		return CreateUser(ctx, user)
	})
}

// UpdateUserHandler handles the UpdateUser messages
func UpdateUserHandler() nats.MsgHandler {
	return handleMessage(func(ctx context.Context, data []byte) (interface{}, error) {
		var user entity.User
		if err := json.Unmarshal(data, &user); err != nil {
			return nil, err
		}
		return UpdateUser(ctx, user)
	})
}

// ReadUserHandler handles the ReadUser messages
func ReadUserHandler() nats.MsgHandler {
	return handleMessage(func(ctx context.Context, data []byte) (interface{}, error) {
		id, err := uuid.Parse(string(data))
		if err != nil {
			return nil, err
		}
		return ReadUser(ctx, id)
	})
}

// ReadAllUsersHandler handles the ReadAllUsers messages
func ReadAllUsersHandler() nats.MsgHandler {
	return handleMessage(func(ctx context.Context, data []byte) (interface{}, error) {
		id := string(data)
		return ReadAllUsers(ctx, id)
	})
}

// DeleteUserHandler handles the DeleteUser messages
func DeleteUserHandler() nats.MsgHandler {
	return handleMessage(func(ctx context.Context, data []byte) (interface{}, error) {
		id, err := uuid.Parse(string(data))
		if err != nil {
			return nil, err
		}
		return DeleteUser(ctx, id)
	})
}
