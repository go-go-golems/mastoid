package pkg

import (
	"context"
	"fmt"
	"github.com/mattn/go-mastodon"
	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
	"github.com/yardbirdsax/bubblewrap"
)

func Authorize(ctx context.Context, credentials_ *Credentials) error {
	// open app.AuthURI in browser
	err := browser.OpenURL(credentials_.Application.AuthURI)
	if err != nil {
		return fmt.Errorf("Error opening browser: %w", err)
	}

	isCodeValid := false

	var grantToken string

	for !isCodeValid {
		grantToken, err = bubblewrap.Input("Enter the code from the browser: ")
		if err != nil {
			return fmt.Errorf("Error reading user input: %w", err)
		}

		credentials_.GrantToken = grantToken
		log.Debug().Str("GrantToken", credentials_.GrantToken).Msg("Grant Token")

		client := mastodon.NewClient(&mastodon.Config{
			Server:       credentials_.Server,
			ClientID:     credentials_.Application.ClientID,
			ClientSecret: credentials_.Application.ClientSecret,
		})

		err = client.AuthenticateApp(ctx)
		if err != nil {
			fmt.Printf("Error authenticating app: %s\n", err)
			isCodeValid = false
			continue
		}

		err = client.AuthenticateToken(ctx, grantToken, credentials_.Application.RedirectURI)
		if err != nil {
			fmt.Printf("Error authenticating token: %s\n", err)
			isCodeValid = false
			continue
		}
		credentials_.AccessToken = client.Config.AccessToken
		fmt.Printf("Access Token: %s\n", credentials_.AccessToken)

		credentials, err := client.VerifyAppCredentials(ctx)
		if err != nil {
			fmt.Printf("Error verifying credentials: %s\n", err)
			isCodeValid = false
			continue
		}

		if credentials.Website != "" {
			log.Info().Str("Website", credentials.Website).Msg("Website")
		}
		if credentials.Name != "" {
			log.Info().Str("Name", credentials.Name).Msg("Name")
		}

		isCodeValid = true
	}

	return nil
}
