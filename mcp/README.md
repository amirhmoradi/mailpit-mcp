# Mailpit MCP Server - Source Code

This directory contains the source code for the Mailpit MCP server.

For full documentation, installation instructions, and usage examples, see the [main README](../README.md).

## Directory Structure

```
mcp/
├── cmd/
│   └── mailpit-mcp-server/   # Main entry point
├── internal/
│   ├── client/               # Mailpit HTTP API client
│   ├── server/               # MCP server setup and config
│   ├── tools/                # MCP tool implementations
│   ├── resources/            # MCP resource implementations
│   └── prompts/              # MCP prompt templates
├── Dockerfile                # Docker build configuration
├── go.mod                    # Go module definition
└── go.sum                    # Go dependencies lock file
```

## Quick Build

```bash
# Build the binary
go build -o mailpit-mcp-server ./cmd/mailpit-mcp-server

# Run tests
go test ./...

# Run locally
MAILPIT_URL=http://localhost:8025 ./mailpit-mcp-server
```

## Docker Build

```bash
# Build the Docker image
docker build -t mailpit-mcp:local .

# Run the container
docker run -d \
  --name mailpit-mcp \
  -e MAILPIT_URL=http://host.docker.internal:8025 \
  -e MCP_TRANSPORT=http \
  -p 3000:3000 \
  mailpit-mcp:local
```

## Available MCP Tools

### Message Management
- `list_messages` - List messages with pagination
- `search_messages` - Search using Mailpit query syntax
- `get_message` - Get full message details
- `get_message_headers` - Get message headers
- `get_message_source` - Get raw RFC 2822 source
- `delete_messages` - Delete messages by ID
- `delete_search` - Delete messages matching a search
- `set_read_status` - Mark messages as read/unread

### Content Tools
- `get_message_html` - Get rendered HTML content
- `get_message_text` - Get plain text content
- `get_attachment` - Download an attachment

### Validation Tools
- `check_html` - Check HTML email client compatibility
- `check_links` - Validate links in an email
- `check_spam` - Get SpamAssassin analysis

### Tag Management
- `list_tags` - List all tags
- `set_tags` - Set tags on messages
- `rename_tag` - Rename a tag
- `delete_tag` - Delete a tag

### Testing Tools
- `send_message` - Send a test email
- `release_message` - Relay to external SMTP
- `get_chaos` - Get chaos testing configuration
- `set_chaos` - Configure chaos testing triggers

### System Tools
- `get_info` - Get server information and stats
- `get_webui_config` - Get Mailpit configuration

## MCP Resources

| Resource URI | Description |
|--------------|-------------|
| `mailpit://info` | Server information (JSON) |
| `mailpit://messages/latest` | Latest message summary |
| `mailpit://tags` | All current tags (JSON) |
| `mailpit://config` | Mailpit configuration (JSON) |

## MCP Prompts

| Prompt | Description |
|--------|-------------|
| `analyze_latest_email` | Comprehensive email analysis |
| `debug_email_delivery` | Debug delivery issues |
| `check_email_quality` | Full quality check (HTML, links, spam) |
| `search_emails` | Help construct search queries |
| `compose_test_email` | Create and send test emails |
| `analyze_email_headers` | Deep header analysis |
| `compare_emails` | Compare two emails (A/B testing) |
| `summarize_inbox` | Summarize inbox state |

## License

MIT License
