package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/knowledge/dictionary"
	"github.com/sterlingcodes/alpha-cli/internal/knowledge/stackexchange"
	"github.com/sterlingcodes/alpha-cli/internal/knowledge/wikipedia"
)

func NewKnowledgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "knowledge",
		Aliases: []string{"k", "know"},
		Short:   "Knowledge and research commands",
		Long:    `Access knowledge bases: Wikipedia, dictionaries, etc.`,
	}

	cmd.AddCommand(wikipedia.NewCmd())
	cmd.AddCommand(stackexchange.NewCmd())
	cmd.AddCommand(dictionary.NewCmd())

	return cmd
}
