package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/social/mastodon"
	"github.com/sterlingcodes/alpha-cli/internal/social/reddit"
	"github.com/sterlingcodes/alpha-cli/internal/social/spotify"
	"github.com/sterlingcodes/alpha-cli/internal/social/twitter"
	"github.com/sterlingcodes/alpha-cli/internal/social/youtube"
)

func NewSocialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "social",
		Aliases: []string{"s"},
		Short:   "Social media commands",
		Long:    `Interact with social media platforms: Twitter/X, Reddit, Mastodon, YouTube, etc.`,
	}

	cmd.AddCommand(twitter.NewCmd())
	cmd.AddCommand(reddit.NewCmd())
	cmd.AddCommand(mastodon.NewCmd())
	cmd.AddCommand(youtube.NewCmd())
	cmd.AddCommand(spotify.NewCmd())

	return cmd
}
