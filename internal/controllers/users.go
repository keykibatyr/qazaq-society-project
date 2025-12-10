package controllers

import (
	"fmt"
	"net/http"
	"strconv"

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
		New          Template
		SignIn       Template
		Events       Template
		Home         Template
		AllElections Template
		Election Template
	}

	UserService      *models.UserService
	SessionService   *models.SessionService
	EventService     *models.EventService
	CandidateService *models.CandidateService
	ElectionService  *models.ElectionService
	VoteService      *models.VoteService
}

func (u Users) Home(c *gin.Context) {
	event, err := u.EventService.GetLatest()
	if err != nil {
		Render(c, u.Templates.Home, nil)
	}

	data := gin.H{
		"Title":    event.Title,
		"Time":     event.Time,
		"Month":    event.Month,
		"Day":      event.Day,
		"ImageURL": event.ImageURL,
		"Location": event.Location,
	}

	fmt.Println(data)

	Render(c, u.Templates.Home, data)
}

func (u Users) Events(c *gin.Context) {
	events, err := u.EventService.GetAll()
	if err != nil {
		c.String(500, "error")
		return
	}

	data := gin.H{
		"Events": events,
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

	if !(utils.ValidLen(password)) || !(utils.ValidPassword(password)) {
		Render(c, u.Templates.New, gin.H{
			"Error": "The Password must contain special charachters and longer be than 8 chars",
		})
		return
	}

	if !u.UserService.EmailCheck(email) {
		Render(c, u.Templates.New, gin.H{
			"Error": "That Email is already registered",
		})
		return
	}

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
	c.Redirect(http.StatusFound, "/")

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
		Render(c, u.Templates.SignIn, gin.H{
			"Error": "The Password or Email are Incorrect",
		})
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
	c.Redirect(http.StatusFound, "/")
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

func (u Users) AllElections(c *gin.Context) {
	elections, err := u.ElectionService.GetAllElections()
	if err != nil {
		c.String(500, "error")
		return
	}

	data := gin.H{
		"Elections": elections,
	}

	Render(c, u.Templates.AllElections, data)

}

func (u Users) Election(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(500, "error")
		return
	}
	candidates, err := u.CandidateService.GetAllCandidates(id)
	if err != nil {
		c.String(500, "error")
		return
	}

	data := gin.H{
		"Candidates": candidates,
	}

	Render(c, u.Templates.Election, data)
}
