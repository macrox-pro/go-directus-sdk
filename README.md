# Go Directus SDK

A pure Go library for interacting with Directus CMS API. Provides type-safe client methods for CRUD operations, authentication, filtering, and aggregation.

## Features

- **Type-safe**: Uses Go generics for compile-time type checking
- **Complete API coverage**: Supports all major Directus REST API endpoints
- **Builder pattern**: Fluent API for constructing complex queries
- **Context-aware**: Full support for Go contexts and cancellation
- **Authentication**: Static tokens, context-based tokens, OTP, password reset
- **Filtering**: Type-safe filter rules with logical operators
- **Aggregation**: Built-in support for Directus aggregation functions
- **Error handling**: Structured errors with HTTP status codes
- **No CGO dependencies**: Pure Go implementation

## Installation

```bash
go get github.com/macrox-pro/go-directus-sdk
```

Requires Go 1.25.0 or higher.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/macrox-pro/go-directus-sdk"
)

type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // Create a client
    client, err := directus.NewClient("https://your-directus-instance.com")
    if err != nil {
        log.Fatal(err)
    }
    
    // Set authentication token
    client = directus.WithStaticToken("your-access-token")
    
    // Read items with filtering
    ctx := context.Background()
    users, err := client.ReadItems[User]("users").
        SetFilter(directus.Equals("status", "active")).
        SetLimit(10).
        SetOffset(0).
        Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, user := range users {
        fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
    }
}
```

## Authentication

### Static Token
```go
client, err := directus.NewClient(
    "https://api.example.com",
    directus.WithStaticToken("your-access-token"),
)
```

### Context-Based Token
```go
client, err := directus.NewClient(
    "https://api.example.com",
    directus.WithExtractTokenFromContext(true),
)

ctx := context.WithValue(context.Background(), 
    directus.AccessTokenContextKey, "dynamic-token")
