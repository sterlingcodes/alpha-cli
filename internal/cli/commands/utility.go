package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/utility/crypto"
	"github.com/sterlingcodes/alpha-cli/internal/utility/currency"
	"github.com/sterlingcodes/alpha-cli/internal/utility/dnsbench"
	"github.com/sterlingcodes/alpha-cli/internal/utility/domain"
	"github.com/sterlingcodes/alpha-cli/internal/utility/geocoding"
	"github.com/sterlingcodes/alpha-cli/internal/utility/holidays"
	"github.com/sterlingcodes/alpha-cli/internal/utility/ipinfo"
	"github.com/sterlingcodes/alpha-cli/internal/utility/netdiag"
	"github.com/sterlingcodes/alpha-cli/internal/utility/paste"
	"github.com/sterlingcodes/alpha-cli/internal/utility/speedtest"
	"github.com/sterlingcodes/alpha-cli/internal/utility/stocks"
	"github.com/sterlingcodes/alpha-cli/internal/utility/timezone"
	"github.com/sterlingcodes/alpha-cli/internal/utility/traceroute"
	"github.com/sterlingcodes/alpha-cli/internal/utility/translate"
	"github.com/sterlingcodes/alpha-cli/internal/utility/urlshort"
	"github.com/sterlingcodes/alpha-cli/internal/utility/wayback"
	"github.com/sterlingcodes/alpha-cli/internal/utility/weather"
	"github.com/sterlingcodes/alpha-cli/internal/utility/wifi"
)

func NewUtilityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "utility",
		Aliases: []string{"u", "util"},
		Short:   "Utility commands",
		Long:    `Utility tools: weather, crypto, stocks, currency, DNS/WHOIS, translation, etc.`,
	}

	cmd.AddCommand(weather.NewCmd())
	cmd.AddCommand(crypto.NewCmd())
	cmd.AddCommand(ipinfo.NewCmd())
	cmd.AddCommand(domain.NewCmd())
	cmd.AddCommand(currency.NewCmd())
	cmd.AddCommand(wayback.NewCmd())
	cmd.AddCommand(holidays.NewCmd())
	cmd.AddCommand(translate.NewCmd())
	cmd.AddCommand(stocks.NewCmd())
	cmd.AddCommand(urlshort.NewCmd())
	cmd.AddCommand(geocoding.NewCmd())
	cmd.AddCommand(netdiag.NewCmd())
	cmd.AddCommand(paste.NewCmd())
	cmd.AddCommand(timezone.NewCmd())
	cmd.AddCommand(speedtest.NewCmd())
	cmd.AddCommand(dnsbench.NewCmd())
	cmd.AddCommand(traceroute.NewCmd())
	cmd.AddCommand(wifi.NewCmd())

	return cmd
}
