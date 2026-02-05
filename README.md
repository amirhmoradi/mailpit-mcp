# Mailpit MCP Server

A standalone [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that enables AI assistants to interact with [Mailpit](https://mailpit.axllent.org/) email testing environments.

## Overview

This project provides an MCP server that connects AI assistants (Claude, Cursor, VS Code Copilot, etc.) to Mailpit's API, enabling natural language interaction with your email testing workflow.

**Key Features:**
- Browse, search, and analyze emails using natural language
- Check HTML compatibility across email clients
- Validate links and detect broken images
- Run spam analysis on messages (requires SpamAssassin)
- Send test emails and manage tags
- Debug email delivery issues

## Prerequisites

- A running [Mailpit](https://mailpit.axllent.org/) instance (the MCP server connects to Mailpit via its HTTP API)
- Go 1.24+ (if building from source)

## Installation

### Option 1: Docker (Recommended)

```bash
docker run -d \
  --name mailpit-mcp \
  -e MAILPIT_URL=http://your-mailpit-host:8025 \
  -e MCP_TRANSPORT=http \
  -p 3000:3000 \
  ghcr.io/amirhmoradi/mailpit-mcp:latest
```

### Option 2: Docker Compose

Create a `docker-compose.yml`:

```yaml
services:
  mailpit:
    image: axllent/mailpit:latest
    ports:
      - "8025:8025"
      - "1025:1025"

  mailpit-mcp:
    image: ghcr.io/amirhmoradi/mailpit-mcp:latest
    environment:
      MAILPIT_URL: http://mailpit:8025
      MCP_TRANSPORT: http
    ports:
      - "3000:3000"
    depends_on:
      - mailpit
```

Run with:
```bash
docker-compose up -d
```

### Option 3: Build from Source

```bash
# Clone the repository
git clone https://github.com/amirhmoradi/mailpit-mcp.git
cd mailpit-mcp

# Build the binary
cd mcp
go build -o mailpit-mcp-server ./cmd/mailpit-mcp-server

# Run the server
MAILPIT_URL=http://localhost:8025 ./mailpit-mcp-server
```

### Option 4: Go Install

```bash
go install github.com/amirhmoradi/mailpit-mcp/cmd/mailpit-mcp-server@latest
```

## Configuration

Configure the MCP server using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `MAILPIT_URL` | Mailpit server URL | `http://localhost:8025` |
| `MAILPIT_AUTH_USER` | Basic auth username (if Mailpit requires auth) | (none) |
| `MAILPIT_AUTH_PASS` | Basic auth password | (none) |
| `MAILPIT_TIMEOUT` | Request timeout in seconds | `30` |
| `MCP_TRANSPORT` | Transport mode: `stdio` or `http` | `stdio` |
| `MCP_HTTP_HOST` | HTTP server bind address | `0.0.0.0` |
| `MCP_HTTP_PORT` | HTTP server port | `3000` |
| `MCP_LOG_LEVEL` | Log level: `debug`, `info`, `warn`, `error` | `info` |

## Usage with AI Tools

### Claude Desktop

Add to your Claude Desktop configuration file:

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "mailpit": {
      "command": "/path/to/mailpit-mcp-server",
      "env": {
        "MAILPIT_URL": "http://localhost:8025"
      }
    }
  }
}
```

Or if using Docker:

```json
{
  "mcpServers": {
    "mailpit": {
      "command": "docker",
      "args": ["run", "-i", "--rm",
        "-e", "MAILPIT_URL=http://host.docker.internal:8025",
        "ghcr.io/amirhmoradi/mailpit-mcp:latest"]
    }
  }
}
```

### Cursor / VS Code

For IDE integration, run the MCP server in HTTP mode:

```bash
# Start the server in HTTP mode
MAILPIT_URL=http://localhost:8025 MCP_TRANSPORT=http ./mailpit-mcp-server
```

Then add to your IDE's MCP configuration:

```json
{
  "mcp.servers": {
    "mailpit": {
      "url": "http://localhost:3000/mcp"
    }
  }
}
```

## Available Tools

| Tool | Description |
|------|-------------|
| `list_messages` | List messages with pagination |
| `search_messages` | Search using Mailpit query syntax |
| `get_message` | Get full message details |
| `get_message_headers` | Get message headers |
| `get_message_source` | Get raw RFC 2822 source |
| `get_message_html` | Get rendered HTML content |
| `get_message_text` | Get plain text content |
| `get_attachment` | Download an attachment |
| `delete_messages` | Delete messages by ID |
| `delete_search` | Delete messages matching a search |
| `set_read_status` | Mark messages as read/unread |
| `check_html` | Check HTML email client compatibility |
| `check_links` | Validate links in an email |
| `check_spam` | Get SpamAssassin analysis |
| `list_tags` | List all tags |
| `set_tags` | Set tags on messages |
| `rename_tag` | Rename a tag |
| `delete_tag` | Delete a tag |
| `send_message` | Send a test email |
| `release_message` | Relay to external SMTP |
| `get_chaos` | Get chaos testing config |
| `set_chaos` | Configure chaos testing |
| `get_info` | Get server information |
| `get_webui_config` | Get Mailpit configuration |

## Available Prompts

| Prompt | Description |
|--------|-------------|
| `analyze_latest_email` | Comprehensive analysis of the most recent email |
| `debug_email_delivery` | Debug delivery issues for a message |
| `check_email_quality` | Full quality check (HTML, links, spam) |
| `search_emails` | Help construct search queries |
| `compose_test_email` | Create and send test emails |
| `analyze_email_headers` | Deep header analysis |
| `compare_emails` | Compare two emails (A/B testing) |
| `summarize_inbox` | Summarize inbox state |

## Search Syntax

The search tool supports Mailpit's query syntax:

- `from:email@example.com` - Sender address
- `to:email@example.com` - Recipient address
- `subject:keyword` - Subject contains text
- `message-id:id` - Specific Message-ID
- `tag:tagname` - Has specific tag
- `is:read` / `is:unread` - Read status
- `has:attachment` - Has attachments
- `before:YYYY-MM-DD` - Before date
- `after:YYYY-MM-DD` - After date

Multiple terms are combined with AND logic.

## Development

```bash
# Build
cd mcp
go build ./...

# Test
go test ./...

# Run locally (STDIO mode for Claude Desktop)
MAILPIT_URL=http://localhost:8025 go run ./cmd/mailpit-mcp-server

# Run locally (HTTP mode for IDE integration)
MAILPIT_URL=http://localhost:8025 MCP_TRANSPORT=http go run ./cmd/mailpit-mcp-server
```

## License

MIT License - See [LICENSE](LICENSE) for details.

## Acknowledgments

This project provides MCP integration for [Mailpit](https://mailpit.axllent.org/), the excellent email testing tool created by [@axllent](https://github.com/axllent).
