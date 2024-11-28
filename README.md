# MURL

Download webpages as markdown files.

## Setup

Start the headless Chrome container:

```bash
docker run -p 3000:3000 ghcr.io/browserless/chromium
```

## Build

```bash
go build
```

## Usage

```bash
murl https://example.com
```
