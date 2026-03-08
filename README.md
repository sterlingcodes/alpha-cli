# 🛠️ Alpha CLI

<p align="center">
  <img src="https://raw.githubusercontent.com/sterlingcodes/alpha-cli/main/assets/icon_rounded_1024.png" alt="Alpha CLI" width="200">
</p>

<p align="center">
  <strong>Give your AI assistant hands to interact with the internet.</strong>
</p>

<p align="center">
  <a href="https://github.com/sterlingcodes/alpha-cli/releases/latest"><img src="https://img.shields.io/github/v/release/sterlingcodes/alpha-cli?include_prereleases&style=for-the-badge" alt="GitHub release"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" alt="MIT License"></a>
  <a href="https://youtube.com/@kenkaidoesai"><img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" alt="YouTube"></a>
  <a href="https://skool.com/kenkai"><img src="https://img.shields.io/badge/Skool-Community-7C3AED?style=for-the-badge" alt="Skool"></a>
</p>

**Alpha CLI** gives your AI assistant the power to actually *do things* on the internet — check emails, browse social media, get news, look up information, and more.

Think of it as hands for your AI. Instead of just chatting, your AI can now reach out and interact with real services like Twitter, YouTube, Hacker News, Wikipedia, and dozens more.

No coding required. Just install it, and your AI assistant instantly gains superpowers to help you with real tasks across the web.

---

## 🚀 Install

One command. That's it.

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/sterlingcodes/alpha-cli/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/sterlingcodes/alpha-cli/main/scripts/install.ps1 | iex
```

Works on **macOS** (Intel & Apple Silicon), **Linux**, and **Windows**.

The installer automatically:
- Downloads the right version for your system
- Installs it globally
- Configures your PATH
- Works immediately in new terminals

To update later, just run the same command again.

---

## 🧠 Why this exists

AI assistants are smart but powerless. They can answer questions, but they can't actually *do* anything.

Alpha CLI changes that. It's a universal interface that lets any AI agent interact with the real world:
- Check your emails and send replies
- Send SMS via Twilio, push notifications via ntfy
- Message on Slack, Discord, Telegram
- Search YouTube, get video stats
- Browse Hacker News, Reddit, Twitter
- Look up weather, crypto prices, currency rates
- Query Wikipedia, StackOverflow, dictionaries
- Manage Todoist tasks, Notion pages, Obsidian vaults
- Control macOS apps: Calendar, Reminders, Notes, Contacts, Finder, Safari
- **80 integrations** across 10 categories

All with simple commands that return clean JSON — perfect for AI to understand and act on.

---

## ✨ What you can do

### No setup required (works immediately)
```bash
alpha news hn top -l 5              # Top 5 Hacker News stories
alpha utility weather now "Tokyo"   # Current weather in Tokyo
alpha knowledge wiki summary "AI"   # Wikipedia summary
alpha utility crypto price bitcoin  # Bitcoin price
alpha utility currency convert 100 USD EUR  # Currency conversion
alpha utility translate text "Hello" --to es # Translate to Spanish
alpha dev npm info react            # npm package info
alpha dev dockerhub search nginx    # Search Docker images
alpha comms notify ntfy mytopic "Hello!"    # Push notification (no auth)
alpha comms webhook slack [url] "Message"   # Slack webhook
alpha security crtsh lookup example.com      # Certificate transparency logs
alpha utility netdiag ping example.com      # Network diagnostics
alpha utility timezone get "America/New_York" # Timezone info
alpha utility paste create "code snippet"   # Create a paste

