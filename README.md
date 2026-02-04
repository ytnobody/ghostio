# ghostio

GitHub repository activity monitor CLI tool.

## Overview

ghostio watches a GitHub repository and announces activities to stdout.

### Monitored Activities

- Issues (created/updated/closed)
- Pull Requests (created/merged/closed)
- Comments
- Releases

## Usage

```bash
ghostio watch owner/repo
```

## Installation

```bash
go install github.com/ytnobody/ghostio/cmd/ghostio@latest
```

## License

MIT
