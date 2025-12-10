package handlers

import (
	"net/http"
	"skillKonnect/app/db"
	"skillKonnect/app/models"
	"strconv"
	"strings"

	"github.com/anthdm/superkit/kit"
)

// WorkerProfileResponse is the response structure for worker profiles
type WorkerProfileResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Profession string  `json:"profession"`
	Rating     float64 `json:"rating"`
	Distance   string  `json:"distance"`
	Reviews    int     `json:"reviews"`
	Price      string  `json:"price"`
	Available  bool    `json:"available"`
}

// @Summary List workers by distance
// @Description Get list of workers sorted by distance for clients
// @Tags client, workers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category query string false "Filter by profession/category (e.g., Plumber, Electrician, Carpenter)"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/client/workers [get]
func ClientListWorkers(kit *kit.Kit) error {
	// Get pagination params
	page, _ := strconv.Atoi(kit.Request.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(kit.Request.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get category filter
	category := kit.Request.URL.Query().Get("category")

	var workers []models.WorkerProfile
	var total int64

	dbConn := db.Get()

	// Build query
	query := dbConn.Model(&models.WorkerProfile{})

	// Apply category filter if provided (case-insensitive)
	if category != "" {
		query = query.Where("LOWER(profession) = ?", strings.ToLower(category))
	}

	// Count total workers with filter
	query.Count(&total)

	// Get workers sorted by distance
	if err := query.
		Order("distance ASC").
		Limit(limit).
		Offset(offset).
		Find(&workers).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch workers",
		})
	}

	// Format the response
	response := make([]WorkerProfileResponse, len(workers))
	for i, w := range workers {
		response[i] = WorkerProfileResponse{
			ID:         w.ID,
			Name:       w.Name,
			Profession: w.Profession,
			Rating:     w.Rating,
			Distance:   strconv.FormatFloat(w.Distance, 'f', 1, 64) + " km",
			Reviews:    w.Reviews,
			Price:      strconv.FormatFloat(w.Price, 'f', 0, 64) + " TND/hr",
			Available:  w.Available,
		}
	}

	return kit.JSON(http.StatusOK, map[string]interface{}{
		"workers": response,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}
