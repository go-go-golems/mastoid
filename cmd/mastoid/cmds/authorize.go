package cmds

import (
	"context"
	"fmt"
	"github.com/go-go-golems/mastoid/cmd/mastoid/pkg"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var AuthorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Authorize the app with the Mastodon instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		credentials, err := pkg.LoadCredentials()
		if err != nil {
			return fmt.Errorf("Error loading credentials: %w", err)
		}

		log.Info().Msg("Authorizing app...")
		err = pkg.Authorize(context.Background(), credentials)
		if err != nil {
			log.Error().Err(err).Msgf("Error authorizing app")
			return fmt.Errorf("Error authorizing app: %w", err)
		}

		log.Info().Msg("App authorization successful!\n")

		err = pkg.StoreCredentials(credentials)
		if err != nil {
			return fmt.Errorf("Error storing credentials: %w", err)
		}
		log.Debug().Str("GrantToken", credentials.GrantToken).Msgf("Grant Token")
		log.Debug().Str("AccessToken", credentials.AccessToken).Msgf("Access Token")

		return nil
	},
}
