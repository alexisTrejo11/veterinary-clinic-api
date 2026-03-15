// Package middleware contains middleware for authentication and user context management.
// Auth is stateless: credentials in Gin context come only from JWT claims (no repository calls).
package middleware

import (
	"context"
	"strconv"
	"strings"

	"clinic-vet-api/internal/core/auth"
	"clinic-vet-api/internal/shared"
	"clinic-vet-api/internal/shared/errors"
	sharedhttp "clinic-vet-api/internal/shared/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const bearerPrefix = "Bearer "

// UserContext holds claims put on context from the JWT (stateless, no DB).
type UserContext struct {
	UserID uint
	Email  string
	Role   string
}

func (uc *UserContext) UserIDValueObject() shared.UserID {
	return shared.NewUserID(uc.UserID)
}

// AuthMiddleware parses JWT and sets user claims on context. No repository.
type AuthMiddleware struct {
	jwtKey []byte
}

// NewAuthMiddleware creates middleware that validates JWT and sets claims on context.
// jwtSecret must match the secret used by the token service to sign access tokens.
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtKey: []byte(jwtSecret)}
}

// Authenticate requires a valid Bearer JWT. Sets userID, userEmail, userRole and user (UserContext) from token claims.
func (am *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		if raw == "" {
			sharedhttp.Unauthorized(c, errors.UnauthorizedError(context.Background(), "Authenticate", "auth token required"))
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(raw, bearerPrefix)
		if tokenStr == raw {
			sharedhttp.Unauthorized(c, errors.UnauthorizedError(context.Background(), "Authenticate", "invalid authorization header"))
			c.Abort()
			return
		}

		claims, err := am.parseAccessToken(tokenStr)
		if err != nil {
			sharedhttp.Unauthorized(c, errors.UnauthorizedError(context.Background(), "Authenticate", err.Error()))
			c.Abort()
			return
		}

		userID, email, role := claims[auth.ClaimUserID], claims[auth.ClaimEmail], claims[auth.ClaimRole]
		userIDStr, _ := userID.(string)
		emailStr, _ := email.(string)
		roleStr, _ := role.(string)

		if userIDStr == "" {
			sharedhttp.Unauthorized(c, errors.UnauthorizedError(context.Background(), "Authenticate", "invalid token claims"))
			c.Abort()
			return
		}

		idUint, err := strconv.ParseUint(userIDStr, 10, 0)
		if err != nil {
			sharedhttp.Unauthorized(c, errors.UnauthorizedError(context.Background(), "Authenticate", "invalid user id in token"))
			c.Abort()
			return
		}

		uc := &UserContext{
			UserID: uint(idUint),
			Email:  emailStr,
			Role:   roleStr,
		}

		c.Set("jwtToken", tokenStr)
		c.Set("user", uc)
		c.Set("userID", uint(idUint))
		c.Set("userEmail", emailStr)
		c.Set("userRole", roleStr)
		c.Next()
	}
}

// OptionalAuth parses JWT if present and sets claims on context; does not abort if missing.
func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		if raw == "" {
			c.Next()
			return
		}

		tokenStr := strings.TrimPrefix(raw, bearerPrefix)
		if tokenStr == raw {
			c.Next()
			return
		}

		claims, err := am.parseAccessToken(tokenStr)
		if err != nil {
			c.Next()
			return
		}

		userID, email, role := claims[auth.ClaimUserID], claims[auth.ClaimEmail], claims[auth.ClaimRole]
		userIDStr, _ := userID.(string)
		emailStr, _ := email.(string)
		roleStr, _ := role.(string)

		if userIDStr == "" {
			c.Next()
			return
		}

		idUint, err := strconv.ParseUint(userIDStr, 10, 0)
		if err != nil {
			c.Next()
			return
		}

		uc := &UserContext{
			UserID: uint(idUint),
			Email:  emailStr,
			Role:   roleStr,
		}

		c.Set("jwtToken", tokenStr)
		c.Set("user", uc)
		c.Set("userID", uint(idUint))
		c.Set("userEmail", emailStr)
		c.Set("userRole", roleStr)
		c.Next()
	}
}

// parseAccessToken parses and validates the JWT, returns map claims or error.
func (am *AuthMiddleware) parseAccessToken(tokenStr string) (map[string]interface{}, error) {
	tok, err := jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		return am.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tok.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	out := make(map[string]interface{})
	for k, v := range claims {
		out[k] = v
	}
	return out, nil
}

// GetUserFromContext returns the user context set by Authenticate/OptionalAuth (from JWT claims).
func GetUserFromContext(c *gin.Context) (*UserContext, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	uc, ok := user.(*UserContext)
	return uc, ok
}

// GetUserIDFromContext returns the user ID from context (set from JWT claims).
func GetUserIDFromContext(c *gin.Context) (shared.UserID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return shared.UserID{}, errors.UnauthorizedError(context.Background(), "GetUserIDFromContext", "user ID not found in context")
	}
	idUint, ok := userID.(uint)
	if !ok {
		return shared.UserID{}, errors.UnauthorizedError(context.Background(), "GetUserIDFromContext", "invalid user ID type in context")
	}
	return shared.NewUserID(idUint), nil
}

// GetUserEmailFromContext returns the email from context (from JWT claims).
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("userEmail")
	if !exists {
		return "", false
	}
	s, ok := email.(string)
	return s, ok
}

// GetUserRoleFromContext returns the role from context (from JWT claims).
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get("userRole")
	if !exists {
		return "", false
	}
	s, ok := role.(string)
	return s, ok
}

// HasRole returns true if the context user has the given role.
func HasRole(c *gin.Context, role string) bool {
	r, ok := GetUserRoleFromContext(c)
	return ok && r == role
}

// RequireAnyRole aborts with 403 if the context user's role is not in the given list.
// Must be used after Authenticate().
func (am *AuthMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := GetUserRoleFromContext(c)
		if !exists {
			sharedhttp.Unauthorized(c, errors.UnauthorizedError(context.Background(), "RequireAnyRole", "user role not found in context"))
			c.Abort()
			return
		}

		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		sharedhttp.Forbidden(c, errors.ForbiddenError(context.Background(), "RequireAnyRole", "access", "insufficient permissions: role "+userRole))
		c.Abort()
	}
}
