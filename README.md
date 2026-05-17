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
    // Create a client with static token
    client, err := directus.NewClient(
        "https://your-directus-instance.com",
        directus.WithStaticToken("your-access-token"),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Read items (without filtering)
    ctx := context.Background()
    users, err := directus.NewReadItems[User]("users").
        SetLimit(10).
        SetOffset(0).
        SendBy(client)
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

ctx := directus.WithAccessTokenContext(context.Background(), "dynamic-token")
// Use ctx in requests
```

### Login with Credentials
```go
resp, err := client.AuthLogin(ctx, directus.AuthLoginParams{
    Email:    "admin@example.com",
    Password: "password",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Access token: %s\n", resp.AccessToken)
```

### Refresh Token
```go
resp, err := client.AuthRefresh(ctx, directus.AuthRefreshParams{
    RefreshToken: "refresh-token",
})
```

### OTP Verification
```go
resp, err := client.AuthOTPVerify(ctx, directus.OTPVerifyParams{
    OTP: "otp-code",
})
```

### Password Reset Request
```go
err := client.AuthResetPasswordRequest(ctx, directus.PasswordResetRequestParams{
    Email: "user@example.com",
})
```

### Password Reset
```go
err := client.AuthPasswordReset(ctx, directus.PasswordResetParams{
    Token:    "reset-token",
    Password: "new-password",
})
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
created, err := directus.NewCreateItem[User]("users", newUser).
    SendBy(client)
```

### Read Item
```go
user, err := directus.NewReadItem[User]("users", "user-123").
    SendBy(client)
```

### Read Items with Filtering
```go
// Create a filter using ByField and Equal
filter := directus.ByField{
    Name: "status",
    Filter: directus.Equal[string]{Value: "active"},
}
users, err := directus.NewReadItems[User]("users").
    SetFilter(filter).
    SetSort("-created_at").
    SetLimit(100).
    SetOffset(0).
    SendBy(client)
```

### Update Item
```go
updated, err := directus.NewUpdateItem[User]("users", "user-123").
    SetChanges(map[string]any{"name": "John Updated"}).
    SendBy(client)
```

### Update Items (Batch)
```go
err := directus.NewUpdateItems[User]("users").
    SetFilter(directus.ByField{
        Name: "status",
        Filter: directus.Equal[string]{Value: "inactive"},
    }).
    SetChanges(map[string]any{"archived": true}).
    SendBy(client)
```

### Delete Item
```go
err := directus.NewDeleteItem("users", "user-123").
    SendBy(client)
```

### Delete Items (Batch)
```go
err := directus.NewDeleteItems[string]("users").
    SetFilter(directus.ByField{
        Name: "status",
        Filter: directus.Equal[string]{Value: "deleted"},
    }).
    SendBy(client)
```

## Filtering

The SDK provides type-safe filter rules via structs. You can compose filters using the `ByField` wrapper.

### Basic comparisons
```go
// Equal
filter := directus.ByField{
    Name: "status",
    Filter: directus.Equal[string]{Value: "active"},
}

// Not equal
filter := directus.ByField{
    Name: "status",
    Filter: directus.NotEqual[string]{Value: "inactive"},
}

// Greater than
filter := directus.ByField{
    Name: "age",
    Filter: directus.GreaterThan[int]{Value: 18},
}

// Less than
filter := directus.ByField{
    Name: "price",
    Filter: directus.LessThan[float64]{Value: 100.0},
}
```

### String operations
```go
// Contains
filter := directus.ByField{
    Name: "title",
    Filter: directus.Contains{Value: "important"},
}

// Starts with
filter := directus.ByField{
    Name: "email",
    Filter: directus.StartsWith{Value: "admin@"},
}

// Ends with
filter := directus.ByField{
    Name: "filename",
    Filter: directus.EndsWith{Value: ".pdf"},
}
```

### Logical operators
```go
// AND
filter := directus.AND{
    Filters: []directus.FilterRule{
        directus.ByField{Name: "status", Filter: directus.Equal[string]{Value: "active"}},
        directus.ByField{Name: "age", Filter: directus.GreaterThan[int]{Value: 18}},
    },
}

// OR
filter := directus.OR{
    Filters: []directus.FilterRule{
        directus.ByField{Name: "role", Filter: directus.Equal[string]{Value: "admin"}},
        directus.ByField{Name: "role", Filter: directus.Equal[string]{Value: "moderator"}},
    },
}

// NOT
filter := directus.NOT{
    Filter: directus.ByField{Name: "status", Filter: directus.Equal[string]{Value: "deleted"}},
}
```

### Null checks
```go
filter := directus.ByField{
    Name: "deleted_at",
    Filter: directus.IsNull{},
}

filter := directus.ByField{
    Name: "updated_at",
    Filter: directus.IsNotNull{},
}
```

### In / NotIn
```go
filter := directus.ByField{
    Name: "category",
    Filter: directus.In[string]{Values: []string{"news", "blog", "tutorial"}},
}

filter := directus.ByField{
    Name: "status",
    Filter: directus.NotIn[string]{Values: []string{"draft", "archived"}},
}
```

### Between
```go
filter := directus.ByField{
    Name: "price",
    Filter: directus.Between[float64]{Low: 10.0, High: 100.0},
}
```

## Aggregation

Aggregation is performed using `directus.NewAggregate` and `SetAggregate` with an aggregate rule.

### Count
```go
result, err := directus.NewAggregate[any]("users").
    SetAggregate(directus.Count{}).
    SetFilter(directus.ByField{
        Name: "status",
        Filter: directus.Equal[string]{Value: "active"},
    }).
    SendBy(client)
```

### Sum
```go
result, err := directus.NewAggregate[any]("orders").
    SetAggregate(directus.Sum{Fields: []string{"total_amount"}}).
    SetFilter(directus.ByField{
        Name: "status",
        Filter: directus.Equal[string]{Value: "completed"},
    }).
    SendBy(client)
```

### Multiple aggregations
```go
result, err := directus.NewAggregate[any]("products").
    SetAggregate(directus.Many{
        Rules: []directus.AggregateRule{
            directus.Count{},
            directus.Sum{Fields: []string{"stock"}},
            directus.Avg{Fields: []string{"price"}},
            directus.Min{Fields: []string{"price"}},
            directus.Max{Fields: []string{"price"}},
        },
    }).
    SendBy(client)
```

Note: Use `directus.Many` to combine multiple aggregate rules in a single request.

## Server Information

### Ping
```go
ping, err := client.ServerPing(ctx)
```

### Health
```go
health, err := client.ServerHealth(ctx)
```

### Info
```go
info, err := client.ServerInfo(ctx)
```

## Singleton Operations

### Read Singleton
```go
settings, err := directus.NewReadSingleton[Settings]("settings").
    SendBy(client)
```

### Update Singleton
Update of singletons is currently not implemented as a separate method. Use `UpdateItem` with the singleton collection name.

## Error Handling

The SDK returns structured errors:

```go
items, err := directus.NewReadItems[User]("users").SendBy(client)
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
users, err := directus.NewReadItems[User]("users").
    SetDeep(map[string]directus.DeepQuery{
        "posts": {
            Filter: directus.ByField{
                Name: "status",
                Filter: directus.Equal[string]{Value: "published"},
            },
            Limit: 5,
        },
    }).
    SendBy(client)
```

### System Collections
```go
// Access Directus system collections
roles, err := directus.NewReadItems[Role]("directus_roles").
    SetIsSystem(true).
    SendBy(client)
```

## Configuration

### Client Options
```go
client, err := directus.NewClient(
    "https://api.example.com",
    directus.WithStaticToken("your-access-token"),
    directus.WithExtractTokenFromContext(true),
)
```

### Custom HTTP Client
The SDK uses `github.com/go-resty/resty/v2` internally. You can customize the underlying HTTP client by accessing the resty client directly:

```go
// Create client first
client, err := directus.NewClient("https://api.example.com")
if err != nil {
    log.Fatal(err)
}

// Access the underlying resty client
restyClient := client.resty // Note: field is not exported, you need to use reflection or provide a getter.
// Currently there is no public method to access the resty client.
```

## Requirements

- Go 1.25.0 or higher
- Directus API v11+

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
