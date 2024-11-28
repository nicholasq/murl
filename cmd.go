package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().BoolVar(&cfg.Summarize, "s", false, "Enable AI Summarization")
}

var cfg Config

var rootCmd = &cobra.Command{
	Use:   "murl",
	Short: "Fetch markdown from a url",
	Long:  "Fetch page from url and convert the body into markdown",
	Example: ` murl https://example.com
	murl https://example.com -s
	`,
	RunE: runQuery,
}

type Config struct {
	Summarize bool
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

	fmt.Println(markdown)
	return nil
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
