# Authentication Methods in Distributed Architecture

This project implements two authentication strategies in a distributed architecture using Go:
1. **Token-based Authentication (JWT)**
2. **Username/Password Authentication**

The components include token generation/validation, credential verification, handlers, and middleware for securing endpoints.

---

## 1. Token-based Authentication

### Overview
Token-based authentication utilizes **JSON Web Tokens (JWT)** to ensure stateless and secure communication between the client and server. A JWT is issued upon successful login and included in the `Authorization` header for subsequent requests.

### Key Features
- **GenerateToken**: Issues a signed JWT with claims including `Issuer`, `Subject`, and `Expiration`.
- **ValidateToken**: Validates incoming JWTs to ensure they are not expired or tampered with.

### File: `pkg/auth/token.go`
```go
// GenerateToken generates a JWT token for a given user.
func GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "iss": "auth_service",
        "sub": userID,
        "exp": time.Now().Add(time.Hour * 1).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// ValidateToken parses and validates a JWT token.
func ValidateToken(tokenStr string) (*jwt.Token, error) {
    return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })
}
```

### Token Workflow
1. User logs in with valid credentials.
2. A JWT is generated and sent as a response.
3. The client includes the JWT in the `Authorization` header of subsequent requests.
4. Middleware validates the JWT before granting access to secured routes.

---

## 2. Username/Password Authentication

### Overview
This strategy authenticates users by validating provided credentials (username/password). It serves as the initial step before issuing tokens.

### Key Features
- **Authenticate**: Verifies provided credentials against the system.
- Hardcoded credentials are placeholders for integration with databases.

### File: `pkg/auth/username_password.go`
```go
// Authenticate verifies username and password.
func Authenticate(username, password string) error {
    if username == "admin" && password == "password123" {
        return nil // Success
    }
    return errors.New("invalid credentials")
}
```

---

## 3. Handlers

### Token Handler: `pkg/handlers/token_handler.go`
- Issues a JWT upon successful authentication.
```go
func TokenHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.FormValue("userID")
    token, err := auth.GenerateToken(userID)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
```

### Username/Password Handler: `pkg/handlers/username_handler.go`
- Authenticates credentials and returns a JWT if successful.
```go
func UsernameHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")

    if err := auth.Authenticate(username, password); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    
    token, _ := auth.GenerateToken(username)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
```

---

## 4. Authentication Middleware

### Overview
Middleware validates JWTs on protected routes and injects user data into the request context.

### File: `pkg/middleware/auth_middleware.go`
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenStr := r.Header.Get("Authorization")
        
        token, err := auth.ValidateToken(tokenStr)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        ctx := context.WithValue(r.Context(), "user", token.Claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Middleware Workflow
1. Extract JWT from the `Authorization` header.
2. Validate the JWT using `ValidateToken`.
3. If valid, attach user claims to the request context.
4. Reject requests with invalid or expired tokens.

---

## 5. Summary of Authentication Flow
1. **Login**:
    - User provides `username` and `password`.
    - Upon success, a JWT is issued.
2. **Token Usage**:
    - The client includes the JWT in the `Authorization` header for protected endpoints.
3. **Middleware Validation**:
    - Token is validated before granting access.

---

## Trade-offs and Challenges

### Token-based Authentication (JWT)
| Pros                           | Cons                                |
|--------------------------------|-------------------------------------|
| Stateless and scalable         | Token storage on the client side    |
| Reduces server load            | Tokens must be carefully secured    |
| Faster authentication          | Token expiration requires renewal   |
| Suitable for distributed systems | JWT size adds overhead to requests |

### Username/Password Authentication
| Pros                           | Cons                                |
|--------------------------------|-------------------------------------|
| Simple and easy to implement   | Requires server-side state management |
| No reliance on external libraries | Harder to scale in distributed systems |
| Familiar approach for users    | Passwords must be securely stored  |

---

## Design Decisions
- **JWT**: Chosen for its stateless and scalable nature in distributed systems.
- **Username/Password**: Used as a preliminary step to authenticate users before issuing tokens.
- **Middleware**: Centralized validation of tokens for securing endpoints.

---

## Component Structure
```
├── pkg/
│   ├── auth/
│   │   ├── token.go                # Token-based authentication logic
│   │   └── username_password.go    # Username/password authentication
│   ├── handlers/
│   │   ├── token_handler.go        # Token issuance handler
│   │   └── username_handler.go     # Username/password handler
│   └── middleware/
│       └── auth_middleware.go      # Middleware for token validation
└── main.go                         # Entry point for the application
```

---

## Future Improvements
- Integrate a database for storing user credentials.
- Add refresh token support for extending JWT lifespans.
- Implement stronger password hashing (e.g., bcrypt).