// Use ctx in requests
```

### Login with Credentials
```go
resp, err := client.AuthLogin("admin@example.com", "password").Do(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Access token: %s\n", resp.AccessToken)
```

### Refresh Token
```go
resp, err := client.AuthRefresh("refresh-token").Do(ctx)
```

### OTP Verification
```go
err := client.AuthOTPVerify("otp-code").Do(ctx)
```

### Password Reset
```go
err := client.AuthResetPassword("token", "new-password").Do(ctx)
```

### Auth Providers
```go
providers, err := client.AuthProviders(ctx)
if err != nil {
    log.Fatal(err)
}
for _, provider := range providers {
    fmt.Printf("Provider: %s (%s)\n", provider.Name, provider.Driver)
}
```

## CRUD Operations

### Create Item
```go
newUser := User{Name: "John Doe", Email: "john@example.com"}
created, err := client.CreateItem[User]("users").
    SetData(newUser).
    Do(ctx)
```

### Read Item
```go
user, err := client.ReadItem[User]("users", "user-123").Do(ctx)
```

### Read Items with Filtering
```go
users, err := client.ReadItems[User]("users").
    SetFilter(directus.And(
        directus.Equals("status", "active"),
        directus.GreaterThan("age", 18),
    )).
    SetSort("-created_at").
    SetLimit(100).
    SetOffset(0).
    SetSearch("john").
    Do(ctx)
```

### Update Item
```go
updated, err := client.UpdateItem[User]("users", "user-123").
    SetData(map[string]any{"name": "John Updated"}).
    Do(ctx)
```

### Update Items (Batch)
```go
err := client.UpdateItems("users").
    SetFilter(directus.Equals("status", "inactive")).
    SetData(map[string]any{"archived": true}).
    Do(ctx)
```

### Delete Item
```go
err := client.DeleteItem("users", "user-123").Do(ctx)
```

### Delete Items (Batch)
```go
err := client.DeleteItems("users").
    SetFilter(directus.Equals("status", "deleted")).
    Do(ctx)
```

## Filtering

The SDK provides type-safe filter rules:

```go
// Basic comparisons
directus.Equals("status", "active")
directus.NotEquals("status", "inactive")
directus.GreaterThan("age", 18)
directus.GreaterThanOrEqual("score", 80)
directus.LessThan("price", 100)
directus.LessThanOrEqual("quantity", 10)

// String operations
directus.Contains("title", "important")
directus.StartsWith("email", "admin@")
directus.EndsWith("filename", ".pdf")

// Logical operators
directus.And(
    directus.Equals("status", "active"),
    directus.GreaterThan("age", 18),
)

directus.Or(
    directus.Equals("role", "admin"),
    directus.Equals("role", "moderator"),
)

directus.Not(directus.Equals("status", "deleted"))

// Null checks
directus.IsNull("deleted_at")
directus.IsNotNull("updated_at")

// In/NotIn
directus.In("category", []string{"news", "blog", "tutorial"})
directus.NotIn("status", []string{"draft", "archived"})

// Between
directus.Between("price", 10, 100)
```

## Aggregation

```go
// Count items
count, err := client.Aggregate("users").
    SetFilter(directus.Equals("status", "active")).
    Count().
    Do(ctx)

// Sum, Average, Min, Max
result, err := client.Aggregate("orders").
    SetFilter(directus.Equals("status", "completed")).
    Sum("total_amount").
    Do(ctx)

// Multiple aggregations
result, err := client.Aggregate("products").
    Count().
    Sum("stock").
    Average("price").
    Min("price").
    Max("price").
    Do(ctx)
```

## Server Information

```go
// Ping
err := client.ServerPing().Do(ctx)

// Health
health, err := client.ServerHealth().Do(ctx)

// Info
info, err := client.ServerInfo().Do(ctx)
```

## Singleton Operations

```go
type Settings struct {
    SiteName string `json:"site_name"`
    Theme    string `json:"theme"`
}

// Read singleton
settings, err := client.ReadSingleton[Settings]("settings").Do(ctx)

// Update singleton
updated, err := client.UpdateSingleton[Settings]("settings").
    SetData(map[string]any{"theme": "dark"}).
    Do(ctx)
```

## Error Handling

The SDK returns structured errors:

```go
items, err := client.ReadItems[User]("users").Do(ctx)
if err != nil {
    if errs, ok := err.(directus.Errors); ok {
        for _, e := range errs {
            fmt.Printf("Error %d: %s\n", e.Status(), e.Message)
        }
    } else {
        log.Fatal(err)
    }
}
```

## Advanced Usage

### Deep Query
```go
users, err := client.ReadItems[User]("users").
    SetDeep(map[string]directus.DeepQuery{
        "posts": {
            Filter: directus.Equals("status", "published"),
            Limit:  5,
        },
    }).
    Do(ctx)
```

### Fields Selection
```go
users, err := client.ReadItems[User]("users").
    SetFields([]string{"id", "name", "email"}).
    Do(ctx)
```

### System Collections
```go
// Access Directus system collections
roles, err := client.ReadItems[Role]("directus_roles").
    SetSystem(true).
    Do(ctx)
```

## Configuration

### Client Options
```go
client, err := directus.NewClient(
    "https://api.example.com",
    directus.WithStaticToken("token"),
    directus.WithExtractTokenFromContext(true),
)
```

### Custom HTTP Client
The SDK uses `github.com/go-resty/resty/v2` internally. You can customize the underlying HTTP client:

```go
// Create client first
client, err := directus.NewClient("https://api.example.com")
if err != nil {
    log.Fatal(err)
}

// Access the underlying resty client
restyClient := client.RestyClient()
restyClient.SetTimeout(30 * time.Second)
restyClient.SetRetryCount(3)
```

## Project Structure

```
.
├── client.go              # Core client and constructor
├── client_options.go      # Client configuration options
├── auth_*.go             # Authentication methods
├── *_item.go             # CRUD operations
├── aggregate*.go         # Aggregation functionality
├── filter_rules.go       # Filter rule types
├── errors.go             # Error definitions
├── utils.go              # Internal utilities
├── helpers/              # Internal helper packages
│   ├── fields_extractor.go
│   ├── join_parts_url.go
│   └── url_param_json.go
└── README.md             # This file
```

## Requirements

- Go 1.25.0 or higher
- Directus API v9+

## Dependencies

- `github.com/go-resty/resty/v2` - HTTP client
- `github.com/google/go-querystring` - URL query parameter encoding

## Testing

```bash
# Run tests
go test ./...

# Check for linting issues
go vet ./...

# Format code
go fmt ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Ensure `go build ./...` and `go test ./...` pass
6. Submit a pull request

Please follow the Go code style and ensure all exported types/functions are documented.

## License

MIT License - see LICENSE file for details.

## References

- [Directus API Documentation](https://docs.directus.io/reference/introduction/)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [Resty Documentation](https://github.com/go-resty/resty)