# macOS only (no auth needed)
alpha system reminders today        # Today's reminders
alpha system notes list             # List Apple Notes
alpha system calendar today         # Today's calendar events
alpha system clipboard get          # Get clipboard content
alpha system finder search "query"  # Spotlight search
```

### With credentials (one-time setup)
```bash
alpha comms email list -l 10        # Your latest emails
alpha comms slack channels          # List Slack channels
alpha comms discord guilds          # List Discord servers
alpha comms telegram send 123 "Hi"  # Send Telegram message
alpha comms twilio send +1234 "SMS" # Send SMS via Twilio
alpha social youtube search "AI"    # Search YouTube
alpha social twitter timeline       # Your Twitter feed
alpha social spotify search "Lo-fi"  # Search Spotify tracks
alpha productivity todoist tasks    # Your todo list
alpha productivity trello boards    # Your Trello boards
alpha productivity obsidian notes   # List Obsidian notes
alpha productivity logseq pages     # List Logseq pages
alpha productivity gdrive search "q" # Google Drive files
alpha productivity gsheets read ID  # Read a Google Sheet
alpha dev github repos              # Your GitHub repos
alpha dev jira issues               # Your Jira issues
alpha dev sentry issues             # Sentry error tracking
alpha dev kube pods                 # Kubernetes pods
alpha security virustotal url URL   # Scan URL for threats
alpha security shodan lookup 1.2.3.4 # Shodan host lookup
```

---

## 🔧 Quick start

### See what's available
```bash
alpha commands                      # All commands (for AI agents)
alpha integrations list             # All integrations + auth status
alpha integrations list --no-auth   # Services that work without setup
```

### Set up credentials
```bash
alpha setup list                    # What needs configuration
alpha setup show email              # Step-by-step setup guide
alpha setup set email imap_server imap.gmail.com
```

### Example workflow
```bash
# Check what integrations work without auth
$ alpha integrations list --no-auth

# Get top tech news
$ alpha news hn top -l 3

# Look up a term
$ alpha knowledge dict define "API"

# Check the weather
$ alpha utility weather now "San Francisco"

# Send yourself a notification (no auth needed!)
$ alpha comms notify ntfy my-alerts "Task completed!"
```

### Communication examples
```bash
# Send SMS (requires Twilio setup)
alpha comms twilio send "+15551234567" "Hello from Alpha CLI"

# Discord bot commands
alpha comms discord guilds              # List servers
alpha comms discord channels 123456     # List channels in server
alpha comms discord send 789 "Hello!"   # Send message to channel

# Slack integration
alpha comms slack channels              # List channels
alpha comms slack send general "Hi!"    # Post to channel
alpha comms slack search "important"    # Search messages

# Telegram bot
alpha comms telegram chats              # List chats
alpha comms telegram send 123 "Hello"   # Send message

# Push notifications (ntfy.sh - no auth!)
alpha comms notify ntfy alerts "Server is down!" --priority 5

# Webhooks (no auth)
alpha comms webhook discord [url] "Deployment complete"
```

### macOS system examples (no auth needed)
```bash
# Apple Reminders
alpha system reminders lists            # List all reminder lists
alpha system reminders today            # Today's reminders
alpha system reminders add "Buy milk"   # Add a reminder
alpha system reminders complete "Buy milk"  # Mark complete

# Apple Notes
alpha system notes folders              # List folders
alpha system notes list                 # List all notes
alpha system notes read "Shopping"      # Read a note
alpha system notes create "Ideas" "My brilliant idea"

# Apple Calendar
alpha system apple-calendar calendars   # List calendars
alpha system apple-calendar today       # Today's events
alpha system apple-calendar upcoming    # Next 7 days

# Apple Contacts
alpha system contacts search "John"     # Search contacts
alpha system contacts get "John Doe"    # Get full details

# Finder & Clipboard
alpha system finder search "project"    # Spotlight search
alpha system finder info ~/Documents    # Get folder info
alpha system clipboard get              # Get clipboard
alpha system clipboard set "Hello"      # Set clipboard

# Safari (requires Safari to be running)
alpha system safari tabs                # List open tabs
alpha system safari bookmarks           # List bookmarks
alpha system safari history --limit 10  # Recent history
```

### Obsidian & Logseq examples
```bash
# Obsidian (configure vault path first)
alpha config set obsidian_vault ~/Documents/MyVault
alpha productivity obsidian notes       # List all notes
alpha productivity obsidian daily       # Today's daily note
alpha productivity obsidian search "AI" # Search notes
alpha productivity obsidian read "Ideas"  # Read a note

