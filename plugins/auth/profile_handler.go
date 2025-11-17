package auth

import (
	"fmt"
	"skillKonnect/app/db"
	"skillKonnect/app/models"

	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
)

var profileSchema = v.Schema{
	"firstName": v.Rules(v.Min(3), v.Max(50)),
	"lastName":  v.Rules(v.Min(3), v.Max(50)),
}

type ProfileFormValues struct {
	ID        uint   `form:"id"`
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
	Email     string
	Success   string
}

func HandleProfileShow(kit *kit.Kit) error {
	auth := kit.Auth().(models.AuthPayload)

	var user models.User
	if err := db.Get().First(&user, auth.User.ID).Error; err != nil {
		return err
	}

	formValues := ProfileFormValues{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return kit.Render(ProfileShow(formValues))
}

func HandleProfileUpdate(kit *kit.Kit) error {
	var values ProfileFormValues
	errors, ok := v.Request(kit.Request, &values, profileSchema)
	if !ok {
		return kit.Render(ProfileForm(values, errors))
	}

	auth := kit.Auth().(models.AuthPayload)
	if auth.User.ID != values.ID {
		return fmt.Errorf("unauthorized request for profile %d", values.ID)
	}
	err := db.Get().Model(&models.User{}).
		Where("id = ?", auth.User.ID).
		Updates(&models.User{
			FirstName: values.FirstName,
			LastName:  values.LastName,
		}).Error
	if err != nil {
		return err
	}

	values.Success = "Profile successfully updated!"
	values.Email = auth.User.Email

	return kit.Render(ProfileForm(values, v.Errors{}))
}
