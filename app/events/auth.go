package events

import (
	"context"
	"encoding/json"
	"fmt"
	"skillKonnect/plugins/auth"
)

// Event handlers
func OnUserSignup(ctx context.Context, event any) {
	userWithToken, ok := event.(auth.UserWithVerificationToken)
	if !ok {
		return
	}
	b, _ := json.MarshalIndent(userWithToken, "   ", "    ")
	fmt.Println(string(b))
}

func OnResendVerificationToken(ctx context.Context, event any) {
	userWithToken, ok := event.(auth.UserWithVerificationToken)
	if !ok {
		return
	}
	b, _ := json.MarshalIndent(userWithToken, "   ", "    ")
	fmt.Println(string(b))
}
