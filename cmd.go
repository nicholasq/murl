package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().BoolVar(&cfg.Summarize, "s", false, "Enable AI Summarization")
	rootCmd.Flags().BoolVarP(&cfg.Summarize, "summarize", "s", false, "Enable AI Summarization")
	rootCmd.Flags().StringVarP(&cfg.Model, "model", "m", "", "Ollama model to use (default: llama3.2)")
}

var cfg Config

var rootCmd = &cobra.Command{
	Use:   "murl",
	Short: "Fetch markdown from a url",
	Long:  "Fetch page from url and convert the body into markdown",
	Example: `murl https://example.com
murl https://example.com -s
	`,
	RunE: runQuery,
}

type Config struct {
	Summarize bool
	Model     string
}

func runQuery(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Usage()
		os.Exit(1)
	}

	url := args[0]

	markdown, err := MarkdownFromUrl(url)

	if err != nil {
		return fmt.Errorf("Failed to convert HTML to Markdown: %w", err)
	}

	if cfg.Summarize {
		summary, err := summarizeContent(markdown, cfg.Model)
		if err != nil {
			return fmt.Errorf("Failed to summarize content: %w", err)
		}
		fmt.Println(summary)
		return nil
	}

	fmt.Println(markdown)
	return nil
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func summarizeContent(markdown string, model string) (string, error) {
	if model == "" {
		model = "llama3.2"
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`Analyze the following content and provide a detailed summary following this format:

Summary:
[High level overview of the content]

Key Points:
- [Key point 1]
  - [Supporting detail]
  - [Supporting detail]
- [Key point 2]
  - [Supporting detail]
  - [Supporting detail]

Content to summarize:

%s`, markdown)

	req := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: new(bool),
	}

	var summary string
	ctx := context.Background()

	err = client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		summary += resp.Response
		return nil
	})

	if err != nil {
		return "", err
	}

	return summary, nil
}
