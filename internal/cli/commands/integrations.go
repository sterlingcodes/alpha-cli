package commands

import (
	"github.com/spf13/cobra"

	"github.com/sterlingcodes/alpha-cli/internal/common/config"
	"github.com/sterlingcodes/alpha-cli/pkg/output"
)

const (
	statusNoAuth = "no_auth"
	statusReady  = "ready"
)

// Integration describes an available integration
type Integration struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Group       string   `json:"group"`
	Description string   `json:"desc"`
	AuthNeeded  bool     `json:"auth_needed"`
	Status      string   `json:"status"` // "ready", "needs_setup", "no_auth"
	Commands    []string `json:"commands"`
	SetupCmd    string   `json:"setup_cmd,omitempty"`
}

var allIntegrations = []Integration{
	// News - No Auth
	{
		ID:          "hackernews",
		Name:        "Hacker News",
		Group:       "news",
		Description: "Tech news, stories, and discussions from Hacker News",
		AuthNeeded:  false,
		Commands:    []string{"alpha news hn top", "alpha news hn new", "alpha news hn best", "alpha news hn ask", "alpha news hn show", "alpha news hn item [id]"},
	},
	{
		ID:          "rss",
		Name:        "RSS/Atom Feeds",
		Group:       "news",
		Description: "Fetch and manage RSS/Atom feeds from any source",
		AuthNeeded:  false,
		Commands:    []string{"alpha news feeds fetch [url]", "alpha news feeds list", "alpha news feeds add [url]", "alpha news feeds read [name]", "alpha news feeds remove [name]"},
	},
	{
		ID:          "newsapi",
		Name:        "NewsAPI",
		Group:       "news",
		Description: "Search news articles and get headlines from 80,000+ sources",
		AuthNeeded:  true,
		Commands:    []string{"alpha news newsapi headlines", "alpha news newsapi search [query]", "alpha news newsapi sources"},
		SetupCmd:    "alpha setup show newsapi",
	},

	// Knowledge - No Auth
	{
		ID:          "wikipedia",
		Name:        "Wikipedia",
		Group:       "knowledge",
		Description: "Search and read Wikipedia articles",
		AuthNeeded:  false,
		Commands:    []string{"alpha knowledge wiki search [query]", "alpha knowledge wiki summary [title]", "alpha knowledge wiki article [title]"},
	},
	{
		ID:          "stackexchange",
		Name:        "StackOverflow",
		Group:       "knowledge",
		Description: "Search programming Q&A from StackOverflow and StackExchange sites",
		AuthNeeded:  false,
		Commands:    []string{"alpha knowledge so search [query]", "alpha knowledge so question [id]", "alpha knowledge so answers [id]"},
	},
	{
		ID:          "dictionary",
		Name:        "Dictionary",
		Group:       "knowledge",
		Description: "Word definitions, synonyms, antonyms, and pronunciations",
		AuthNeeded:  false,
		Commands:    []string{"alpha knowledge dict define [word]", "alpha knowledge dict synonyms [word]", "alpha knowledge dict antonyms [word]"},
	},

	// Utility - No Auth
	{
		ID:          "weather",
		Name:        "Weather",
		Group:       "utility",
		Description: "Current weather and forecasts for any location",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility weather now [location]", "alpha utility weather forecast [location]"},
	},
	{
		ID:          "crypto",
		Name:        "CoinGecko",
		Group:       "utility",
		Description: "Cryptocurrency prices, market data, and trending coins",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility crypto price [coins...]", "alpha utility crypto info [coin]", "alpha utility crypto top", "alpha utility crypto trending", "alpha utility crypto search [query]"},
	},
	{
		ID:          "ipinfo",
		Name:        "IP Geolocation",
		Group:       "utility",
		Description: "IP address lookup with geolocation data",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility ip me", "alpha utility ip lookup [ip]"},
	},
	{
		ID:          "domain",
		Name:        "DNS/WHOIS/SSL",
		Group:       "utility",
		Description: "DNS lookups, WHOIS domain info, and SSL certificate inspection",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility domain dns [domain]", "alpha utility domain whois [domain]", "alpha utility domain ssl [domain]"},
	},
	{
		ID:          "currency",
		Name:        "Currency Exchange",
		Group:       "utility",
		Description: "Real-time currency exchange rates and conversion",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility currency rate [from] [to]", "alpha utility currency convert [amount] [from] [to]", "alpha utility currency list"},
	},
	{
		ID:          "wayback",
		Name:        "Wayback Machine",
		Group:       "utility",
		Description: "Check archived versions of websites via Internet Archive",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility wayback check [url]", "alpha utility wayback latest [url]", "alpha utility wayback snapshots [url]"},
	},
	{
		ID:          "holidays",
		Name:        "Public Holidays",
		Group:       "utility",
		Description: "Public holidays by country and year",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility holidays list [country] [year]", "alpha utility holidays next [country]", "alpha utility holidays countries"},
	},
	{
		ID:          "translate",
		Name:        "Translation",
		Group:       "utility",
		Description: "Translate text between languages",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility translate text [text] --from [lang] --to [lang]", "alpha utility translate languages"},
	},
	{
		ID:          "urlshort",
		Name:        "URL Shortener",
		Group:       "utility",
		Description: "Shorten and expand URLs",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility url shorten [url]", "alpha utility url expand [short-url]"},
	},
	// Utility - Auth Required
	{
		ID:          "stocks",
		Name:        "Stock Market",
		Group:       "utility",
		Description: "Stock quotes, search, and company info via Alpha Vantage",
		AuthNeeded:  true,
		Commands:    []string{"alpha utility stocks quote [symbol]", "alpha utility stocks search [query]", "alpha utility stocks info [symbol]"},
		SetupCmd:    "alpha setup show alphavantage",
	},

	// Dev - No Auth
	{
		ID:          "npm",
		Name:        "npm Registry",
		Group:       "dev",
		Description: "Search npm packages, get info, versions, and dependencies",
		AuthNeeded:  false,
		Commands:    []string{"alpha dev npm search [query]", "alpha dev npm info [package]", "alpha dev npm versions [package]", "alpha dev npm deps [package]"},
	},
	{
		ID:          "pypi",
		Name:        "PyPI Registry",
		Group:       "dev",
		Description: "Search Python packages, get info, versions, and dependencies",
		AuthNeeded:  false,
		Commands:    []string{"alpha dev pypi search [query]", "alpha dev pypi info [package]", "alpha dev pypi versions [package]", "alpha dev pypi deps [package]"},
	},
	{
		ID:          "dockerhub",
		Name:        "Docker Hub",
		Group:       "dev",
		Description: "Search Docker images, get tags, and inspect manifests",
		AuthNeeded:  false,
		Commands:    []string{"alpha dev dockerhub search [query]", "alpha dev dockerhub image [name]", "alpha dev dockerhub tags [name]", "alpha dev dockerhub inspect [name:tag]"},
	},

	// Dev - Auth Required
	{
		ID:          "github",
		Name:        "GitHub",
		Group:       "dev",
		Description: "Repos, issues, PRs, notifications, and search on GitHub",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev github repos", "alpha dev github repo [owner/name]", "alpha dev github issues", "alpha dev github issue [repo] [num]", "alpha dev github prs -r [repo]", "alpha dev github pr [repo] [num]", "alpha dev github notifications", "alpha dev github search [query]"},
		SetupCmd:    "alpha setup show github",
	},
	{
		ID:          "gitlab",
		Name:        "GitLab",
		Group:       "dev",
		Description: "Projects, issues, and merge requests on GitLab",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev gitlab projects", "alpha dev gitlab issues", "alpha dev gitlab mrs"},
		SetupCmd:    "alpha setup show gitlab",
	},
	{
		ID:          "linear",
		Name:        "Linear",
		Group:       "dev",
		Description: "Issues and project management with Linear",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev linear issues", "alpha dev linear teams", "alpha dev linear create [desc]"},
		SetupCmd:    "alpha setup show linear",
	},
	{
		ID:          "jira",
		Name:        "Jira",
		Group:       "dev",
		Description: "Issues, projects, and sprint management with Atlassian Jira",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev jira issues", "alpha dev jira issue [key]", "alpha dev jira projects", "alpha dev jira create [summary]", "alpha dev jira transition [key] [status]"},
		SetupCmd:    "alpha setup show jira",
	},
	{
		ID:          "cloudflare",
		Name:        "Cloudflare",
		Group:       "dev",
		Description: "DNS, zones, cache purge, and analytics via Cloudflare",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev cloudflare zones", "alpha dev cloudflare zone [id]", "alpha dev cloudflare dns [zone-id]", "alpha dev cloudflare purge [zone-id]", "alpha dev cloudflare analytics [zone-id]"},
		SetupCmd:    "alpha setup show cloudflare",
	},
	{
		ID:          "vercel",
		Name:        "Vercel",
		Group:       "dev",
		Description: "Projects, deployments, domains, and environment variables on Vercel",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev vercel projects", "alpha dev vercel project [name]", "alpha dev vercel deployments [project]", "alpha dev vercel domains", "alpha dev vercel env [project]"},
		SetupCmd:    "alpha setup show vercel",
	},
	{
		ID:          "sentry",
		Name:        "Sentry",
		Group:       "dev",
		Description: "Error tracking: projects, issues, and events from Sentry",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev sentry projects", "alpha dev sentry issues [project-slug]", "alpha dev sentry issue [issue-id]", "alpha dev sentry events [issue-id]"},
		SetupCmd:    "alpha setup show sentry",
	},
	{
		ID:          "s3",
		Name:        "AWS S3",
		Group:       "dev",
		Description: "List buckets, browse objects, upload/download, and generate presigned URLs",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev s3 buckets", "alpha dev s3 ls [s3-path]", "alpha dev s3 get [s3-path] [local-path]", "alpha dev s3 put [local-path] [s3-path]", "alpha dev s3 presign [s3-path]"},
		SetupCmd:    "alpha setup show s3",
	},
	{
		ID:          "redis",
		Name:        "Redis",
		Group:       "dev",
		Description: "Get/set keys, list keys, and view server info on Redis",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev redis get [key]", "alpha dev redis set [key] [value]", "alpha dev redis del [key...]", "alpha dev redis keys [pattern]", "alpha dev redis info"},
		SetupCmd:    "alpha setup show redis",
	},
	{
		ID:          "prometheus",
		Name:        "Prometheus",
		Group:       "dev",
		Description: "PromQL queries, alerts, and scrape targets from Prometheus",
		AuthNeeded:  true,
		Commands:    []string{"alpha dev prometheus query [promql]", "alpha dev prometheus range [promql]", "alpha dev prometheus alerts", "alpha dev prometheus targets"},
		SetupCmd:    "alpha setup show prometheus",
	},
	// Dev - No Auth
	{
		ID:          "gist",
		Name:        "GitHub Gists",
		Group:       "dev",
		Description: "Create, list, and read GitHub Gists",
		AuthNeeded:  false,
		Commands:    []string{"alpha dev gist list", "alpha dev gist get [id]", "alpha dev gist create [content]"},
	},
	{
		ID:          "kubernetes",
		Name:        "Kubernetes",
		Group:       "dev",
		Description: "Pods, logs, deployments, services, and resource descriptions via kubectl",
		AuthNeeded:  false,
		Commands:    []string{"alpha dev kube pods", "alpha dev kube logs [pod]", "alpha dev kube deployments", "alpha dev kube services", "alpha dev kube describe [resource] [name]"},
	},
	{
		ID:          "database",
		Name:        "Database (SQLite)",
		Group:       "dev",
		Description: "Query, inspect schema, and list tables in SQLite databases",
		AuthNeeded:  false,
		Commands:    []string{"alpha dev db query [db-path] [sql]", "alpha dev db schema [db-path]", "alpha dev db tables [db-path]"},
	},

	// Social - Auth Required
	{
		ID:          "twitter",
		Name:        "Twitter/X",
		Group:       "social",
		Description: "Post tweets, delete tweets, get account info (free tier: 1,500 posts/month)",
		AuthNeeded:  true,
		Commands:    []string{"alpha social twitter post [msg]", "alpha social twitter delete [id]", "alpha social twitter me", "alpha social twitter --reply-to [id] [msg]"},
		SetupCmd:    "alpha setup show twitter",
	},
	{
		ID:          "reddit",
		Name:        "Reddit",
		Group:       "social",
		Description: "Browse feeds, subreddits, search, users, and comments (free tier: 100 req/min)",
		AuthNeeded:  true,
		Commands:    []string{"alpha social reddit feed", "alpha social reddit subreddit [name]", "alpha social reddit search [query]", "alpha social reddit user [name]", "alpha social reddit comments [post-id]"},
		SetupCmd:    "alpha setup show reddit",
	},
	{
		ID:          "mastodon",
		Name:        "Mastodon",
		Group:       "social",
		Description: "Fediverse: timelines, posting, and search",
		AuthNeeded:  true,
		Commands:    []string{"alpha social mastodon timeline", "alpha social mastodon post [content]", "alpha social mastodon search [query]"},
		SetupCmd:    "alpha setup show mastodon",
	},
	{
		ID:          "youtube",
		Name:        "YouTube",
		Group:       "social",
		Description: "Search videos, get channel info, video metrics, comments, and trending",
		AuthNeeded:  true,
		Commands:    []string{"alpha social youtube search [query]", "alpha social youtube video [id]", "alpha social youtube channel [id]", "alpha social youtube videos [channel-id]", "alpha social youtube comments [video-id]", "alpha social youtube trending"},
		SetupCmd:    "alpha setup show youtube",
	},
	{
		ID:          "spotify",
		Name:        "Spotify",
		Group:       "social",
		Description: "Search tracks, artists, and albums on Spotify",
		AuthNeeded:  true,
		Commands:    []string{"alpha social spotify search [query]", "alpha social spotify track [id]", "alpha social spotify artist [id]", "alpha social spotify album [id]"},
		SetupCmd:    "alpha setup show spotify",
	},

	// Communication - Auth Required
	{
		ID:          "email",
		Name:        "Email (IMAP/SMTP)",
		Group:       "comms",
		Description: "Read, search, send, and reply to emails via IMAP/SMTP (Gmail, Outlook, Yahoo, etc.)",
		AuthNeeded:  true,
		Commands:    []string{"alpha comms email list", "alpha comms email read [uid]", "alpha comms email send [body]", "alpha comms email reply [uid] [body]", "alpha comms email search [query]", "alpha comms email mailboxes"},
		SetupCmd:    "alpha setup show email",
	},
	{
		ID:          "slack",
		Name:        "Slack",
		Group:       "comms",
		Description: "Channels, messages, users, DMs, and search in Slack workspaces",
		AuthNeeded:  true,
		Commands:    []string{"alpha comms slack channels", "alpha comms slack messages [channel]", "alpha comms slack send [channel] [msg]", "alpha comms slack users", "alpha comms slack dm [user] [msg]", "alpha comms slack search [query]"},
		SetupCmd:    "alpha setup show slack",
	},
	{
		ID:          "discord",
		Name:        "Discord",
		Group:       "comms",
		Description: "Servers (guilds), channels, messages, and DMs in Discord",
		AuthNeeded:  true,
		Commands:    []string{"alpha comms discord guilds", "alpha comms discord channels [guild]", "alpha comms discord messages [channel]", "alpha comms discord send [channel] [msg]", "alpha comms discord dm [user] [msg]"},
		SetupCmd:    "alpha setup show discord",
	},
	{
		ID:          "telegram",
		Name:        "Telegram",
		Group:       "comms",
		Description: "Chats, messages, and forwarding via Telegram bot",
		AuthNeeded:  true,
		Commands:    []string{"alpha comms telegram me", "alpha comms telegram chats", "alpha comms telegram updates", "alpha comms telegram send [chat] [msg]", "alpha comms telegram forward [from] [to] [msg-id]"},
		SetupCmd:    "alpha setup show telegram",
	},
	{
		ID:          "twilio",
		Name:        "Twilio (SMS)",
		Group:       "comms",
		Description: "Send and manage SMS messages via Twilio",
		AuthNeeded:  true,
		Commands:    []string{"alpha comms twilio send [to] [msg]", "alpha comms twilio messages", "alpha comms twilio message [sid]", "alpha comms twilio account"},
		SetupCmd:    "alpha setup show twilio",
	},
	// Communication - No Auth
	{
		ID:          "notify",
		Name:        "Push Notifications",
		Group:       "comms",
		Description: "Send push notifications via ntfy.sh (no auth) or Pushover (auth)",
		AuthNeeded:  false,
		Commands:    []string{"alpha comms notify ntfy [topic] [msg]", "alpha comms notify pushover [msg]"},
	},
	{
		ID:          "webhook",
		Name:        "Webhooks",
		Group:       "comms",
		Description: "Send data to webhooks (generic, Slack, Discord)",
		AuthNeeded:  false,
		Commands:    []string{"alpha comms webhook send [url] [data]", "alpha comms webhook slack [url] [msg]", "alpha comms webhook discord [url] [msg]"},
	},

	// Productivity - Auth Required
	{
		ID:          "calendar",
		Name:        "Google Calendar",
		Group:       "productivity",
		Description: "View and create calendar events",
		AuthNeeded:  true,
		Commands:    []string{"alpha productivity calendar events", "alpha productivity calendar today", "alpha productivity calendar create"},
		SetupCmd:    "alpha setup show calendar",
	},
	{
		ID:          "notion",
		Name:        "Notion",
		Group:       "productivity",
		Description: "Search pages and query databases in Notion",
		AuthNeeded:  true,
		Commands:    []string{"alpha productivity notion search [query]", "alpha productivity notion page [id]", "alpha productivity notion database [id]"},
		SetupCmd:    "alpha setup show notion",
	},
	{
		ID:          "todoist",
		Name:        "Todoist",
		Group:       "productivity",
		Description: "Tasks and projects in Todoist",
		AuthNeeded:  true,
		Commands:    []string{"alpha productivity todoist tasks", "alpha productivity todoist projects", "alpha productivity todoist add [task]", "alpha productivity todoist complete [id]"},
		SetupCmd:    "alpha setup show todoist",
	},
	{
		ID:          "trello",
		Name:        "Trello",
		Group:       "productivity",
		Description: "Boards, lists, and cards in Trello",
		AuthNeeded:  true,
		Commands:    []string{"alpha productivity trello boards", "alpha productivity trello board [id]", "alpha productivity trello cards [board-id]", "alpha productivity trello card [id]", "alpha productivity trello create [name]"},
		SetupCmd:    "alpha setup show trello",
	},
	// Productivity - Local (Path Required)
	{
		ID:          "logseq",
		Name:        "Logseq",
		Group:       "productivity",
		Description: "Local Logseq graphs - read/write pages, search, journals",
		AuthNeeded:  true,
		Commands:    []string{"alpha productivity logseq graphs", "alpha productivity logseq pages", "alpha productivity logseq read [page]", "alpha productivity logseq write [page] [content]", "alpha productivity logseq search [query]", "alpha productivity logseq journal", "alpha productivity logseq recent"},
		SetupCmd:    "alpha setup show logseq",
	},

	// Productivity - Local (Path Required)
	{
		ID:          "obsidian",
		Name:        "Obsidian",
		Group:       "productivity",
		Description: "Local Obsidian vaults - read/write notes, search, daily notes",
		AuthNeeded:  true,
		Commands:    []string{"alpha productivity obsidian vaults", "alpha productivity obsidian notes", "alpha productivity obsidian read [note]", "alpha productivity obsidian write [note] [content]", "alpha productivity obsidian search [query]", "alpha productivity obsidian daily", "alpha productivity obsidian recent"},
		SetupCmd:    "alpha setup show obsidian",
	},

	// Marketing - Auth Required
	{
		ID:          "facebook-ads",
		Name:        "Facebook Ads (Meta)",
		Group:       "marketing",
		Description: "Manage Facebook/Meta ad campaigns, ad sets, ads, and view performance insights",
		AuthNeeded:  true,
		Commands:    []string{"alpha marketing facebook-ads account", "alpha marketing facebook-ads campaigns", "alpha marketing facebook-ads campaign-create", "alpha marketing facebook-ads adsets", "alpha marketing facebook-ads ads", "alpha marketing facebook-ads insights"},
		SetupCmd:    "alpha setup show facebook-ads",
	},
	{
		ID:          "amazon-sp",
		Name:        "Amazon Selling Partner",
		Group:       "marketing",
		Description: "Manage Amazon seller orders, inventory, and reports via SP-API",
		AuthNeeded:  true,
		Commands:    []string{"alpha marketing amazon-sp orders", "alpha marketing amazon-sp order [id]", "alpha marketing amazon-sp order-items [id]", "alpha marketing amazon-sp inventory", "alpha marketing amazon-sp report-create", "alpha marketing amazon-sp report-status [id]"},
		SetupCmd:    "alpha setup show amazon-sp",
	},
	{
		ID:          "shopify",
		Name:        "Shopify",
		Group:       "marketing",
		Description: "Manage Shopify store: orders, products, customers, and inventory",
		AuthNeeded:  true,
		Commands:    []string{"alpha marketing shopify shop", "alpha marketing shopify orders", "alpha marketing shopify order [id]", "alpha marketing shopify products", "alpha marketing shopify product [id]", "alpha marketing shopify customers", "alpha marketing shopify customer-search [query]", "alpha marketing shopify inventory", "alpha marketing shopify inventory-set"},
		SetupCmd:    "alpha setup show shopify",
	},

	// System - macOS Only (No Auth)
	{
		ID:          "reminders",
		Name:        "Apple Reminders",
		Group:       "system",
		Description: "Manage Apple Reminders via AppleScript (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system reminders lists", "alpha system reminders list [name]", "alpha system reminders add [title]", "alpha system reminders complete [id]", "alpha system reminders delete [id]", "alpha system reminders today", "alpha system reminders overdue"},
	},
	{
		ID:          "notes",
		Name:        "Apple Notes",
		Group:       "system",
		Description: "Read and manage Apple Notes via AppleScript (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system notes folders", "alpha system notes list", "alpha system notes read [name]", "alpha system notes search [query]", "alpha system notes create [name] [body]", "alpha system notes append [name] [text]"},
	},
	{
		ID:          "apple-calendar",
		Name:        "Apple Calendar",
		Group:       "system",
		Description: "Manage Apple Calendar events via AppleScript (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system apple-calendar calendars", "alpha system apple-calendar today", "alpha system apple-calendar events", "alpha system apple-calendar event [title]", "alpha system apple-calendar create [title]", "alpha system apple-calendar upcoming", "alpha system apple-calendar week"},
	},
	{
		ID:          "contacts",
		Name:        "Apple Contacts",
		Group:       "system",
		Description: "Search and manage Apple Contacts via AppleScript (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system contacts list", "alpha system contacts search [query]", "alpha system contacts get [name]", "alpha system contacts groups", "alpha system contacts group [name]", "alpha system contacts create [name]"},
	},
	{
		ID:          "finder",
		Name:        "Finder",
		Group:       "system",
		Description: "Finder operations, file info, tags, Spotlight search (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system finder open [path]", "alpha system finder reveal [path]", "alpha system finder info [path]", "alpha system finder list [path]", "alpha system finder tags [path]", "alpha system finder tag [path] [tag]", "alpha system finder trash [path]", "alpha system finder search [query]"},
	},
	{
		ID:          "safari",
		Name:        "Safari",
		Group:       "system",
		Description: "Safari tabs, bookmarks, reading list, history (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system safari tabs", "alpha system safari url", "alpha system safari open [url]", "alpha system safari bookmarks", "alpha system safari reading-list", "alpha system safari add-reading [url]", "alpha system safari history"},
	},
	{
		ID:          "clipboard",
		Name:        "Clipboard",
		Group:       "system",
		Description: "Get/set macOS clipboard content (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system clipboard get", "alpha system clipboard set [text]", "alpha system clipboard clear", "alpha system clipboard copy [file]"},
	},
	{
		ID:          "imessage",
		Name:        "iMessage",
		Group:       "system",
		Description: "Send and read iMessages via Messages.app (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system imessage send [recipient] [message]", "alpha system imessage chats", "alpha system imessage read [contact]", "alpha system imessage search [query]", "alpha system imessage unread"},
	},
	{
		ID:          "apple-mail",
		Name:        "Apple Mail",
		Group:       "system",
		Description: "Read and send emails via Apple Mail (macOS only)",
		AuthNeeded:  false,
		Commands:    []string{"alpha system mail accounts", "alpha system mail mailboxes", "alpha system mail list", "alpha system mail read [id]", "alpha system mail search [query]", "alpha system mail send", "alpha system mail unread", "alpha system mail count"},
	},

	// Security - Auth Required
	{
		ID:          "virustotal",
		Name:        "VirusTotal",
		Group:       "security",
		Description: "Scan URLs, domains, IPs, and file hashes for threats via VirusTotal",
		AuthNeeded:  true,
		Commands:    []string{"alpha security virustotal url [url]", "alpha security virustotal domain [domain]", "alpha security virustotal ip [ip]", "alpha security virustotal hash [hash]"},
		SetupCmd:    "alpha setup show virustotal",
	},
	// Security - No Auth
	{
		ID:          "shodan",
		Name:        "Shodan",
		Group:       "security",
		Description: "IP lookup for open ports and vulnerabilities via Shodan",
		AuthNeeded:  false,
		Commands:    []string{"alpha security shodan lookup [ip]"},
	},
	{
		ID:          "crtsh",
		Name:        "crt.sh",
		Group:       "security",
		Description: "Certificate transparency log lookups via crt.sh",
		AuthNeeded:  false,
		Commands:    []string{"alpha security crtsh lookup [domain]"},
	},
	{
		ID:          "hibp",
		Name:        "Have I Been Pwned",
		Group:       "security",
		Description: "Check passwords against breaches and list public data breaches",
		AuthNeeded:  false,
		Commands:    []string{"alpha security hibp password [password]", "alpha security hibp breaches"},
	},

	// Productivity - Auth Required (Google OAuth)
	{
		ID:          "gdrive",
		Name:        "Google Drive",
		Group:       "productivity",
		Description: "Full read/write access to Google Drive: search, read, upload, download, create folders, update, and delete files",
		AuthNeeded:  true,
		Commands:    []string{"alpha productivity gdrive search [query]", "alpha productivity gdrive info [file-id]", "alpha productivity gdrive list", "alpha productivity gdrive read [file-id]", "alpha productivity gdrive download [file-id]", "alpha productivity gdrive upload [path]", "alpha productivity gdrive mkdir [name]", "alpha productivity gdrive update [file-id]", "alpha productivity gdrive delete [file-id]"},
		SetupCmd:    "alpha setup show google-oauth",
	},
	{
		ID:          "gsheets",
		Name:        "Google Sheets",
		Group:       "productivity",
		Description: "Full read/write access to Google Sheets: read, write, append, search, clear, and create spreadsheets",
		AuthNeeded:  true,
		Commands:    []string{"alpha productivity gsheets get [spreadsheet-id]", "alpha productivity gsheets read [spreadsheet-id] [range]", "alpha productivity gsheets search [spreadsheet-id] [query]", "alpha productivity gsheets write [spreadsheet-id] [range] [values]", "alpha productivity gsheets append [spreadsheet-id] [range] [values]", "alpha productivity gsheets clear [spreadsheet-id] [range]", "alpha productivity gsheets create --title [name]"},
		SetupCmd:    "alpha setup show google-oauth",
	},

	// Utility - No Auth (additional)
	{
		ID:          "geocoding",
		Name:        "Geocoding",
		Group:       "utility",
		Description: "Forward and reverse geocoding (address to coordinates and back)",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility geocode forward [address]", "alpha utility geocode reverse [lat] [lon]"},
	},
	{
		ID:          "timezone",
		Name:        "Timezone",
		Group:       "utility",
		Description: "Get time in timezones, lookup timezone by IP, list all timezones",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility timezone get [timezone]", "alpha utility timezone ip [ip]", "alpha utility timezone list"},
	},
	{
		ID:          "paste",
		Name:        "Paste",
		Group:       "utility",
		Description: "Create and fetch text pastes (pastebin-like)",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility paste create [content]", "alpha utility paste get [url]"},
	},
	{
		ID:          "netdiag",
		Name:        "Network Diagnostics",
		Group:       "utility",
		Description: "HTTP headers, port scanning, and DNS/ping diagnostics",
		AuthNeeded:  false,
		Commands:    []string{"alpha utility netdiag headers [url]", "alpha utility netdiag ports [host]", "alpha utility netdiag ping [host]"},
	},
}

func NewIntegrationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "integrations",
		Aliases: []string{"int", "services"},
		Short:   "List all available integrations",
	}

	cmd.AddCommand(newIntListCmd())
	cmd.AddCommand(newIntReadyCmd())
	cmd.AddCommand(newIntGroupCmd())

	return cmd
}

func newIntListCmd() *cobra.Command {
	var noAuth bool
	var group string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all integrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := make([]Integration, 0)

			for i := range allIntegrations {
				integ := allIntegrations[i]
				// Filter by auth requirement
				if noAuth && integ.AuthNeeded {
					continue
				}

				// Filter by group
				if group != "" && integ.Group != group {
					continue
				}

				// Set status
				integ.Status = getIntegrationStatus(integ)
				result = append(result, integ)
			}

			return output.Print(result)
		},
	}

	cmd.Flags().BoolVar(&noAuth, "no-auth", false, "Only show integrations that don't need authentication")
	cmd.Flags().StringVarP(&group, "group", "g", "", "Filter by group: news, knowledge, utility, dev, social, comms, productivity, system, security, marketing")

	return cmd
}

func newIntReadyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ready",
		Short: "List integrations ready to use (configured or no auth needed)",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := make([]Integration, 0)

			for i := range allIntegrations {
				integ := allIntegrations[i]
				status := getIntegrationStatus(integ)
				if status == statusReady || status == statusNoAuth {
					integ.Status = status
					result = append(result, integ)
				}
			}

			return output.Print(result)
		},
	}

	return cmd
}

func newIntGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "List integration groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			groups := map[string]struct {
				Name  string `json:"name"`
				Desc  string `json:"desc"`
				Count int    `json:"count"`
			}{
				"news":         {Name: "News", Desc: "News feeds and articles", Count: 0},
				"knowledge":    {Name: "Knowledge", Desc: "Research and reference", Count: 0},
				"utility":      {Name: "Utility", Desc: "Weather, tools", Count: 0},
				"dev":          {Name: "Dev", Desc: "Developer tools and package registries", Count: 0},
				"social":       {Name: "Social", Desc: "Social media platforms", Count: 0},
				"comms":        {Name: "Comms", Desc: "Email and messaging", Count: 0},
				"productivity": {Name: "Productivity", Desc: "Calendar, tasks, notes", Count: 0},
				"system":       {Name: "System", Desc: "macOS system integrations", Count: 0},
				"security":     {Name: "Security", Desc: "Security scanning and threat intelligence", Count: 0},
				"marketing":    {Name: "Marketing", Desc: "Ad platforms and marketing tools", Count: 0},
			}

			for i := range allIntegrations {
				integ := allIntegrations[i]
				if g, ok := groups[integ.Group]; ok {
					g.Count++
					groups[integ.Group] = g
				}
			}

			type GroupInfo struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Desc  string `json:"desc"`
				Count int    `json:"count"`
			}

			result := []GroupInfo{
				{ID: "news", Name: groups["news"].Name, Desc: groups["news"].Desc, Count: groups["news"].Count},
				{ID: "knowledge", Name: groups["knowledge"].Name, Desc: groups["knowledge"].Desc, Count: groups["knowledge"].Count},
				{ID: "utility", Name: groups["utility"].Name, Desc: groups["utility"].Desc, Count: groups["utility"].Count},
				{ID: "dev", Name: groups["dev"].Name, Desc: groups["dev"].Desc, Count: groups["dev"].Count},
				{ID: "social", Name: groups["social"].Name, Desc: groups["social"].Desc, Count: groups["social"].Count},
				{ID: "comms", Name: groups["comms"].Name, Desc: groups["comms"].Desc, Count: groups["comms"].Count},
				{ID: "productivity", Name: groups["productivity"].Name, Desc: groups["productivity"].Desc, Count: groups["productivity"].Count},
				{ID: "system", Name: groups["system"].Name, Desc: groups["system"].Desc, Count: groups["system"].Count},
				{ID: "security", Name: groups["security"].Name, Desc: groups["security"].Desc, Count: groups["security"].Count},
				{ID: "marketing", Name: groups["marketing"].Name, Desc: groups["marketing"].Desc, Count: groups["marketing"].Count},
			}

			return output.Print(result)
		},
	}

	return cmd
}

