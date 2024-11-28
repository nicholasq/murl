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

### AI Summarization

To use the AI summarization feature:
1. Install Ollama locally: https://ollama.ai
2. Pull your preferred model: `ollama pull llama3.2`
3. Use the `-s` or `--summarize` flag when running murl

Default model is llama3.2, but you can specify others with the `-m` flag.

## Usage

```bash
murl https://example.com
murl https://example.com -s     # Generate AI summary
murl https://example.com --summarize --model mistral   # Use specific model for summary
```
