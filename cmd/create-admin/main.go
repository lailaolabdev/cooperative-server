package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"cooperative-service/internal/config"
	"cooperative-service/internal/database"
	authmodule "cooperative-service/internal/modules/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	usernameFlag := flag.String("username", "", "admin username")
	passwordFlag := flag.String("password", "", "admin password (local seeding only)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	username := strings.TrimSpace(*usernameFlag)
	password := []byte(*passwordFlag)
	if username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Admin username: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
	}
	if len(password) == 0 {
		fmt.Print("Admin password: ")
		password, err = term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			log.Fatal(err)
		}
	}
	if username == "" || len(password) < 8 {
		log.Fatal("username is required and password must contain at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now().UTC()
	repository := authmodule.NewRepository(client.Database(cfg.Database))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = repository.EnsureIndexes(ctx); err != nil {
		log.Fatal(err)
	}
	err = repository.Create(ctx, authmodule.Admin{ID: primitive.NewObjectID(), Username: username, PasswordHash: string(hash), Active: true, CreatedAt: now, UpdatedAt: now})
	if err != nil {
		log.Fatalf("cannot create admin (username may already exist): %v", err)
	}
	fmt.Printf("Admin %q created successfully.\n", username)
}