# Logseq (configure graph path first)
alpha config set logseq_graph ~/Documents/MyGraph
alpha productivity logseq pages         # List pages
alpha productivity logseq journal       # Today's journal
alpha productivity logseq search "todo" # Search pages
```

---

## 📦 All 80 integrations

| Category | Services |
|----------|----------|
| **Social** (5) | Twitter/X, Reddit, Mastodon, YouTube, Spotify |
| **Communication** (7) | Email (IMAP/SMTP), Slack, Discord, Telegram, Twilio SMS, Push Notifications (ntfy/Pushover), Webhooks |
| **News** (3) | Hacker News, RSS feeds, NewsAPI |
| **Knowledge** (3) | Wikipedia, StackOverflow, Dictionary |
| **Dev Tools** (16) | GitHub, GitLab, Gist, Linear, Jira, Sentry, Cloudflare, Vercel, npm, PyPI, Docker Hub, Redis, Prometheus, Kubernetes, Database, S3 |
| **Productivity** (8) | Todoist, Notion, Google Calendar, Google Drive, Google Sheets, Trello, Obsidian, Logseq |
| **Utility** (18) | Weather, Crypto, Currency, IP lookup, DNS/WHOIS/SSL, Wayback Machine, Holidays, Translation, URL Shortener, Stocks, Geocoding, Network Diagnostics, Pastebin, Timezone, DNS Benchmark, Speed Test, Traceroute, WiFi Info |
| **Security** (4) | VirusTotal, Shodan, Certificate Transparency (crt.sh), Have I Been Pwned |
| **Marketing** (3) | Facebook Ads (Meta), Amazon Selling Partner, Shopify |
| **System** (13) | Apple Calendar, Apple Reminders, Apple Notes, Apple Contacts, Apple Mail, Safari, Finder, Clipboard, iMessage, Battery, System Cleanup, Disk Health, System Info *(macOS only)* |

### 46 integrations work without any setup:
Hacker News, RSS, Wikipedia, StackOverflow, Dictionary, Weather, Crypto, Currency, IP lookup, Domain tools, Wayback Machine, Holidays, Translation, URL Shortener, npm, PyPI, Docker Hub, Gist, Kubernetes, Database, Geocoding, Timezone, Network Diagnostics, Pastebin, DNS Benchmark, Speed Test, Traceroute, WiFi Info, Shodan, Certificate Transparency, Have I Been Pwned, ntfy notifications, Webhooks, plus all 13 macOS System integrations

---

## 🤖 Built for AI agents

Every command outputs clean JSON:

```json
{
  "success": true,
  "data": {
    "title": "Show HN: I built a CLI for AI agents",
    "score": 142,
    "url": "https://..."
  }
}
```

Errors are structured too:

```json
{
  "success": false,
  "error": {
    "code": "setup_required",
    "message": "Email not configured",
    "setup_cmd": "alpha setup show email"
  }
}
```

Your AI knows exactly what went wrong and how to fix it.

---

## 🔒 Privacy

- Credentials stored locally in `~/.config/alpha/config.json`
- No telemetry, no analytics
- API calls go directly to the services you configure
- Open source — inspect every line

---

## 🛠️ For developers

```bash
git clone https://github.com/sterlingcodes/alpha-cli.git
cd alpha-cli
make install
```

Build releases for all platforms:
```bash
make release
```

Stack: Go + Cobra CLI + zero external dependencies at runtime

---

## 👥 Community

- [YouTube @kenkaidoesai](https://youtube.com/@kenkaidoesai) — tutorials and demos
- [Skool community](https://skool.com/kenkai) — come hang out

---

## 📄 License

MIT

---

<p align="center">
  <strong>Give your AI the power to actually do things.</strong>
</p>

<p align="center">
  <a href="https://github.com/sterlingcodes/alpha-cli/releases/latest"><img src="https://img.shields.io/badge/Install-One%20Command-blue?style=for-the-badge" alt="Install"></a>
</p>