//nolint:gocyclo,gocritic // complex but clear sequential logic; Integration is read-only value type
func getIntegrationStatus(integ Integration) string {
	if !integ.AuthNeeded {
		return statusNoAuth
	}

	// Check if required keys are set
	switch integ.ID {
	case "github":
		if v, _ := config.Get("github_token"); v != "" {
			return statusReady
		}
	case "gitlab":
		if v, _ := config.Get("gitlab_token"); v != "" {
			return statusReady
		}
	case "linear":
		if v, _ := config.Get("linear_token"); v != "" {
			return statusReady
		}
	case "twitter":
		if v, _ := config.Get("x_client_id"); v != "" {
			return statusReady
		}
	case "reddit":
		if v, _ := config.Get("reddit_client_id"); v != "" {
			return statusReady
		}
	case "mastodon":
		if v, _ := config.Get("mastodon_token"); v != "" {
			return statusReady
		}
	case "youtube":
		if v, _ := config.Get("youtube_api_key"); v != "" {
			return statusReady
		}
	case "email":
		addr, _ := config.Get("email_address")
		pass, _ := config.Get("email_password")
		imap, _ := config.Get("imap_server")
		smtp, _ := config.Get("smtp_server")
		if addr != "" && pass != "" && imap != "" && smtp != "" {
			return statusReady
		}
	case "slack":
		if v, _ := config.Get("slack_token"); v != "" {
			return statusReady
		}
	case "discord":
		if v, _ := config.Get("discord_token"); v != "" {
			return statusReady
		}
	case "telegram":
		if v, _ := config.Get("telegram_token"); v != "" {
			return statusReady
		}
	case "twilio":
		sid, _ := config.Get("twilio_sid")
		token, _ := config.Get("twilio_token")
		phone, _ := config.Get("twilio_phone")
		if sid != "" && token != "" && phone != "" {
			return statusReady
		}
	case "calendar":
		cid, _ := config.Get("google_client_id")
		secret, _ := config.Get("google_client_secret")
		refresh, _ := config.Get("google_refresh_token")
		if cid != "" && secret != "" && refresh != "" {
			return statusReady
		}
	case "notion":
		if v, _ := config.Get("notion_token"); v != "" {
			return statusReady
		}
	case "todoist":
		if v, _ := config.Get("todoist_token"); v != "" {
			return statusReady
		}
	case "newsapi":
		if v, _ := config.Get("newsapi_key"); v != "" {
			return statusReady
		}
	case "stocks":
		if v, _ := config.Get("alphavantage_key"); v != "" {
			return statusReady
		}
	case "jira":
		url, _ := config.Get("jira_url")
		email, _ := config.Get("jira_email")
		token, _ := config.Get("jira_token")
		if url != "" && email != "" && token != "" {
			return statusReady
		}
	case "cloudflare":
		if v, _ := config.Get("cloudflare_token"); v != "" {
			return statusReady
		}
	case "vercel":
		if v, _ := config.Get("vercel_token"); v != "" {
			return statusReady
		}
	case "trello":
		key, _ := config.Get("trello_key")
		token, _ := config.Get("trello_token")
		if key != "" && token != "" {
			return statusReady
		}
	case "logseq":
		if v, _ := config.Get("logseq_graph"); v != "" {
			return statusReady
		}
	case "obsidian":
		if v, _ := config.Get("obsidian_vault"); v != "" {
			return statusReady
		}
	case "facebook-ads":
		token, _ := config.Get("facebook_ads_token")
		acctID, _ := config.Get("facebook_ads_account_id")
		if token != "" && acctID != "" {
			return statusReady
		}
	case "amazon-sp":
		cid, _ := config.Get("amazon_sp_client_id")
		secret, _ := config.Get("amazon_sp_client_secret")
		refresh, _ := config.Get("amazon_sp_refresh_token")
		seller, _ := config.Get("amazon_sp_seller_id")
		if cid != "" && secret != "" && refresh != "" && seller != "" {
			return statusReady
		}
	case "shopify":
		store, _ := config.Get("shopify_store")
		token, _ := config.Get("shopify_token")
		if store != "" && token != "" {
			return statusReady
		}
	case "spotify":
		cid, _ := config.Get("spotify_client_id")
		secret, _ := config.Get("spotify_client_secret")
		if cid != "" && secret != "" {
			return statusReady
		}
	case "sentry":
		if v, _ := config.Get("sentry_auth_token"); v != "" {
			return statusReady
		}
	case "s3":
		profile, _ := config.Get("aws_profile")
		region, _ := config.Get("aws_region")
		if profile != "" && region != "" {
			return statusReady
		}
	case "redis":
		if v, _ := config.Get("redis_url"); v != "" {
			return statusReady
		}
	case "prometheus":
		if v, _ := config.Get("prometheus_url"); v != "" {
			return statusReady
		}
	case "virustotal":
		if v, _ := config.Get("virustotal_api_key"); v != "" {
			return statusReady
		}
	case "gdrive":
		cid, _ := config.Get("google_client_id")
		secret, _ := config.Get("google_client_secret")
		refresh, _ := config.Get("google_refresh_token")
		if cid != "" && secret != "" && refresh != "" {
			return statusReady
		}
	case "gsheets":
		cid, _ := config.Get("google_client_id")
		secret, _ := config.Get("google_client_secret")
		refresh, _ := config.Get("google_refresh_token")
		if cid != "" && secret != "" && refresh != "" {
			return statusReady
		}
	}

	return "needs_setup"
}
