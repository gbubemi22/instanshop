package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"instashop/internal/utils"
)

func ErrorHandlerMiddleware(c *gin.Context) {
	c.Next() // Continue middleware pipeline and handlers

	serviceName := os.Getenv("SERVICE_NAME")

	// Check if there's an error to handle
	err := c.Errors.Last()
	if err != nil {
		var statusCode int
		var response gin.H

		switch customErr := err.Err.(type) {
		case *utils.CustomError:
			// Handle custom errors
			statusCode = customErr.HTTPStatusCode
			response = gin.H{
				"success":        customErr.Success,
				"message":        customErr.Message,
				"httpStatusCode": customErr.HTTPStatusCode,
				"error":          http.StatusText(customErr.HTTPStatusCode),
				"service":        customErr.Service,
			}
		case *mysql.MySQLError:
			// Handle MySQL-specific errors
			switch customErr.Number {
			case 1062:
				statusCode = http.StatusConflict
				response = gin.H{
					"success":        false,
					"message":        "Duplicate value entered",
					"httpStatusCode": statusCode,
					"error":          "DUPLICATE_ENTRY",
					"service":        serviceName,
				}
			case 1451, 1452:
				statusCode = http.StatusBadRequest
				response = gin.H{
					"success":        false,
					"message":        "Foreign key constraint error",
					"httpStatusCode": statusCode,
					"error":          "FOREIGN_KEY_CONSTRAINT",
					"service":        serviceName,
				}
			default:
				statusCode = http.StatusInternalServerError
				response = gin.H{
					"success":        false,
					"message":        strings.TrimSpace(customErr.Message),
					"httpStatusCode": statusCode,
					"error":          "DATABASE_ERROR",
					"service":        serviceName,
				}
			}
		default:
			// Default to internal server error
			statusCode = http.StatusInternalServerError
			response = gin.H{
				"success":        false,
				"message":        "Internal Server Error",
				"httpStatusCode": statusCode,
				"error":          "INTERNAL_SERVER_ERROR",
				"service":        serviceName,
			}
		}

		// Send JSON response
		c.JSON(statusCode, response)
	}
}
