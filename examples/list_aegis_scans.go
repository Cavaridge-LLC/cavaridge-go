// Quickstart: list AEGIS scans.
//
// Build + run (against sandbox):
//   CAVARIDGE_API_KEY=sk_test_... go run examples/list_aegis_scans.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Cavaridge-LLC/cavaridge-go/cavaridge"
)

func main() {
	apiKey := os.Getenv("CAVARIDGE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "CAVARIDGE_API_KEY env var required")
		os.Exit(1)
	}
	baseURL := os.Getenv("CAVARIDGE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api-sandbox.cavaridge.app"
	}
	client, err := cavaridge.NewClient(
		cavaridge.WithAPIKey(apiKey),
		cavaridge.WithBaseURL(baseURL),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "FAIL:", err)
		os.Exit(1)
	}
	page, err := client.Aegis().ListScans(context.Background(), struct {
		ClientID      string
		Limit         int
		StartingAfter string
	}{Limit: 5})
	if err != nil {
		fmt.Fprintln(os.Stderr, "FAIL:", err)
		os.Exit(1)
	}
	for _, scan := range page.Data {
		raw, adj := "-", "-"
		if scan.RawScore != nil {
			raw = fmt.Sprintf("%d", *scan.RawScore)
		}
		if scan.AdjustedScore != nil {
			adj = fmt.Sprintf("%d", *scan.AdjustedScore)
		}
		fmt.Printf("  %s  status=%s  raw=%s  adj=%s\n", scan.ID, scan.Status, raw, adj)
	}
	if page.HasMore {
		fmt.Println("(more pages available)")
	}
	fmt.Println("OK")
}
