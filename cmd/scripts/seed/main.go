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
		{
			Email:        "client@example.com",
			PasswordHash: string(hashpassword),
			FirstName:    "Sarah",
			LastName:     "Johnson",
			Role:         "client",
			Suspended:    false,
			EmailVerifiedAt: sql.NullTime{Time: time.Now(),
				Valid: true},
		},
		{
			Email:        "foued@hamrouni.com",
			PasswordHash: string(hashpassword),
			FirstName:    "Foued",
			LastName:     "Hamrouni",
			Role:         "client",
			Suspended:    false,
			EmailVerifiedAt: sql.NullTime{Time: time.Now(),
				Valid: true},
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

	// Seed worker profiles with sample data
	workerProfiles := []struct {
		Email      string
		Name       string
		Profession string
		Rating     float64
		Distance   float64
		Reviews    int
		Price      float64
		Available  bool
		Latitude   float64
		Longitude  float64
	}{
		{"mohamed.hassan@example.com", "Mohamed Hassan", "Plumber", 4.9, 0.8, 127, 180, true, 36.8065, 10.1815},
		{"ahmed.elsayed@example.com", "Ahmed El-Sayed", "Electrician", 4.8, 1.2, 98, 200, true, 36.8115, 10.1915},
		{"hossam.abid@example.com", "Hossam Abid", "Painter", 4.5, 2.1, 55, 150, false, 36.8265, 10.2015},
		{"youssef.mansour@example.com", "Youssef Mansour", "AC Technician", 4.7, 1.5, 89, 190, true, 36.8165, 10.1965},
		{"karim.saidi@example.com", "Karim Saidi", "Cleaner", 4.6, 0.5, 112, 120, true, 36.8015, 10.1765},
		{"ali.bensalem@example.com", "Ali Ben Salem", "Plumber", 4.8, 3.2, 145, 170, true, 36.8365, 10.2115},
		{"sofiane.gharbi@example.com", "Sofiane Gharbi", "Electrician", 4.9, 2.8, 203, 210, false, 36.8315, 10.2065},
		{"mehdi.trabelsi@example.com", "Mehdi Trabelsi", "Carpenter", 4.7, 1.9, 78, 165, true, 36.8215, 10.1965},
		{"nabil.chebbi@example.com", "Nabil Chebbi", "Locksmith", 4.6, 2.5, 92, 140, true, 36.8265, 10.2015},
		{"rami.bouazizi@example.com", "Rami Bouazizi", "Painter", 4.5, 3.0, 67, 155, true, 36.8315, 10.2065},
		{"farid.jelassi@example.com", "Farid Jelassi", "AC Technician", 4.8, 1.8, 134, 195, true, 36.8165, 10.1915},
		{"tarek.maatoug@example.com", "Tarek Maatoug", "Plumber", 4.7, 2.2, 101, 175, false, 36.8215, 10.1965},
		{"walid.hamdi@example.com", "Walid Hamdi", "Electrician", 4.9, 0.9, 187, 205, true, 36.8065, 10.1865},
		{"sami.ayari@example.com", "Sami Ayari", "Cleaner", 4.6, 1.7, 88, 125, true, 36.8165, 10.1915},
		{"bassem.jribi@example.com", "Bassem Jribi", "Carpenter", 4.5, 3.5, 72, 160, true, 36.8415, 10.2165},
		// Additional 15 workers
		{"omar.khelifi@example.com", "Omar Khelifi", "Plumber", 4.8, 4.0, 156, 185, true, 36.8465, 10.2215},
		{"amine.ben.ali@example.com", "Amine Ben Ali", "Electrician", 4.7, 3.8, 142, 195, false, 36.8415, 10.2165},
		{"khaled.mansouri@example.com", "Khaled Mansouri", "Carpenter", 4.9, 2.3, 189, 170, true, 36.8215, 10.1965},
		{"bilal.gharbi@example.com", "Bilal Gharbi", "AC Technician", 4.6, 4.5, 98, 185, true, 36.8515, 10.2265},
		{"hamza.jemni@example.com", "Hamza Jemni", "Painter", 4.8, 1.3, 123, 155, true, 36.8115, 10.1865},
		{"adel.chebbi@example.com", "Adel Chebbi", "Locksmith", 4.5, 3.7, 76, 145, true, 36.8415, 10.2165},
		{"slim.meddeb@example.com", "Slim Meddeb", "Cleaner", 4.7, 2.6, 104, 130, false, 36.8265, 10.2015},
		{"fares.azizi@example.com", "Fares Azizi", "Plumber", 4.9, 1.1, 198, 190, true, 36.8085, 10.1865},
		{"houssem.triki@example.com", "Houssem Triki", "Electrician", 4.6, 4.2, 87, 200, true, 36.8465, 10.2215},
		{"mourad.dali@example.com", "Mourad Dali", "Carpenter", 4.8, 2.9, 167, 175, true, 36.8315, 10.2065},
		{"zied.souissi@example.com", "Zied Souissi", "AC Technician", 4.7, 3.3, 111, 180, false, 36.8365, 10.2115},
		{"wassim.ghariani@example.com", "Wassim Ghariani", "Painter", 4.5, 4.8, 64, 160, true, 36.8565, 10.2315},
		{"samir.chaabane@example.com", "Samir Chaabane", "Locksmith", 4.9, 0.6, 201, 150, true, 36.8025, 10.1775},
		{"aymen.messaoudi@example.com", "Aymen Messaoudi", "Cleaner", 4.6, 3.6, 91, 125, true, 36.8405, 10.2155},
		{"chokri.lahmar@example.com", "Chokri Lahmar", "Plumber", 4.8, 2.0, 145, 175, true, 36.8215, 10.1965},
	}

	for _, wp := range workerProfiles {
		// First, create or get the user
		var user models.User
		result := dbConn.Where("email = ?", wp.Email).First(&user)
		if result.Error != nil {
			// Create the user if doesn't exist
			names := []string{wp.Name}
			if len(names) > 0 {
				parts := []rune(wp.Name)
				firstName := string(parts[:len(parts)/2])
				lastName := string(parts[len(parts)/2:])

				user = models.User{
					Email:        wp.Email,
					PasswordHash: string(hashpassword),
					FirstName:    firstName,
					LastName:     lastName,
					Role:         "worker",
					Suspended:    false,
					Rating:       wp.Rating,
				}
				if err := dbConn.Create(&user).Error; err != nil {
					fmt.Printf("Failed to insert worker user %s: %v\n", wp.Email, err)
					continue
				}
				fmt.Printf("Inserted worker user: %s\n", wp.Email)
			}
		}

		// Now create the worker profile
		var existingProfile models.WorkerProfile
		result = dbConn.Where("user_id = ?", user.ID).First(&existingProfile)
		if result.Error == nil {
			fmt.Printf("Worker profile for %s already exists, skipping\n", wp.Name)
			continue
		}

		profile := models.WorkerProfile{
			UserID:     user.ID,
			Name:       wp.Name,
			Profession: wp.Profession,
			Rating:     wp.Rating,
			Distance:   wp.Distance,
			Reviews:    wp.Reviews,
			Price:      wp.Price,
			Available:  wp.Available,
			Latitude:   wp.Latitude,
			Longitude:  wp.Longitude,
		}

		if err := dbConn.Create(&profile).Error; err != nil {
			fmt.Printf("Failed to insert worker profile %s: %v\n", wp.Name, err)
		} else {
			fmt.Printf("Inserted worker profile: %s\n", wp.Name)
		}
	}
}
