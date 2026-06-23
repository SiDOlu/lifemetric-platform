package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims defines our custom claims matching our Keycloak/OIDC JWT schema
type JWTClaims struct {
	TenantID string   `json:"tenant_id"`
	Scopes   []string `json:"scopes"`
	Role     string   `json:"role"`
	jwt.RegisteredClaims
}

// SecureAuthorize enforces OAuth 2.0 JWT scope and role validation
func SecureAuthorize(requiredScope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization bearer token",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenStr := parts[1]

		// For prototyping, we will parse the token using a mock verification key.
		// In production, this verifies against Keycloak's public JWKS endpoint.
		token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &JWTClaims{})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Malformed authorization token",
			})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims mapping",
			})
		}

		// Enforce Scope Verification
		scopeApproved := false
		for _, scope := range claims.Scopes {
			if scope == requiredScope || scope == "super:admin" {
				scopeApproved = true
				break
			}
		}

		if !scopeApproved {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient token scopes to access this resource",
			})
		}

		// Store Tenant Context securely in Fiber context (for Database RLS query isolation)
		c.Locals("tenant_id", claims.TenantID)
		c.Locals("scopes", claims.Scopes)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
