package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project/internal/models"
	"github.com/keykibatyr/qazaq-society-project/internal/utils"
)


const (
	CookieSession = "session"
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
		val, exist := c.Get("user")
		if !exist{
			c.Redirect(http.StatusFound, um.SignInPage)
			c.Abort()
			return
		}
		user := val.(*models.User)
		expired, err := um.SessionService.IsExpired(user.ID)
		if  err != nil || expired{
			token, _ := utils.ReadCookie(c.Request, CookieSession)
			_ = um.SessionService.Delete(token)
			utils.DeleteCookie(c.Writer, CookieSession)
			
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

	user, ok := val.(*models.User) 
	if !ok {
		return nil
	}

	return user
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