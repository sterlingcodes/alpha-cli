package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/pkg/output"
)

// Cmd represents a command for LLM consumption
type Cmd struct {
	Command string `json:"cmd"`
	Desc    string `json:"desc"`
	Args    string `json:"args,omitempty"`
	Flags   string `json:"flags,omitempty"`
}

// Group represents a command group
type Group struct {
	Name     string `json:"name"`
	Commands []Cmd  `json:"commands"`
}

func NewCommandsCmd() *cobra.Command {
	var group string

	cmd := &cobra.Command{
		Use:     "commands",
		Aliases: []string{"cmds", "ls"},
		Short:   "List all commands (LLM-friendly)",
		RunE: func(cmd *cobra.Command, args []string) error {
			all := getAllCommands()

			if group != "" {
				for _, g := range all {
					if g.Name == group {
						return output.Print(g.Commands)
					}
				}
				return output.PrintError("not_found", "group not found", nil)
			}

			return output.Print(all)
		},
	}

	cmd.Flags().StringVarP(&group, "group", "g", "", "Filter by group: social, comms, dev, productivity, news, knowledge, utility, system, security, marketing")

	return cmd
}

func getAllCommands() []Group {
	return []Group{
		{
			Name: "social",
			Commands: []Cmd{
				{Command: "alpha social twitter post", Desc: "Post a tweet (free tier: 1,500/mo)", Args: "[message]", Flags: "--reply-to"},
				{Command: "alpha social twitter delete", Desc: "Delete a tweet", Args: "[tweet-id]"},
				{Command: "alpha social twitter me", Desc: "Get your account info"},
				{Command: "alpha social reddit feed", Desc: "Get your home feed", Flags: "-l limit, -s sort"},
				{Command: "alpha social reddit subreddit", Desc: "Get subreddit posts", Args: "[name]", Flags: "-l limit, -s sort, -t time"},
				{Command: "alpha social reddit search", Desc: "Search Reddit", Args: "[query]", Flags: "-l limit, -r subreddit, -s sort"},
				{Command: "alpha social reddit user", Desc: "Get user info and posts", Args: "[username]", Flags: "-l limit"},
				{Command: "alpha social reddit comments", Desc: "Get post comments", Args: "[post-id]", Flags: "-l limit, -s sort"},
				{Command: "alpha social mastodon timeline", Desc: "Get timeline", Flags: "-l limit, -t type"},
				{Command: "alpha social mastodon post", Desc: "Post a toot", Args: "[content]", Flags: "-V visibility"},
				{Command: "alpha social mastodon search", Desc: "Search Mastodon", Args: "[query]", Flags: "-l limit, -t type"},
				{Command: "alpha social youtube search", Desc: "Search videos", Args: "[query]", Flags: "-l limit, -s order, -a after, -d duration"},
				{Command: "alpha social youtube video", Desc: "Get video details", Args: "[id]"},
				{Command: "alpha social youtube channel", Desc: "Get channel info", Args: "[id-or-handle]"},
				{Command: "alpha social youtube videos", Desc: "List channel videos", Args: "[channel-id]", Flags: "-l limit"},
				{Command: "alpha social youtube comments", Desc: "Get video comments", Args: "[video-id]", Flags: "-l limit, -s order"},
				{Command: "alpha social youtube trending", Desc: "Get trending videos", Flags: "-l limit, -r region, -c category"},
				{Command: "alpha social spotify search", Desc: "Search tracks/artists/albums", Args: "[query]", Flags: "-t type, -l limit"},
				{Command: "alpha social spotify track", Desc: "Get track details", Args: "[id]"},
				{Command: "alpha social spotify artist", Desc: "Get artist details", Args: "[id]"},
				{Command: "alpha social spotify album", Desc: "Get album details", Args: "[id]"},
			},
		},
		{
			Name: "comms",
			Commands: []Cmd{
				{Command: "alpha comms email list", Desc: "List emails", Flags: "-l limit, -m mailbox"},
				{Command: "alpha comms email read", Desc: "Read an email", Args: "[uid]", Flags: "-m mailbox"},
				{Command: "alpha comms email send", Desc: "Send an email", Args: "[body]", Flags: "--to, --subject, --cc"},
				{Command: "alpha comms email reply", Desc: "Reply to an email", Args: "[uid] [body]", Flags: "-m mailbox, -a reply-all"},
				{Command: "alpha comms email search", Desc: "Search emails", Args: "[query]", Flags: "-l limit, -m mailbox"},
				{Command: "alpha comms email mailboxes", Desc: "List mailboxes/folders"},
				{Command: "alpha comms slack channels", Desc: "List Slack channels"},
				{Command: "alpha comms slack messages", Desc: "Get channel messages", Args: "[channel]", Flags: "-l limit"},
				{Command: "alpha comms slack send", Desc: "Send Slack message", Args: "[message]", Flags: "-c channel"},
				{Command: "alpha comms discord guilds", Desc: "List Discord servers"},
				{Command: "alpha comms discord channels", Desc: "List guild channels", Args: "[guild-id]"},
				{Command: "alpha comms discord messages", Desc: "Get channel messages", Args: "[channel-id]", Flags: "-l limit"},
				{Command: "alpha comms discord send", Desc: "Send Discord message", Args: "[message]", Flags: "-c channel"},
				{Command: "alpha comms telegram chats", Desc: "List Telegram chats"},
				{Command: "alpha comms telegram messages", Desc: "Get chat messages", Args: "[chat-id]", Flags: "-l limit"},
				{Command: "alpha comms telegram send", Desc: "Send Telegram message", Args: "[message]", Flags: "-c chat"},
			},
		},
		{
			Name: "dev",
			Commands: []Cmd{
				{Command: "alpha dev github repos", Desc: "List repositories", Flags: "-l limit, -s sort, -u user"},
				{Command: "alpha dev github repo", Desc: "Get repo details", Args: "[owner/name]"},
				{Command: "alpha dev github issues", Desc: "List issues", Flags: "-r repo, -s state, -l limit, --labels"},
				{Command: "alpha dev github issue", Desc: "Get issue details", Args: "[owner/repo] [number]"},
				{Command: "alpha dev github prs", Desc: "List pull requests", Flags: "-r repo, -s state, -l limit"},
				{Command: "alpha dev github pr", Desc: "Get PR details", Args: "[owner/repo] [number]"},
				{Command: "alpha dev github notifications", Desc: "List notifications", Flags: "-l limit, -a all"},
				{Command: "alpha dev github search", Desc: "Search GitHub", Args: "[query]", Flags: "-t type, -l limit"},
				{Command: "alpha dev gitlab projects", Desc: "List projects", Flags: "-l limit"},
				{Command: "alpha dev gitlab issues", Desc: "List issues", Flags: "-p project, -s state, -l limit"},
				{Command: "alpha dev gitlab mrs", Desc: "List merge requests", Flags: "-p project, -s state, -l limit"},
				{Command: "alpha dev linear issues", Desc: "List Linear issues", Flags: "-t team, -s status, -l limit"},
				{Command: "alpha dev linear teams", Desc: "List Linear teams"},
				{Command: "alpha dev linear create", Desc: "Create Linear issue", Args: "[description]", Flags: "-t team, --title"},
				{Command: "alpha dev npm search", Desc: "Search npm packages", Args: "[query]", Flags: "-l limit"},
				{Command: "alpha dev npm info", Desc: "Get package info", Args: "[package]"},
				{Command: "alpha dev npm versions", Desc: "List package versions", Args: "[package]", Flags: "-l limit"},
				{Command: "alpha dev npm deps", Desc: "List dependencies", Args: "[package]", Flags: "-d dev"},
				{Command: "alpha dev pypi search", Desc: "Search PyPI packages", Args: "[query]"},
				{Command: "alpha dev pypi info", Desc: "Get package info", Args: "[package]"},
				{Command: "alpha dev pypi versions", Desc: "List package versions", Args: "[package]", Flags: "-l limit"},
				{Command: "alpha dev pypi deps", Desc: "List dependencies", Args: "[package]"},
				{Command: "alpha dev gist list", Desc: "List your gists", Flags: "-l limit"},
				{Command: "alpha dev gist get", Desc: "Get gist details", Args: "[id]"},
				{Command: "alpha dev gist create", Desc: "Create a gist", Args: "[content]", Flags: "-d desc, -f filename, --public"},
				{Command: "alpha dev sentry projects", Desc: "List Sentry projects", Flags: "-l limit, -o org"},
				{Command: "alpha dev sentry issues", Desc: "List project issues", Args: "[project-slug]", Flags: "-l limit, -o org, -q query"},
				{Command: "alpha dev sentry issue", Desc: "Get issue details", Args: "[issue-id]"},
				{Command: "alpha dev sentry events", Desc: "List events for issue", Args: "[issue-id]", Flags: "-l limit"},
				{Command: "alpha dev redis get", Desc: "Get key value", Args: "[key]"},
				{Command: "alpha dev redis set", Desc: "Set key value", Args: "[key] [value]", Flags: "--ttl"},
				{Command: "alpha dev redis del", Desc: "Delete keys", Args: "[key...]"},
				{Command: "alpha dev redis keys", Desc: "List keys matching pattern", Args: "[pattern]", Flags: "-l limit"},
				{Command: "alpha dev redis info", Desc: "Get Redis server info"},
				{Command: "alpha dev prometheus query", Desc: "Instant PromQL query", Args: "[promql]"},
				{Command: "alpha dev prometheus range", Desc: "Range PromQL query", Args: "[promql]", Flags: "-s start, -e end, --step"},
				{Command: "alpha dev prometheus alerts", Desc: "List alerts", Flags: "--state"},
				{Command: "alpha dev prometheus targets", Desc: "List scrape targets", Flags: "--state"},
				{Command: "alpha dev kube pods", Desc: "List pods", Flags: "-n namespace, -a all"},
				{Command: "alpha dev kube logs", Desc: "Get pod logs", Args: "[pod]", Flags: "-n namespace, -t tail, -c container"},
				{Command: "alpha dev kube deployments", Desc: "List deployments", Flags: "-n namespace, -a all"},
				{Command: "alpha dev kube services", Desc: "List services", Flags: "-n namespace, -a all"},
				{Command: "alpha dev kube describe", Desc: "Describe a resource", Args: "[resource] [name]", Flags: "-n namespace"},
				{Command: "alpha dev db query", Desc: "Execute SQL query (read-only)", Args: "[db-path] [sql]", Flags: "-l limit"},
				{Command: "alpha dev db schema", Desc: "Show database schema", Args: "[db-path]"},
				{Command: "alpha dev db tables", Desc: "List tables", Args: "[db-path]"},
				{Command: "alpha dev s3 buckets", Desc: "List S3 buckets"},
				{Command: "alpha dev s3 ls", Desc: "List objects in bucket", Args: "[s3-path]", Flags: "-r recursive"},
				{Command: "alpha dev s3 get", Desc: "Download object", Args: "[s3-path] [local-path]"},
				{Command: "alpha dev s3 put", Desc: "Upload object", Args: "[local-path] [s3-path]"},
				{Command: "alpha dev s3 presign", Desc: "Generate presigned URL", Args: "[s3-path]", Flags: "--expires"},
			},
		},
		{
			Name: "productivity",
			Commands: []Cmd{
				{Command: "alpha productivity notion search", Desc: "Search Notion", Args: "[query]", Flags: "-l limit"},
				{Command: "alpha productivity notion page", Desc: "Get page content", Args: "[page-id]"},
				{Command: "alpha productivity notion database", Desc: "Query database", Args: "[database-id]", Flags: "-l limit"},
				{Command: "alpha productivity todoist tasks", Desc: "List tasks", Flags: "-p project, -f filter"},
				{Command: "alpha productivity todoist projects", Desc: "List projects"},
				{Command: "alpha productivity todoist add", Desc: "Add a task", Args: "[content]", Flags: "-p project, -d due, --priority"},
				{Command: "alpha productivity todoist complete", Desc: "Complete a task", Args: "[task-id]"},
				{Command: "alpha productivity logseq graphs", Desc: "List configured graphs"},
				{Command: "alpha productivity logseq pages", Desc: "List pages in graph", Flags: "-g graph, -l limit"},
				{Command: "alpha productivity logseq read", Desc: "Read page content", Args: "[page]", Flags: "-g graph"},
				{Command: "alpha productivity logseq write", Desc: "Create/update page", Args: "[page] [content]", Flags: "-g graph, -a append"},
				{Command: "alpha productivity logseq search", Desc: "Search pages by content", Args: "[query]", Flags: "-g graph, -l limit, -c case-sensitive"},
				{Command: "alpha productivity logseq journal", Desc: "Get/create journal entry", Flags: "-g graph, -d date, -c content"},
				{Command: "alpha productivity logseq recent", Desc: "List recently modified pages", Flags: "-g graph, -l limit, -d days"},
			},
		},
		{
			Name: "news",
			Commands: []Cmd{
				{Command: "alpha news hn top", Desc: "HN top stories", Flags: "-l limit"},
				{Command: "alpha news hn new", Desc: "HN new stories", Flags: "-l limit"},
				{Command: "alpha news hn best", Desc: "HN best stories", Flags: "-l limit"},
				{Command: "alpha news hn ask", Desc: "Ask HN stories", Flags: "-l limit"},
				{Command: "alpha news hn show", Desc: "Show HN stories", Flags: "-l limit"},
				{Command: "alpha news hn item", Desc: "Get item with comments", Args: "[id]", Flags: "-c comments"},
				{Command: "alpha news feeds fetch", Desc: "Fetch RSS/Atom feed", Args: "[url]", Flags: "-l limit, -s summary-len"},
				{Command: "alpha news feeds list", Desc: "List saved feeds"},
				{Command: "alpha news feeds add", Desc: "Save a feed", Args: "[url]", Flags: "-n name"},
				{Command: "alpha news feeds read", Desc: "Fetch saved feed by name", Args: "[name]", Flags: "-l limit, -s summary-len"},
				{Command: "alpha news feeds remove", Desc: "Remove saved feed", Args: "[name-or-url]"},
				{Command: "alpha news newsapi headlines", Desc: "Get top headlines", Flags: "--country, --category, -l limit"},
				{Command: "alpha news newsapi search", Desc: "Search news", Args: "[query]", Flags: "--sort, -l limit"},
				{Command: "alpha news newsapi sources", Desc: "List news sources", Flags: "--category, --country"},
			},
		},
		{
			Name: "knowledge",
			Commands: []Cmd{
				{Command: "alpha knowledge wiki search", Desc: "Search Wikipedia", Args: "[query]", Flags: "-l limit"},
				{Command: "alpha knowledge wiki summary", Desc: "Get article summary", Args: "[title]", Flags: "-s sentences"},
				{Command: "alpha knowledge wiki article", Desc: "Get full article", Args: "[title]", Flags: "-c chars"},
				{Command: "alpha knowledge so search", Desc: "Search StackOverflow", Args: "[query]", Flags: "-l limit, -t tagged, -s site"},
				{Command: "alpha knowledge so question", Desc: "Get question details", Args: "[id]", Flags: "-s site"},
				{Command: "alpha knowledge so answers", Desc: "Get answers", Args: "[question-id]", Flags: "-l limit, -s site"},
				{Command: "alpha knowledge dict define", Desc: "Get word definition", Args: "[word]", Flags: "-l limit"},
				{Command: "alpha knowledge dict synonyms", Desc: "Get synonyms", Args: "[word]"},
				{Command: "alpha knowledge dict antonyms", Desc: "Get antonyms", Args: "[word]"},
			},
		},
		{
			Name: "utility",
			Commands: []Cmd{
				{Command: "alpha utility weather now", Desc: "Current weather", Args: "[location]"},
				{Command: "alpha utility weather forecast", Desc: "Weather forecast", Args: "[location]", Flags: "-d days"},
				{Command: "alpha utility crypto price", Desc: "Get crypto prices", Args: "[coins...]"},
				{Command: "alpha utility crypto info", Desc: "Get coin details", Args: "[coin]"},
				{Command: "alpha utility crypto top", Desc: "Top coins by market cap", Flags: "-l limit"},
				{Command: "alpha utility crypto trending", Desc: "Trending coins"},
				{Command: "alpha utility crypto search", Desc: "Search for coins", Args: "[query]", Flags: "-l limit"},
				{Command: "alpha utility ip me", Desc: "Get your public IP and location"},
				{Command: "alpha utility ip lookup", Desc: "Lookup IP geolocation", Args: "[ip]"},
				{Command: "alpha utility geocode forward", Desc: "Address to coordinates", Args: "[address]"},
				{Command: "alpha utility geocode reverse", Desc: "Coordinates to address", Args: "[lat] [lon]"},
				{Command: "alpha utility timezone get", Desc: "Get time in timezone", Args: "[timezone]"},
				{Command: "alpha utility timezone ip", Desc: "Get timezone by IP", Args: "[ip]"},
				{Command: "alpha utility timezone list", Desc: "List all timezones"},
				{Command: "alpha utility paste create", Desc: "Create a paste", Args: "[content]", Flags: "-e expiry, -t title"},
				{Command: "alpha utility paste get", Desc: "Fetch a paste", Args: "[url]"},
				{Command: "alpha utility netdiag headers", Desc: "Get HTTP response headers", Args: "[url]"},
				{Command: "alpha utility netdiag ports", Desc: "Scan common ports", Args: "[host]", Flags: "-t timeout"},
				{Command: "alpha utility netdiag ping", Desc: "DNS resolve and timing", Args: "[host]"},
				{Command: "alpha utility domain whois", Desc: "WHOIS lookup for domain", Args: "[domain]"},
				{Command: "alpha utility domain dns", Desc: "DNS records for domain", Args: "[domain]", Flags: "-t type"},
				{Command: "alpha utility currency convert", Desc: "Convert between currencies", Args: "[amount] [from] [to]"},
				{Command: "alpha utility currency rate", Desc: "Get exchange rate between two currencies", Args: "[from] [to]"},
				{Command: "alpha utility wayback check", Desc: "Check if URL has archived snapshots", Args: "[url]"},
				{Command: "alpha utility holidays list", Desc: "List holidays for a country", Args: "[country-code] [year]"},
				{Command: "alpha utility translate text", Desc: "Translate text", Args: "[text]", Flags: "-f from, -t to"},
				{Command: "alpha utility stocks quote", Desc: "Get stock quote", Args: "[symbol]"},
				{Command: "alpha utility stocks search", Desc: "Search stocks", Args: "[query]"},
				{Command: "alpha utility urlshort shorten", Desc: "Shorten a URL", Args: "[url]"},
				{Command: "alpha utility speedtest run", Desc: "Full speed test (download, upload, latency)", Flags: "--size"},
				{Command: "alpha utility speedtest download", Desc: "Download speed test only", Flags: "--size"},
				{Command: "alpha utility speedtest upload", Desc: "Upload speed test only", Flags: "--size"},
				{Command: "alpha utility speedtest latency", Desc: "Latency test only"},
				{Command: "alpha utility dnsbench run", Desc: "Benchmark all DNS resolvers"},
				{Command: "alpha utility dnsbench test", Desc: "Test a specific DNS resolver", Args: "[resolver-ip]"},
				{Command: "alpha utility traceroute run", Desc: "Trace network path to host", Args: "[host]", Flags: "--max-hops"},
				{Command: "alpha utility wifi scan", Desc: "Scan nearby WiFi networks"},
				{Command: "alpha utility wifi current", Desc: "Show current WiFi connection"},
			},
		},
		{
			Name: "security",
			Commands: []Cmd{
				{Command: "alpha security virustotal url", Desc: "Scan a URL", Args: "[url]"},
				{Command: "alpha security virustotal domain", Desc: "Get domain report", Args: "[domain]"},
				{Command: "alpha security virustotal ip", Desc: "Get IP report", Args: "[ip]"},
				{Command: "alpha security virustotal hash", Desc: "Get file hash report", Args: "[hash]"},
				{Command: "alpha security shodan lookup", Desc: "Lookup IP (ports, vulns)", Args: "[ip]"},
				{Command: "alpha security crtsh lookup", Desc: "Certificate transparency lookup", Args: "[domain]", Flags: "-l limit, --expired"},
				{Command: "alpha security hibp password", Desc: "Check password in breaches", Args: "[password]"},
				{Command: "alpha security hibp breaches", Desc: "List public data breaches", Flags: "-l limit"},
			},
		},
		{
			Name: "marketing",
			Commands: []Cmd{
				{Command: "alpha marketing facebook-ads account", Desc: "Get ad account details"},
				{Command: "alpha marketing facebook-ads campaigns", Desc: "List campaigns", Flags: "-s status, -l limit"},
				{Command: "alpha marketing facebook-ads campaign-create", Desc: "Create a campaign", Flags: "--name, --objective, --status, --daily-budget, --special-categories"},
				{Command: "alpha marketing facebook-ads campaign-update", Desc: "Update a campaign", Args: "[campaign-id]", Flags: "--name, --status, --daily-budget"},
				{Command: "alpha marketing facebook-ads adsets", Desc: "List ad sets", Flags: "-c campaign-id, -s status, -l limit"},
				{Command: "alpha marketing facebook-ads adset-create", Desc: "Create an ad set", Flags: "--name, --campaign-id, --billing-event, --optimization-goal"},
				{Command: "alpha marketing facebook-ads ads", Desc: "List ads", Flags: "--adset-id, -s status, -l limit"},
				{Command: "alpha marketing facebook-ads insights", Desc: "Get performance insights", Flags: "--object-id, --date-start, --date-stop, --level, -f fields"},
				{Command: "alpha marketing amazon-sp orders", Desc: "List orders", Flags: "-s status, --after, --before, --marketplace, -l limit"},
				{Command: "alpha marketing amazon-sp order", Desc: "Get order details", Args: "[order-id]"},
				{Command: "alpha marketing amazon-sp order-items", Desc: "Get items for an order", Args: "[order-id]"},
				{Command: "alpha marketing amazon-sp inventory", Desc: "List FBA inventory", Flags: "--sku, --marketplace, -l limit"},
				{Command: "alpha marketing amazon-sp report-create", Desc: "Create a report request", Flags: "--type, --start, --end, --marketplace"},
				{Command: "alpha marketing amazon-sp report-status", Desc: "Get report status", Args: "[report-id]"},
				{Command: "alpha marketing shopify shop", Desc: "Get store info"},
				{Command: "alpha marketing shopify orders", Desc: "List orders", Flags: "-l limit, --status, --since, --financial"},
				{Command: "alpha marketing shopify order", Desc: "Get order details", Args: "[id]"},
				{Command: "alpha marketing shopify products", Desc: "List products", Flags: "-l limit, --status, --vendor, --collection"},
				{Command: "alpha marketing shopify product", Desc: "Get product details", Args: "[id]"},
				{Command: "alpha marketing shopify customers", Desc: "List customers", Flags: "-l limit"},
				{Command: "alpha marketing shopify customer-search", Desc: "Search customers", Args: "[query]", Flags: "-l limit"},
				{Command: "alpha marketing shopify inventory", Desc: "List inventory levels", Flags: "--location, -l limit"},
				{Command: "alpha marketing shopify inventory-set", Desc: "Set inventory level", Flags: "--item, --location, --available"},
			},
		},
		{
			Name: "setup",
			Commands: []Cmd{
				{Command: "alpha setup list", Desc: "List services needing setup", Flags: "-a all"},
				{Command: "alpha setup show", Desc: "Show setup instructions", Args: "[service]"},
				{Command: "alpha setup set", Desc: "Set credential for service", Args: "[service] [key] [value]"},
			},
		},
		{
			Name: "config",
			Commands: []Cmd{
				{Command: "alpha config path", Desc: "Show config file path"},
				{Command: "alpha config list", Desc: "List all config (redacted)"},
				{Command: "alpha config set", Desc: "Set a config value", Args: "[key] [value]"},
				{Command: "alpha config get", Desc: "Get a config value", Args: "[key]"},
			},
		},
		{
			Name: "system",
			Commands: []Cmd{
				{Command: "alpha system clipboard get", Desc: "Get clipboard content", Flags: "-m max-length, -r raw"},
				{Command: "alpha system clipboard set", Desc: "Set clipboard content", Args: "[text]"},
				{Command: "alpha system clipboard clear", Desc: "Clear clipboard"},
				{Command: "alpha system clipboard copy", Desc: "Copy file to clipboard", Args: "[file]"},
				{Command: "alpha system clipboard history", Desc: "Clipboard history info"},
				{Command: "alpha system notes list", Desc: "List all notes", Flags: "-f folder, -l limit"},
				{Command: "alpha system notes folders", Desc: "List all folders"},
				{Command: "alpha system notes read", Desc: "Read a note by name", Args: "[name]", Flags: "-f folder"},
				{Command: "alpha system notes create", Desc: "Create a new note", Args: "[name] [body]", Flags: "-f folder"},
				{Command: "alpha system notes search", Desc: "Search notes", Args: "[query]", Flags: "-l limit"},
				{Command: "alpha system notes append", Desc: "Append text to note", Args: "[name] [text]", Flags: "-f folder"},
				{Command: "alpha system mail accounts", Desc: "List mail accounts"},
				{Command: "alpha system mail mailboxes", Desc: "List mailboxes/folders", Flags: "-a account"},
				{Command: "alpha system mail list", Desc: "List recent messages", Flags: "-m mailbox, -a account, -l limit, -u unread"},
				{Command: "alpha system mail read", Desc: "Read a message by ID", Args: "[id]"},
				{Command: "alpha system mail search", Desc: "Search messages", Args: "[query]", Flags: "-l limit, -m mailbox, -a account"},
				{Command: "alpha system mail send", Desc: "Send an email", Flags: "--to, --subject, --body, --cc, --bcc, -a account"},
				{Command: "alpha system mail unread", Desc: "List unread messages", Flags: "-l limit, -a account"},
				{Command: "alpha system mail count", Desc: "Get unread message count", Flags: "-a account"},
				{Command: "alpha system safari tabs", Desc: "List all open tabs", Flags: "-w window"},
				{Command: "alpha system safari url", Desc: "Get URL of current tab", Flags: "-w window, -t tab"},
				{Command: "alpha system safari title", Desc: "Get title of current tab"},
				{Command: "alpha system safari open", Desc: "Open URL in new tab", Args: "[url]", Flags: "-n new-window"},
				{Command: "alpha system safari close", Desc: "Close current tab", Flags: "-w window, -t tab, --window-close"},
				{Command: "alpha system safari bookmarks", Desc: "List bookmarks", Flags: "-f folder, -l limit"},
				{Command: "alpha system safari reading-list", Desc: "List Reading List items", Flags: "-l limit"},
				{Command: "alpha system safari add-reading", Desc: "Add URL to Reading List", Args: "[url]", Flags: "-t title"},
				{Command: "alpha system safari history", Desc: "Get recent history", Flags: "-l limit, -d days, -s search"},
				{Command: "alpha system apple-calendar events", Desc: "List upcoming events", Flags: "-d days, -c calendar"},
				{Command: "alpha system apple-calendar today", Desc: "List today's events"},
				{Command: "alpha system apple-calendar create", Desc: "Create event", Flags: "--title, --start, --end, --desc"},
				{Command: "alpha system contacts list", Desc: "List contacts", Flags: "-l limit, -g group"},
				{Command: "alpha system contacts search", Desc: "Search contacts", Args: "[query]", Flags: "-l limit"},
				{Command: "alpha system finder info", Desc: "File/folder metadata", Args: "[path]"},
				{Command: "alpha system finder search", Desc: "Search with Spotlight", Args: "[query]", Flags: "-l limit, -d dir"},
				{Command: "alpha system finder recent", Desc: "Recently modified files", Flags: "-l limit, -d dir"},
				{Command: "alpha system finder tags", Desc: "Get/set Finder tags", Args: "[path]"},
				{Command: "alpha system imessage chats", Desc: "List recent conversations", Flags: "-l limit"},
				{Command: "alpha system imessage read", Desc: "Read messages from chat", Args: "[chat-id]", Flags: "-l limit"},
				{Command: "alpha system imessage search", Desc: "Search messages", Args: "[query]", Flags: "-l limit"},
				{Command: "alpha system reminders list", Desc: "List reminders in a list (or all)", Args: "[name]", Flags: "--completed"},
				{Command: "alpha system reminders lists", Desc: "List reminder lists"},
				{Command: "alpha system reminders create", Desc: "Create a reminder", Args: "[title]", Flags: "--list, --due, --notes"},
				{Command: "alpha system sysinfo overview", Desc: "System overview (CPU, memory, disk, uptime)"},
				{Command: "alpha system sysinfo cpu", Desc: "CPU info and load averages"},
				{Command: "alpha system sysinfo memory", Desc: "Memory usage details"},
				{Command: "alpha system sysinfo disk", Desc: "Disk usage per mount"},
				{Command: "alpha system sysinfo processes", Desc: "Top processes by CPU/memory", Flags: "-l limit"},
				{Command: "alpha system sysinfo network", Desc: "Network interfaces and IPs"},
				{Command: "alpha system battery status", Desc: "Current battery charge and state"},
				{Command: "alpha system battery health", Desc: "Battery health, cycle count, capacity"},
				{Command: "alpha system diskhealth status", Desc: "Disk health and S.M.A.R.T. overview"},
				{Command: "alpha system diskhealth info", Desc: "Detailed disk info", Args: "[disk]"},
				{Command: "alpha system cleanup scan", Desc: "Full scan of cleanable directories (read-only)"},
				{Command: "alpha system cleanup caches", Desc: "Scan cache directories only"},
				{Command: "alpha system cleanup logs", Desc: "Scan log directories only"},
				{Command: "alpha system cleanup temp", Desc: "Scan temp directories only"},
			},
		},
	}
}
