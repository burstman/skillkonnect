package handlers

import (
	"skillKonnect/app/models"

	"github.com/anthdm/superkit/kit"
)

func HandleAuthentication(kit *kit.Kit) (kit.Auth, error) {
	return models.AuthPayload{}, nil
}
