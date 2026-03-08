package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/dev/cloudflare"
	"github.com/sterlingcodes/alpha-cli/internal/dev/database"
	"github.com/sterlingcodes/alpha-cli/internal/dev/dockerhub"
	"github.com/sterlingcodes/alpha-cli/internal/dev/gist"
	"github.com/sterlingcodes/alpha-cli/internal/dev/github"
	"github.com/sterlingcodes/alpha-cli/internal/dev/gitlab"
	"github.com/sterlingcodes/alpha-cli/internal/dev/jira"
	"github.com/sterlingcodes/alpha-cli/internal/dev/kubernetes"
	"github.com/sterlingcodes/alpha-cli/internal/dev/linear"
	"github.com/sterlingcodes/alpha-cli/internal/dev/npm"
	"github.com/sterlingcodes/alpha-cli/internal/dev/prometheus"
	"github.com/sterlingcodes/alpha-cli/internal/dev/pypi"
	"github.com/sterlingcodes/alpha-cli/internal/dev/redis"
	"github.com/sterlingcodes/alpha-cli/internal/dev/s3"
	"github.com/sterlingcodes/alpha-cli/internal/dev/sentry"
	"github.com/sterlingcodes/alpha-cli/internal/dev/vercel"
)

func NewDevCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dev",
		Aliases: []string{"d"},
		Short:   "Developer tool commands",
		Long:    `Interact with developer tools: GitHub, GitLab, Jira, Cloudflare, Vercel, Docker Hub, etc.`,
	}

	cmd.AddCommand(github.NewCmd())
	cmd.AddCommand(gitlab.NewCmd())
	cmd.AddCommand(linear.NewCmd())
	cmd.AddCommand(npm.NewCmd())
	cmd.AddCommand(pypi.NewCmd())
	cmd.AddCommand(jira.NewCmd())
	cmd.AddCommand(cloudflare.NewCmd())
	cmd.AddCommand(vercel.NewCmd())
	cmd.AddCommand(dockerhub.NewCmd())
	cmd.AddCommand(sentry.NewCmd())
	cmd.AddCommand(redis.NewCmd())
	cmd.AddCommand(prometheus.NewCmd())
	cmd.AddCommand(kubernetes.NewCmd())
	cmd.AddCommand(database.NewCmd())
	cmd.AddCommand(s3.NewCmd())
	cmd.AddCommand(gist.NewCmd())

	return cmd
}
