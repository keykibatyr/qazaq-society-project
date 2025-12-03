package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project/internal/middleware"
	"github.com/keykibatyr/qazaq-society-project/internal/models"
	"github.com/keykibatyr/qazaq-society-project/internal/utils"
)

const (
	CookieSession = "session"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
		Events Template
		Home Template
	}

	UserService    *models.UserService
	SessionService *models.SessionService
	EventService *models.EventService
}

func (u Users) Home(c *gin.Context) {
	event, err := u.EventService.GetLatest()
	if err != nil {
		c.String(500, "error")
		return
	}

	data := gin.H{
		"Title" : event.Title,
		"Time" : event.Time,
		"Month": event.Month,
		"Day": event.Day,
		"ImageURL": event.ImageURL,
	}

	Render(c, u.Templates.Home, data)
}

func (u Users) Events(c *gin.Context) {
	events, err := u.EventService.GetAll()
	if err != nil {
		c.String(500, "error")
		return
	}

	data := gin.H{
		"Events" : events,
	}

	Render(c, u.Templates.Events, data)
}

func (u Users) New(c *gin.Context) {
	data := gin.H{
        "Email": "",
	}
	Render(c, u.Templates.New, data)
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
	data := gin.H{
        "Email": "",
	}
	Render(c, u.Templates.SignIn, data)
}

func (u Users) ProcessSignIn(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not sign in")
		return
	}

	fmt.Print(user)

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not create a session token")
		return
	}

	utils.SetCookie(c.Writer, CookieSession, session.Token)
	// fmt.Fprintf(c.Writer, "Signed In as %+v", user)
	c.Redirect(http.StatusFound, "/users/me")
}

func (u Users) ProcessSignOut(c *gin.Context) {
	token, err := utils.ReadCookie(c.Request, CookieSession)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not read a cookie")
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not delete a session")
		return
	}

	utils.DeleteCookie(c.Writer, CookieSession)
	c.Redirect(http.StatusFound, "/signin")
	
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
