# cavaridge-go

Official Go SDK for the Cavaridge™ public API.

## Install

```bash
go get github.com/Cavaridge-LLC/cavaridge-go/cavaridge
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/Cavaridge-LLC/cavaridge-go/cavaridge"
)

func main() {
    client, _ := cavaridge.NewClient(cavaridge.WithAPIKey(os.Getenv("CAVARIDGE_API_KEY")))
    page, err := client.Aegis().ListScans(context.Background(), struct {
        ClientID, StartingAfter string; Limit int
    }{Limit: 5})
    if err != nil { panic(err) }
    for _, s := range page.Data {
        fmt.Println(s.ID, s.Status)
    }
}
```

## Auth

- **API key**: `cavaridge.WithAPIKey("sk_live_...")` (sandbox: `sk_test_...`).
- **OAuth bearer**: `cavaridge.WithBearerToken("...")`.

## Idempotency

```go
client.Aegis().CreateScan(ctx, body, &cavaridge.RequestOptions{IdempotencyKey: "scn_create_2026-05-04_001"})
```

## Errors

Failures return a `*cavaridge.Error` (implements `error`) with `.Type`, `.Status`, `.Detail`, `.RequestID`, `.DocumentationURL`.

```go
if err := client.Reclaim().ExecutePolicy(ctx, "pol_abc"); err != nil {
    var apiErr *cavaridge.Error
    if errors.As(err, &apiErr) && apiErr.Status == 403 {
        fmt.Println("need reclaim.policies.manage scope")
    }
}
```

## Retries

Default: 3 retries on `429` and 5xx with exponential backoff (250ms → 500ms → 1000ms). Override with `cavaridge.WithMaxRetries(n)`.

## Sandbox

```go
cavaridge.NewClient(
    cavaridge.WithAPIKey(os.Getenv("CAVARIDGE_TEST_KEY")),
    cavaridge.WithBaseURL("https://api-sandbox.cavaridge.app"),
)
```

## Versioning

SDK semver maps to spec dated versions:

- **major** → API dated version
- **minor** → additive features within a dated version
- **patch** → bug fixes only

## Docs

- [API overview](https://docs.cavaridge.app/api/overview/)
- [Patterns](https://github.com/Cavaridge-LLC/cavaridge/blob/main/packages/public-api-spec/PATTERNS.md)
- [Go SDK guide](https://docs.cavaridge.app/sdks/go/)
