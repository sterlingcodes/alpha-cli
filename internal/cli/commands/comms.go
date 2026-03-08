package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/communication/discord"
	"github.com/sterlingcodes/alpha-cli/internal/communication/email"
	"github.com/sterlingcodes/alpha-cli/internal/communication/notify"
	"github.com/sterlingcodes/alpha-cli/internal/communication/slack"
	"github.com/sterlingcodes/alpha-cli/internal/communication/telegram"
	"github.com/sterlingcodes/alpha-cli/internal/communication/twilio"
	"github.com/sterlingcodes/alpha-cli/internal/communication/webhook"
)

func NewCommsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "comms",
		Aliases: []string{"c", "comm"},
		Short:   "Communication commands",
		Long:    `Interact with communication platforms: Email, Slack, Discord, Telegram, etc.`,
	}

	cmd.AddCommand(email.NewCmd())
	cmd.AddCommand(slack.NewCmd())
	cmd.AddCommand(discord.NewCmd())
	cmd.AddCommand(telegram.NewCmd())
	cmd.AddCommand(twilio.NewCmd())
	cmd.AddCommand(webhook.NewCmd())
	cmd.AddCommand(notify.NewCmd())

	return cmd
}
