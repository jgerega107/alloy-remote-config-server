# alloy-remote-config-server

![build workflow](https://github.com/jgerega107/alloy-remote-config-server/actions/workflows/docker-publish.yml/badge.svg)
![license](http://img.shields.io/badge/license-MIT-blue.svg)

## Description

A lightweight, secure remote configuration server for [Grafana Alloy](https://grafana.com/docs/alloy/). Implements the Alloy remote config API to serve dynamic configurations via gRPC with template-based rendering.

See: [Grafana Alloy remote config API](https://github.com/grafana/alloy-remote-config)

### Features

- **Template-based configuration** - Use Go templates with `.conf.tmpl` extension for dynamic configs
- **gRPC API** - Implements Alloy's `remotecfg` protocol for agent communication
- **REST API** - HTTP endpoints for health checks and template discovery
- **Zero dependencies** - Uses only Go standard library and essential gRPC packages
- **Secure by default** - Runs as non-root user in minimal Alpine container
- **Multi-arch support** - Docker images for AMD64 and ARM64

## How It Works

The server renders Alloy configurations on-the-fly using Go templates:

1. **Template files** (`.conf.tmpl`) define configuration patterns with variables
2. **Alloy agents** connect via gRPC and send their `id` and `local_attributes`
3. The server matches the `template` attribute to a template file (defaults to `default`)
4. Template variables are substituted:
   - `{{ .Id }}` - Agent's unique identifier
   - `{{ .Attributes }}` - Map of all `local_attributes` sent by the agent

### Example Template

```alloy
// File: conf/default.conf.tmpl
prometheus.scrape "{{ .Id }}" {
  targets = ["localhost:9090"]
  forward_to = [prometheus.remote_write.default.receiver]
}

{{ if eq .Attributes.env "prod" }}
// Production-specific settings
{{ end }}
```

### Example Alloy Agent Config

```alloy
remotecfg {
  url = "http://your-server:8888"
  id = "web-server-01"
  attributes = {
    "template" = "web-server",  // Selects web-server.conf.tmpl
    "env" = "prod",
    "region" = "us-east-1"
  }
}
```

### Multiple Templates

You can combine multiple templates by specifying a comma-separated list in the `templates` attribute (note the plural):

```alloy
remotecfg {
  url = "http://your-server:8888"
  id = "web-server-01"
  attributes = {
    "templates" = "base,web-server,monitoring"  // Combines 3 templates
    "env" = "prod",
    "region" = "us-east-1"
  }
}
```

Templates are rendered in order and concatenated with blank line separators. This allows you to:
- Compose configurations from reusable components
- Layer settings (base + environment-specific + service-specific)
- Keep templates DRY (Don't Repeat Yourself)

**Note:** You can still use the singular `template` attribute for single template selection.

## Configuration

Environment variables (all optional):

| Variable | Default | Description |
|----------|---------|-------------|
| `CONFIG_FOLDER` | `conf` | Directory containing template files (`.conf.tmpl`) |
| `GRPC_PORT` | `8888` | gRPC port for Alloy agent connections |
| `HTTP_PORT` | `8080` | HTTP REST API port |
| `LISTEN_ADDR` | `0.0.0.0` | Network interface to bind |

## API Endpoints

### REST API (HTTP)

- `GET /health` - Health check endpoint
- `GET /templates` - List available template names

### gRPC API

Implements the Alloy CollectorService protocol on port 8888.

### Local Development

```bash
# Ensure templates exist
mkdir -p conf
cat > conf/default.conf.tmpl << 'EOF'
prometheus.scrape "{{ .Id }}" {
  targets = ["localhost:9090"]
}
EOF

# Run
go mod tidy
go run cmd/config/main.go
```

## Template Directory

Create a directory (default: `conf/`) containing template files with `.conf.tmpl` extension:

```
conf/
├── default.conf.tmpl      # Required - fallback template
├── web-server.conf.tmpl   # Optional - for web servers
└── database.conf.tmpl     # Optional - for databases
```

**Note:** Template names must be unique. If multiple files resolve to the same name (e.g., `web-server.conf.tmpl` in different subdirectories), the last one loaded will overwrite the previous one and a warning will be logged.

## Security

- Runs as non-root user (`appuser` UID/GID)
- Minimal Alpine Linux base image
- Read-only filesystem support
- No external database dependencies
- For TLS/mTLS/OAuth2, deploy behind a reverse proxy (e.g., Istio, Nginx, Traefik)

## Architecture

- **Language**: Go 1.23+
- **Dependencies**: Minimal - only `connectrpc.com/connect`, `github.com/grafana/alloy-remote-config`, and `golang.org/x/net`
- **Storage**: In-memory template cache (no persistent storage)
- **Container**: Multi-stage build with distroless Alpine

## Building

```bash
# Local build
go build -ldflags="-w -s" -trimpath -o alloy-remote-config-server cmd/config/main.go

# Docker build
docker build -t alloy-remote-config-server .

# Multi-arch build
docker buildx build --platform linux/amd64,linux/arm64 -t alloy-remote-config-server .
```

## License

MIT License - see LICENSE file for details
