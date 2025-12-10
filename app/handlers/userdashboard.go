package handlers

import (
	"skillKonnect/app/db"
	"skillKonnect/app/models"

	"github.com/anthdm/superkit/kit"
)

// @Summary User dashboard stats
// @Description Get user dashboard statistics
// @Tags user, dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/user/dashboard [get]
func UserDashboardStats(kit *kit.Kit) error {
	// Get user from context
	authPayload := kit.Request.Context().Value(models.AuthContextKey{}).(models.AuthPayload)
	user := authPayload.User

	var skillsCount, applicationsCount int64
	dbConn := db.Get()

	// Count user's skills
	dbConn.Model(&models.Skill{}).Where("user_id = ?", user.ID).Count(&skillsCount)

	// You can add more user-specific stats here
	// For example: applications, messages, etc.

	return kit.JSON(200, map[string]interface{}{
		"user_id":      user.ID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"role":         user.Role,
		"skills_count": skillsCount,
		"applications": applicationsCount,
	})
}

// @Summary Client dashboard stats
// @Description Get client dashboard statistics
// @Tags client, dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/client/dashboard [get]
func ClientDashboardStats(kit *kit.Kit) error {
	// Get user from context
	authPayload := kit.Request.Context().Value(models.AuthContextKey{}).(models.AuthPayload)
	user := authPayload.User

	var activeJobs, completedJobs, totalSpent int64
	// dbConn := db.Get()

	// Add your client-specific queries here
	// For example: jobs posted, workers hired, etc.

	return kit.JSON(200, map[string]interface{}{
		"user_id":        user.ID,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"email":          user.Email,
		"role":           user.Role,
		"active_jobs":    activeJobs,
		"completed_jobs": completedJobs,
		"total_spent":    totalSpent,
	})
}

// @Summary Worker dashboard stats
// @Description Get worker dashboard statistics
// @Tags worker, dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/worker/dashboard [get]
func WorkerDashboardStats(kit *kit.Kit) error {
	// Get user from context
	authPayload := kit.Request.Context().Value(models.AuthContextKey{}).(models.AuthPayload)
	user := authPayload.User

	var skillsCount, activeGigs, completedGigs int64
	dbConn := db.Get()

	// Count worker's skills
	dbConn.Model(&models.Skill{}).Where("user_id = ?", user.ID).Count(&skillsCount)

	// Add your worker-specific queries here
	// For example: applications, earnings, ratings, etc.

	return kit.JSON(200, map[string]interface{}{
		"user_id":        user.ID,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"email":          user.Email,
		"role":           user.Role,
		"rating":         user.Rating,
		"bio":            user.Bio,
		"approved":       user.Approved,
		"skills_count":   skillsCount,
		"active_gigs":    activeGigs,
		"completed_gigs": completedGigs,
	})
}
