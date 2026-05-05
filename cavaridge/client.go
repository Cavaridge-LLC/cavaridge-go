// Package cavaridge is the official Go SDK for the Cavaridge™ public API.
//
// The transport stub is auto-generated from openapi.yaml via oapi-codegen
// into ./generated/. This package exposes the ergonomic wrappers
// documented in the README.
package cavaridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	SDKVersion        = "0.1.0"
	DefaultAPIVersion = "2026-04-27"
	DefaultBaseURL    = "https://api.cavaridge.app"
)

type Client struct {
	apiKey      string
	bearerToken string
	baseURL     string
	apiVersion  string
	maxRetries  int
	http        *http.Client
}

type Option func(*Client)

func WithAPIKey(key string) Option       { return func(c *Client) { c.apiKey = key } }
func WithBearerToken(t string) Option    { return func(c *Client) { c.bearerToken = t } }
func WithBaseURL(u string) Option        { return func(c *Client) { c.baseURL = u } }
func WithAPIVersion(v string) Option     { return func(c *Client) { c.apiVersion = v } }
func WithMaxRetries(n int) Option        { return func(c *Client) { c.maxRetries = n } }
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }

func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		baseURL:    DefaultBaseURL,
		apiVersion: DefaultAPIVersion,
		maxRetries: 3,
		http:       &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.apiKey == "" && c.bearerToken == "" {
		return nil, errors.New("cavaridge: apiKey or bearerToken required")
	}
	return c, nil
}

// Error is the RFC 7807 + Cavaridge-extended error envelope.
type Error struct {
	Type             string `json:"type"`
	Title            string `json:"title"`
	Status           int    `json:"status"`
	Detail           string `json:"detail,omitempty"`
	RequestID        string `json:"request_id,omitempty"`
	TenantID         string `json:"tenant_id,omitempty"`
	DocumentationURL string `json:"documentation_url,omitempty"`
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("cavaridge: %d %s — %s (%s)", e.Status, e.Title, e.Detail, e.Type)
	}
	return fmt.Sprintf("cavaridge: %d %s (%s)", e.Status, e.Title, e.Type)
}

// RequestOptions allows per-call overrides like Idempotency-Key.
type RequestOptions struct {
	IdempotencyKey string
}

// Request is the low-level transport call. Most callers use the typed
// resource wrappers; this is exposed for forward compatibility.
func (c *Client) Request(ctx context.Context, method, path string, query url.Values, body any, opts *RequestOptions, out any) error {
	full := c.baseURL + path
	if len(query) > 0 {
		full += "?" + query.Encode()
	}
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, full, bytes.NewReader(bodyBytes))
		if err != nil {
			return err
		}
		req.Header.Set("Cavaridge-Version", c.apiVersion)
		req.Header.Set("Content-Type", "application/json")
		if c.apiKey != "" {
			req.Header.Set("Authorization", "Cavaridge-Key "+c.apiKey)
		} else {
			req.Header.Set("Authorization", "Bearer "+c.bearerToken)
		}
		if opts != nil && opts.IdempotencyKey != "" {
			req.Header.Set("Idempotency-Key", opts.IdempotencyKey)
		}

		resp, err := c.http.Do(req)
		if err != nil {
			if attempt < c.maxRetries {
				time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * 250 * time.Millisecond)
				continue
			}
			return err
		}
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if out != nil && len(respBody) > 0 {
				return json.Unmarshal(respBody, out)
			}
			return nil
		}

		// Retry on 429 + 5xx.
		if resp.StatusCode == 429 || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
			if attempt < c.maxRetries {
				delay := 250 * time.Millisecond * time.Duration(math.Pow(2, float64(attempt)))
				if ra := resp.Header.Get("Retry-After"); ra != "" {
					if seconds, err := strconv.Atoi(ra); err == nil && seconds > 0 {
						delay = time.Duration(seconds) * time.Second
					}
				}
				time.Sleep(delay)
				continue
			}
		}

		var apiErr Error
		_ = json.Unmarshal(respBody, &apiErr)
		if apiErr.Type == "" {
			apiErr.Type = "https://docs.cavaridge.app/api/errors/internal_error"
		}
		if apiErr.Title == "" {
			apiErr.Title = http.StatusText(resp.StatusCode)
		}
		if apiErr.Status == 0 {
			apiErr.Status = resp.StatusCode
		}
		return &apiErr
	}
	return errors.New("cavaridge: max retries exhausted")
}

// ── Typed wrappers (subset; full surface in ./generated) ──────────────

type Aegis struct{ c *Client }
type Forge struct{ c *Client }
type Cavalier struct{ c *Client }
type Reclaim struct{ c *Client }
type Meta struct{ c *Client }

func (c *Client) Aegis() *Aegis       { return &Aegis{c} }
func (c *Client) Forge() *Forge       { return &Forge{c} }
func (c *Client) Cavalier() *Cavalier { return &Cavalier{c} }
func (c *Client) Reclaim() *Reclaim   { return &Reclaim{c} }
func (c *Client) Meta() *Meta         { return &Meta{c} }

type Page[T any] struct {
	Object     string `json:"object"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor,omitempty"`
	Data       []T    `json:"data"`
}

type Scan struct {
	ID            string  `json:"id"`
	TenantID      string  `json:"tenant_id"`
	ClientID      string  `json:"client_id,omitempty"`
	Status        string  `json:"status"`
	RawScore      *int    `json:"raw_score,omitempty"`
	AdjustedScore *int    `json:"adjusted_score,omitempty"`
	StartedAt     string  `json:"started_at"`
	CompletedAt   *string `json:"completed_at,omitempty"`
}

func (a *Aegis) ListScans(ctx context.Context, args struct {
	ClientID       string
	Limit          int
	StartingAfter  string
}) (*Page[Scan], error) {
	q := url.Values{}
	if args.ClientID != "" {
		q.Set("client_id", args.ClientID)
	}
	if args.Limit > 0 {
		q.Set("limit", strconv.Itoa(args.Limit))
	}
	if args.StartingAfter != "" {
		q.Set("starting_after", args.StartingAfter)
	}
	var out Page[Scan]
	if err := a.c.Request(ctx, "GET", "/v2026-04-27/aegis/scans", q, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *Aegis) CreateScan(ctx context.Context, body struct {
	ClientID     string   `json:"client_id"`
	ConnectorSet []string `json:"connector_set,omitempty"`
}, opts *RequestOptions) (*Scan, error) {
	var out Scan
	if err := a.c.Request(ctx, "POST", "/v2026-04-27/aegis/scans", nil, body, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *Aegis) GetScan(ctx context.Context, id string) (*Scan, error) {
	var out Scan
	if err := a.c.Request(ctx, "GET", "/v2026-04-27/aegis/scans/"+url.PathEscape(id), nil, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
