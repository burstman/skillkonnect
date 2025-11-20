package main

import (
	"database/sql"
	"fmt"
	"skillKonnect/app/db"
	"skillKonnect/app/models"
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

	// Seed categories for skilled trades
	categories := []models.Category{
		{Name: "plumbing", Description: "Plumbing and pipe repair"},
		{Name: "electrical", Description: "Electrical installation and repair"},
		{Name: "carpentry", Description: "Carpentry and woodwork"},
		{Name: "hvac", Description: "Heating, ventilation, and air conditioning"},
		{Name: "painting", Description: "Painting and wall finishing"},
		{Name: "appliance-repair", Description: "Home appliance repair"},
	}
	for _, c := range categories {
		var existing models.Category
		result := dbConn.Where("name = ?", c.Name).First(&existing)
		if result.Error == nil {
			fmt.Printf("Category %s already exists, skipping\n", c.Name)
			continue
		}
		if err := dbConn.Create(&c).Error; err != nil {
			fmt.Printf("Failed to insert category %s: %v\n", c.Name, err)
		} else {
			fmt.Printf("Inserted category: %s\n", c.Name)
		}
	}

	// Seed skills for each category
	var plumbingCat, electricalCat, carpentryCat, hvacCat, paintingCat, applianceCat models.Category
	dbConn.Where("name = ?", "plumbing").First(&plumbingCat)
	dbConn.Where("name = ?", "electrical").First(&electricalCat)
	dbConn.Where("name = ?", "carpentry").First(&carpentryCat)
	dbConn.Where("name = ?", "hvac").First(&hvacCat)
	dbConn.Where("name = ?", "painting").First(&paintingCat)
	dbConn.Where("name = ?", "appliance-repair").First(&applianceCat)

	skills := []models.Skill{
		{Name: "pipe-installation", Description: "Install and repair pipes", CategoryID: plumbingCat.ID},
		{Name: "drain-cleaning", Description: "Clean and unclog drains", CategoryID: plumbingCat.ID},
		{Name: "leak-detection", Description: "Detect and fix leaks", CategoryID: plumbingCat.ID},

		{Name: "wiring", Description: "Electrical wiring and outlets", CategoryID: electricalCat.ID},
		{Name: "lighting-installation", Description: "Install lighting fixtures", CategoryID: electricalCat.ID},
		{Name: "circuit-breaker-repair", Description: "Repair circuit breakers", CategoryID: electricalCat.ID},

		{Name: "furniture-assembly", Description: "Assemble furniture", CategoryID: carpentryCat.ID},
		{Name: "door-installation", Description: "Install and repair doors", CategoryID: carpentryCat.ID},
		{Name: "cabinet-making", Description: "Build cabinets and shelves", CategoryID: carpentryCat.ID},

		{Name: "ac-installation", Description: "Install air conditioning units", CategoryID: hvacCat.ID},
		{Name: "heater-repair", Description: "Repair heating systems", CategoryID: hvacCat.ID},
		{Name: "ventilation-cleaning", Description: "Clean ventilation ducts", CategoryID: hvacCat.ID},

		{Name: "wall-painting", Description: "Paint interior and exterior walls", CategoryID: paintingCat.ID},
		{Name: "wall-prep", Description: "Prepare walls for painting", CategoryID: paintingCat.ID},
		{Name: "trim-painting", Description: "Paint trim and moldings", CategoryID: paintingCat.ID},

		{Name: "washer-repair", Description: "Repair washing machines", CategoryID: applianceCat.ID},
		{Name: "fridge-repair", Description: "Repair refrigerators", CategoryID: applianceCat.ID},
		{Name: "oven-repair", Description: "Repair ovens and stoves", CategoryID: applianceCat.ID},
	}
	for _, s := range skills {
		var existing models.Skill
		result := dbConn.Where("name = ?", s.Name).First(&existing)
		if result.Error == nil {
			fmt.Printf("Skill %s already exists, skipping\n", s.Name)
			continue
		}
		if err := dbConn.Create(&s).Error; err != nil {
			fmt.Printf("Failed to insert skill %s: %v\n", s.Name, err)
		} else {
			fmt.Printf("Inserted skill: %s\n", s.Name)
		}
	}

	// Seed users
	users := []models.User{
		{
			Email:        "admin@example.com",
			PasswordHash: string(hashpassword), // replace with bcrypt hash in real app
			FirstName:    "Admin",
			LastName:     "User",
			Role:         "admin",
			Suspended:    false,
			EmailVerifiedAt: sql.NullTime{Time: time.Now(),
				Valid: true},
		},
		{
			Email:        "john@example.com",
			PasswordHash: string(hashpassword),
			FirstName:    "John",
			LastName:     "Doe",
			Role:         "worker",
			Suspended:    false,
		},
		{
			Email:        "jane@example.com",
			PasswordHash: string(hashpassword),
			FirstName:    "Jane",
			LastName:     "Smith",
			Role:         "client",
			Suspended:    true,
		},
	}

	for _, u := range users {
		// Use FirstOrCreate to avoid duplicates
		var existing models.User
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
