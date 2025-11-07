package handlers

import (
	"skillKonnect/app/views/components/admin"

	"github.com/anthdm/superkit/kit"
)

func AdminDashboard(kit *kit.Kit) error {
	return kit.Render(admin.Dashboard())
}
