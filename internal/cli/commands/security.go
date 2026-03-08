package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/security/crtsh"
	"github.com/sterlingcodes/alpha-cli/internal/security/hibp"
	"github.com/sterlingcodes/alpha-cli/internal/security/shodan"
	"github.com/sterlingcodes/alpha-cli/internal/security/virustotal"
)

func NewSecurityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security",
		Aliases: []string{"sec"},
		Short:   "Security commands",
		Long:    "Security tools: VirusTotal, Shodan, Certificate Transparency, Have I Been Pwned.",
	}

	cmd.AddCommand(virustotal.NewCmd())
	cmd.AddCommand(shodan.NewCmd())
	cmd.AddCommand(crtsh.NewCmd())
	cmd.AddCommand(hibp.NewCmd())

	return cmd
}
