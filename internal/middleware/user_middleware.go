package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project/internal/models"
	"github.com/keykibatyr/qazaq-society-project/internal/utils"
)

type UserMiddleware struct {
	SessionService *models.SessionService
	CookieName     string
	SignInPage     string
}

func (um *UserMiddleware) SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := utils.ReadCookie(c.Request, um.CookieName)
		if err != nil {
			c.Next()
			return
		}

		user, err := um.SessionService.User(tokenCookie)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func (um *UserMiddleware) RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exist := c.Get("user")
		if !exist {
			c.Redirect(http.StatusFound, um.SignInPage)
			c.Abort()
			return
		}

		c.Next()
	}
}

func CurrentUser(c *gin.Context) *models.User {
	val, exist := c.Get("user")
	if !exist {
		return nil
	}

	return val.(*models.User)
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user.Role != "admin" {
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		c.Next()

	}
}