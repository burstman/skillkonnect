package main

import (
	"database/sql"
	"fmt"
	"skillKonnect/app/db"
	"skillKonnect/plugins/auth"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbConn := db.Get()

	// get your GORM instance

	password := "adminpassword"

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(fmt.Sprintf("failed to hash password: %v", err))
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		panic("failed to get database handle")
	}
	defer sqlDB.Close()

	// Seed users
	users := []auth.User{
		{
			Email:        "admin@example.com",
			PasswordHash: string(hashpassword), // replace with bcrypt hash in real app
			FirstName:    "Admin",
			LastName:     "User",
			Role:         "admin",
			Suspended:    false,
			EmailVerifiedAt: sql.NullTime{Time: time.Now(),
				Valid: true},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:        "john@example.com",
			PasswordHash: string(hashpassword),
			FirstName:    "John",
			LastName:     "Doe",
			Role:         "user",
			Suspended:    false,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Email:        "jane@example.com",
			PasswordHash: string(hashpassword),
			FirstName:    "Jane",
			LastName:     "Smith",
			Role:         "user",
			Suspended:    true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, u := range users {
		// Use FirstOrCreate to avoid duplicates
		var existing auth.User
		result := dbConn.Where("email = ?", u.Email).First(&existing)
		if result.Error == nil {
			fmt.Printf("User %s already exists, skipping\n", u.Email)
			continue
		}
		if err := dbConn.Create(&u).Error; err != nil {
			fmt.Printf("Failed to insert user %s: %v\n", u.Email, err)
		} else {
			fmt.Printf("Inserted user: %s\n", u.Email)
		}
	}
}
