package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/productivity/logseq"
	"github.com/sterlingcodes/alpha-cli/internal/productivity/notion"
	"github.com/sterlingcodes/alpha-cli/internal/productivity/obsidian"
	"github.com/sterlingcodes/alpha-cli/internal/productivity/todoist"
	"github.com/sterlingcodes/alpha-cli/internal/productivity/trello"
)

func NewProductivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "productivity",
		Aliases: []string{"p", "prod"},
		Short:   "Productivity tool commands",
		Long:    `Interact with productivity tools: Notion, Todoist, Trello, etc.`,
	}

	cmd.AddCommand(logseq.NewCmd())
	cmd.AddCommand(notion.NewCmd())
	cmd.AddCommand(obsidian.NewCmd())
	cmd.AddCommand(todoist.NewCmd())
	cmd.AddCommand(trello.NewCmd())

	return cmd
}
