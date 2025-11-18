package handlers

import (
	"skillKonnect/app/db"
	"skillKonnect/app/models"
	"skillKonnect/app/views/components/admin"

	"github.com/anthdm/superkit/kit"
)

func AdminDashboard(kit *kit.Kit) error {
	return kit.Render(admin.Dashboard())
}

// @Summary Dashboard stats
// @Description Get admin dashboard statistics
// @Tags admin, dashboard
// @Produce json
// @Success 200 {object} map[string]int64
// @Router /api/admin/stats/dashboard [get]
func AdminDashboardStats(kit *kit.Kit) error {
	var users, suspendedUsers, workers, pendingVerifications, skills, categories int64
	dbConn := db.Get()
	dbConn.Model(&models.User{}).Count(&users)
	dbConn.Model(&models.User{}).Where("suspended = ?", true).Count(&suspendedUsers)
	dbConn.Model(&models.User{}).Where("role = ?", "worker").Count(&workers)
	dbConn.Model(&models.User{}).Where("email_verified_at IS NULL").Count(&pendingVerifications)
	dbConn.Model(&models.Skill{}).Count(&skills)
	dbConn.Model(&models.Category{}).Count(&categories)

	return kit.JSON(200, map[string]int64{
		"users":                 users,
		"suspended_users":       suspendedUsers,
		"workers":               workers,
		"pending_verifications": pendingVerifications,
		"skills":                skills,
		"categories":            categories,
	})
}
