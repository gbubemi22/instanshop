package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"instashop/internal/utils"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user_id and role in the context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		fmt.Printf("Extracted user_id: %s, role: %s from the token\n", claims.UserID, claims.Role)

		// Continue to the next handler
		c.Next()
	}
}

// func VerifyToken() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Retrieve the access token secret from environment variables
// 		accessTokenSecret := os.Getenv("JWT_SECRET")
// 		serviceName := os.Getenv("SERVICE_NAME")

// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"success":        false,
// 				"message":        "Authorization header is missing or invalid",
// 				"httpStatusCode": http.StatusUnauthorized,
// 				"error":          "VALIDATION_ERROR",
// 				"service":        serviceName,
// 			})
// 			c.Abort()
// 			return
// 		}

// 		token := authHeader[7:]

// 		// Verify token
// 		claims := jwt.MapClaims{}
// 		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(accessTokenSecret), nil
// 		})

// 		if err != nil || !jwtToken.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"success":        false,
// 				"message":        "Invalid token",
// 				"httpStatusCode": http.StatusUnauthorized,
// 				"error":          "VALIDATION_ERROR",
// 				"service":        serviceName,
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Extract user_id and validate
// 		userID, ok := claims["user_id"].(string)
// 		if !ok {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"success":        false,
// 				"message":        "Invalid user ID format",
// 				"httpStatusCode": http.StatusUnauthorized,
// 				"error":          "VALIDATION_ERROR",
// 				"service":        serviceName,
// 			})
// 			c.Abort()
// 			return
// 		}
// 		c.Set("user_id", userID)

// 		// Extract role and validate
// 		role, ok := claims["role"].(string)
// 		if !ok {
// 			role = "guest" // Default role
// 		}
// 		c.Set("role", role)

// 		// Debug logs
// 		fmt.Printf("Extracted user_id: %s, role: %s\n", userID, role)

// 		// Proceed to the next middleware or handler
// 		c.Next()
// 	}
// }
