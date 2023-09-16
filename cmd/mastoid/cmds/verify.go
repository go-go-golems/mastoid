package cmds

import (
	"context"
	"fmt"
	"github.com/go-go-golems/mastoid/cmd/mastoid/pkg"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies the credentials of a Mastodon instance",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		credentials, err := pkg.LoadCredentials()
		cobra.CheckErr(err)

		if credentials.Application.ClientID == "" || credentials.Application.ClientSecret == "" {
			log.Error().Msg("No client ID or secret found")
			return
		}

		client, err := pkg.CreateClient(credentials)
		cobra.CheckErr(err)

		if client.Config.AccessToken == "" {
			fmt.Println("No access token found")
			if credentials.GrantToken != "" {
				log.Info().Msg("Authenticating with grant token")
				err = client.AuthenticateToken(ctx, credentials.GrantToken, credentials.Application.RedirectURI)
				cobra.CheckErr(err)

				log.Info().Msg("Grant token authenticated")
			} else {
				log.Info().Msg("Authenticating with app")
				err = client.AuthenticateApp(ctx)
				cobra.CheckErr(err)
				log.Info().Msg("App authenticated")
			}
		} else {
			log.Info().Str("AccessToken", client.Config.AccessToken).Msg("Access token found")
		}

		app, err := client.VerifyAppCredentials(ctx)
		cobra.CheckErr(err)

		log.Info().Str("AppName", app.Name).Msg("App Name")
		if app.Website != "" {
			log.Info().Str("AppWebsite", app.Website).Msg("App Website")
		}
	},
}
