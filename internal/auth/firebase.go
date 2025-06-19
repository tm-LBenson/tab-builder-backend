package auth

import (
	"context"
	"errors"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var fbClient *auth.Client

func client(ctx context.Context) (*auth.Client, error) {
	if fbClient != nil { return fbClient, nil }
	app, err := firebase.NewApp(ctx, nil,
		option.WithCredentialsFile("/secrets/firebase.json"))
	if err != nil { return nil, err }
	fbClient, err = app.Auth(ctx)
	return fbClient, err
}

func Verify(ctx context.Context, idToken string) (string, error) {
	c, err := client(ctx)
	if err != nil { return "", err }
	tok, err := c.VerifyIDToken(ctx, idToken)
	if err != nil { return "", err }
	uid, ok := tok.Claims["user_id"].(string)
	if !ok { return "", errors.New("uid claim missing") }
	return uid, nil
}
