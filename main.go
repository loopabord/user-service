package main

import (
	"context"
	"log"
	"os"
	"userservice/app"
	"userservice/database"

	"github.com/nats-io/nats.go"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()

	// Initialize the database if needed
	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize the database: ", err)
	}

	// Connect to the database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatal("Failed to close the database: ", err)
		}
	}()

	// Run migrations
	if err := database.RunMigrations(ctx); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	// Connect to NATS server
	natsURL := os.Getenv("NATS_URL")

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Register NATS subscribers for each user operation
	nc.Subscribe("CreateUser", app.CreateUserHandler())
	nc.Subscribe("UpdateUser", app.UpdateUserHandler())
	nc.Subscribe("ReadUser", app.ReadUserHandler())
	nc.Subscribe("ReadAllUsers", app.ReadAllUsersHandler())
	nc.Subscribe("DeleteUser", app.DeleteUserHandler())

	// Keep the connection alive
	select {}
}
