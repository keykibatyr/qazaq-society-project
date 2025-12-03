package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project/internal/models"
)

type Admins struct {
	Templates struct {
		Events    Template
		EventsNew Template
	}

	UserService     *models.UserService
	SerssionService *models.SessionService
	EventService    *models.EventService
}

func (a Admins) Events(c *gin.Context) {
	events, err := a.EventService.GetAll()
	if err != nil {
    	c.String(500, "error")
    	return
	}

	data := gin.H{
    "Events": events,
}
	Render(c, a.Templates.Events, data)
}

func (a Admins) EventsNew(c *gin.Context) {
		events, err := a.EventService.GetAll()
	if err != nil {
    c.String(500, "error")
    return
	}

	data := gin.H{
    "Events": events,
}
	Render(c, a.Templates.EventsNew, data)
}


func (a Admins) AddEvents(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	date := c.PostForm("date")
	clock := c.PostForm("time")
	published := c.PostForm("published") == "true"

	dt := date + " " + clock
	strdate, err := time.Parse("2006-01-02 15:04", dt)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not parse the time")
		return
	}

	fmt.Println(strdate)

	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not extract the file")
		return
	}
	dst := "./assets/uploads/" + file.Filename

	fmt.Println(dst)

	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not upload an image")
		return
	}

	event, err := a.EventService.CreateEvent(title, description, dst, strdate, published)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not create the event")
		return
	}

	c.String(http.StatusOK, "Success")
	fmt.Println(event)
}
