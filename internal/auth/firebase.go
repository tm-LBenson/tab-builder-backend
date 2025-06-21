package auth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var fbClient *auth.Client

func client(ctx context.Context) (*auth.Client, error) {
	if fbClient != nil {
		return fbClient, nil
	}
	app, err := firebase.NewApp(ctx, nil,
		option.WithCredentialsFile("/secrets/firebase.json"))
	if err != nil {
		return nil, err
	}
	fbClient, err = app.Auth(ctx)
	return fbClient, err
}
func Verify(ctx context.Context, idToken string) (uid, name string, err error) {
	c, err := client(ctx)
	if err != nil {
		return
	}

	tok, err := c.VerifyIDToken(ctx, idToken)
	if err != nil {
		return
	}

	uid = tok.UID
	name = tok.Claims["name"].(string)
	return
}
