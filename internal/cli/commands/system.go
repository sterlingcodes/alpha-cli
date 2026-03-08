package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/system/battery"
	"github.com/sterlingcodes/alpha-cli/internal/system/calendar"
	"github.com/sterlingcodes/alpha-cli/internal/system/cleanup"
	"github.com/sterlingcodes/alpha-cli/internal/system/clipboard"
	"github.com/sterlingcodes/alpha-cli/internal/system/contacts"
	"github.com/sterlingcodes/alpha-cli/internal/system/diskhealth"
	"github.com/sterlingcodes/alpha-cli/internal/system/finder"
	"github.com/sterlingcodes/alpha-cli/internal/system/imessage"
	"github.com/sterlingcodes/alpha-cli/internal/system/mail"
	"github.com/sterlingcodes/alpha-cli/internal/system/notes"
	"github.com/sterlingcodes/alpha-cli/internal/system/reminders"
	"github.com/sterlingcodes/alpha-cli/internal/system/safari"
	"github.com/sterlingcodes/alpha-cli/internal/system/sysinfo"
)

func NewSystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "system",
		Aliases: []string{"sys"},
		Short:   "System commands",
		Long:    `System-level integrations: Apple Notes, Calendar, Reminders, Contacts, Finder, Safari, Mail, Clipboard, iMessage (macOS only).`,
	}

	cmd.AddCommand(calendar.NewCmd())
	cmd.AddCommand(clipboard.NewCmd())
	cmd.AddCommand(contacts.NewCmd())
	cmd.AddCommand(finder.NewCmd())
	cmd.AddCommand(imessage.NewCmd())
	cmd.AddCommand(mail.NewCmd())
	cmd.AddCommand(notes.NewCmd())
	cmd.AddCommand(reminders.NewCmd())
	cmd.AddCommand(safari.NewCmd())
	cmd.AddCommand(sysinfo.NewCmd())
	cmd.AddCommand(battery.NewCmd())
	cmd.AddCommand(diskhealth.NewCmd())
	cmd.AddCommand(cleanup.NewCmd())

	return cmd
}
