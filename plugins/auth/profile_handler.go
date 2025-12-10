package auth

import (
	"net/http"
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

	// Read authentication from unified middleware
	payload, ok := kit.Request.Context().Value(models.AuthContextKey{}).(models.AuthPayload)
	if !ok || !payload.Authenticated || payload.User == nil {
		return kit.Redirect(http.StatusSeeOther, "/web/admin/login")
	}

	user := payload.User

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

	payload, ok := kit.Request.Context().Value(models.AuthContextKey{}).(models.AuthPayload)
	if !ok || !payload.Authenticated || payload.User == nil {
		return kit.Redirect(http.StatusSeeOther, "/web/admin/login")
	}
	err := db.Get().Model(&models.User{}).
		Where("id = ?", payload.User.ID).
		Updates(&models.User{
			FirstName: values.FirstName,
			LastName:  values.LastName,
		}).Error
	if err != nil {
		return err
	}

	values.Success = "Profile successfully updated!"
	values.Email = payload.User.Email

	return kit.Render(ProfileForm(values, v.Errors{}))
}
