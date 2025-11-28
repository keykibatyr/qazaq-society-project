package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project.git/internal/middleware"
	"github.com/keykibatyr/qazaq-society-project.git/internal/models"
	"github.com/keykibatyr/qazaq-society-project.git/internal/utils"
)

const (
	CookieSession = "session"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}

	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(c *gin.Context) {
	var data struct {
		Email string
	}

	data.Email = c.PostForm("email")

	u.Templates.New.ExecuteTemplate(c.Writer, c.Request, data)
}

func (u Users) Create(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	firstName := c.PostForm("first_name")
	secondName := c.PostForm("second_name")

	newUser, err := u.UserService.Create(email, password, firstName, secondName)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not create a user")
		return
	}
	// fmt.Fprintf(c.Writer, "the created user is: %+v", newUser)
	session, err := u.SessionService.Create(newUser.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not create a session token")
		return
	}

	utils.SetCookie(c.Writer, CookieSession, session.Token)
	c.Redirect(http.StatusFound, "/users/me")

}

func (u Users) SignIn(c *gin.Context) {
	var data struct {
		Email string
	}

	data.Email = c.PostForm("email")

	u.Templates.SignIn.ExecuteTemplate(c.Writer, c.Request, data)
}

func (u Users) ProcessSignIn(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not sign in")
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not create a session token")
		return
	}

	utils.SetCookie(c.Writer, CookieSession, session.Token)

	// fmt.Fprintf(c.Writer, "Signed In as %+v", user)
	c.Redirect(http.StatusFound, "/users/me")
}

func (u Users) CurrentUserController(c *gin.Context) {
	user := middleware.CurrentUser(c)
	fmt.Println(user)
	if user == nil {
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	c.String(http.StatusOK, "Current User is: %v", user)
	//TODO: Check for expired session token?
	//
	///
	///
	///
	///
}
