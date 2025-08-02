[<a href="/README.md">English</a>] | [<a href="/docs/README.vi.md">Tiếng Việt</a>]

# CoraPanel MCP Server

**CoraPanel MCP Server** is an implementation of the Model Context Protocol (MCP) server for [CoraPanel](https://github.com/cloudora-vn/CoraPanel).

## Installation Methods

### Method 1: Download from Release Page (Recommended)

1. Visit the [Releases Page](https://github.com/cloudora-vn/mcp-corapanel/releases) and download the executable file corresponding to your system.

2. Example installation (for amd64):

```bash
chmod +x mcp-corapanel-linux-amd64
mv mcp-corapanel-linux-amd64 /usr/local/bin/mcp-corapanel
```

### Method 2: Build from Source

Make sure Go 1.23 or later is installed locally. Then run:

1. Clone the repository:

```bash
git clone https://github.com/cloudora-vn/mcp-corapanel.git
cd mcp-corapanel
```

2. Build the executable:

```bash
make build
```

> Move ./build/mcp-corapanel to a directory included in your system's PATH.

### Method 3: Install via go install

Make sure Go 1.23 or later is installed locally. Then run:

```bash
go install github.com/cloudora-vn/mcp-corapanel@latest
```

### Method 4: Install via Docker

Make sure Docker is correctly installed and configured on your machine.

The official image supports the following architectures:

- amd64
- arm64
- arm/v7
- s390x
- ppc64le

## Usage

CoraPanel MCP Server supports two running modes: `stdio` and `sse`.

### stdio Mode

#### Using Local Binary

In the configuration file of Cursor or Windsurf, add:

```json
{
  "mcpServers": {
    "mcp-corapanel": {
      "command": "mcp-corapanel",
      "env": {
        "PANEL_ACCESS_TOKEN": "<your CoraPanel access token>",
        "PANEL_HOST": "such as http://localhost:8080"
      }
    }
  }
}
```

#### Running in Docker

```json
{
  "mcpServers": {
    "mcp-corapanel": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "PANEL_HOST",
        "-e",
        "PANEL_ACCESS_TOKEN",
        "cloudora-vn/corapanel-mcp-server"
      ],
      "env": {
        "PANEL_HOST": "such as http://localhost:8080",
        "PANEL_ACCESS_TOKEN": "<your CoraPanel access token>"
      }
    }
  }
}
```

### sse Mode

1. Start the MCP Server:

```bash
mcp-corapanel -host http://localhost:8080 -token <your CoraPanel access token> -transport sse -addr http://localhost:8000
```

2. Configure in Cursor or Windsurf:

```json
{
  "mcpServers": {
    "mcp-corapanel": {
      "url": "http://localhost:8000/sse"
    }
  }
}
```

#### Command Line Options

- `-token`: CoraPanel access token
- `-host`: CoraPanel access address
- `-transport`: Transport type (stdio or sse, default: stdio)
- `-addr`: Start SSE server address (default: http://localhost:8000)

## Available Tools

The server provides various tools for interacting with CoraPanel:

| Tool                        | Category     | Description               |
|-----------------------------|--------------|---------------------------|
| **get_dashboard_info**      | System       | List dashboard status     |
| **get_system_info**         | System       | Get system information    |
| **list_websites**           | Website      | List all websites         |
| **create_website**          | Website      | Create a website          |
| **list_ssls**               | Certificate  | List all certificates     |
| **create_ssl**              | Certificate  | Create a certificate      |
| **list_installed_apps**     | Application  | List installed apps       |
| **install_openresty**       | Application  | Install OpenResty         |
| **install_mysql**           | Application  | Install MySQL             |
| **list_databases**          | Database     | List all databases        |
| **create_database**         | Database     | Create a database         |
