package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/helpers"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

// RoleMiddleware is a middleware function to check the user's role
func RoleBasedMiddleware() web.FilterFunc {
	return func(ctx *beecontext.Context) {
		// Skip middleware for /auth routes
		if strings.HasPrefix(ctx.Request.URL.Path, "/auth") {
			return
		}

		// Get the user's token from the Authorization header
		authHeader := ctx.Input.Header("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.ErrorResponse(ctx.ResponseWriter, 401, "Unauthorized", errors.New("invalid token format"))
			return
		}

		// Extract the token part (remove the "Bearer " prefix)
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := helpers.ParseJWT(token)
		if err != nil {
			helpers.ErrorResponse(ctx.ResponseWriter, 401, "Unauthorized or expired token", err)
			return
		}

		//  extract the role from claiams
		userRole, err := helpers.GetRoleFromMapClaims(claims)
		if err != nil {
			helpers.ErrorResponse(ctx.ResponseWriter, 500, "Internal server error", err)
			return
		}

		// Get the URL and HTTP method of the current request
		url := ctx.Request.URL.Path
		method := ctx.Request.Method

		// Check if access is restricted for the current role based on the endpoint and method
		if isRestrictedAccess(url, method, userRole) {
			helpers.ErrorResponse(ctx.ResponseWriter, 403, "Forbidden", errors.New("access denied"))
			return
		}

		//  extract the role from claiams
		userId, err := helpers.GetUseridFromMapClaims(claims)
		if err != nil {
			helpers.ErrorResponse(ctx.ResponseWriter, 500, "Internal server error", err)
			return
		}

		// Continue to the next handler
		ctx.Input.SetData(constants.CtxAuthenticatedUserId, userId)
		ctx.Input.SetData(constants.CtxAuthenticatedUserRole, userRole)
	}
}

// isRestrictedAccess checks if the access to the endpoint is restricted based on the role and method
func isRestrictedAccess(url, method, role string) bool {
	// Check if the URL contains "/departments" or "/schedules" and if the method is POST, PUT, or DELETE
	if strings.Contains(url, "/departments") || strings.Contains(url, "/schedules") {
		if method == "POST" || method == "PUT" || method == "DELETE" {
			// Only allow admins to access POST, PUT, and DELETE methods
			return role != constants.RoleAdmin
		}
	}

	// Check if the URL contains "/presences"
	if strings.Contains(url, "/presences") {
		if method == "POST" {
			fmt.Println("role", role)
			fmt.Println("constants.RoleEmployee", constants.RoleEmployee)
			// Only allow users to access the POST method for creating a presence
			return role != constants.RoleEmployee
		} else if method == "PUT" || method == "DELETE" {
			// Only allow admins to access PUT and DELETE methods for updating or deleting presence
			return role != constants.RoleAdmin
		}
	}

	// For GET methods, all users (admin or user) are allowed
	return false
}
