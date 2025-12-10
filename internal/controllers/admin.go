package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project/internal/models"
)

type Admins struct {
	Templates struct {
		Events    Template
		EventsNew Template
		Edit Template
		Elections Template
		ElectionsNew Template
	}

	UserService     *models.UserService
	SerssionService *models.SessionService
	EventService    *models.EventService
	ElectionService *models.ElectionService
	CandidateService *models.CandidateService
	VoteService *models.VoteService
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
	location := c.PostForm("location")

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

	_, err = a.EventService.CreateEvent(title, description, dst, location, strdate, published)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not create the event")
		return
	}

	c.Redirect(http.StatusFound, "/admin/events")

}

func (a Admins) PublishEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot extract id")
		return
	}

	err = a.EventService.Publish(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot publish the event")
		return
	}

	c.Redirect(http.StatusFound, "/admin/events")
}

func (a Admins) UnPublishEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot extract id")
		return
	}

	err = a.EventService.UnPublish(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot publish the event")
		return
	}

	c.Redirect(http.StatusFound, "/admin/events")
}

func (a Admins) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot extract id")
		return
	}

	err = a.EventService.Delete(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot publish the event")
		return
	}

	c.Redirect(http.StatusFound, "/admin/events")
}

func (a Admins) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
 
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot extract id")
		return
	}

	event, err := a.EventService.GetEventByID(id)

	if err != nil {
		c.String(http.StatusInternalServerError, "cannot get the event")
		return
	}

	date := event.StartDate.Format("2006-01-02")
	time := event.StartDate.Format("15:04")

	data := gin.H{
		"Event": event,
		"Date": date,
		"Time": time,
	}
	Render(c, a.Templates.Edit, data)
}


func (a Admins) ProcessUpdateEvent(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	date := c.PostForm("date")
	clock := c.PostForm("time")
	image := c.PostForm("image_url")
	published := c.PostForm("published") == "true"
	location := c.PostForm("location")

	id, err := strconv.Atoi(c.Param("id"))
	
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot extract id")
		return
	}


	dt := date + " " + clock
	strdate, err := time.Parse("2006-01-02 15:04", dt)
	if err != nil {
		c.String(http.StatusInternalServerError, "oops could not parse the time")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		err = a.EventService.Update(id, title, description, image, location, strdate, published)
		if err != nil {
			c.String(http.StatusInternalServerError, "updating error")
			return
		}
	} else {
		dst := "./assets/uploads/" + file.Filename

		err = c.SaveUploadedFile(file, dst)
			if err != nil {
				c.String(http.StatusInternalServerError, "oops could not upload an image")
				return
			}
		
		err = a.EventService.Update(id, title, description, dst, location, strdate, published)
			if err != nil {
				c.String(http.StatusInternalServerError, "updating error")
				return
			}
	}

	c.Redirect(http.StatusFound, "/admin/events")
}

func (a Admins) Elections(c *gin.Context) {
	Render(c, a.Templates.Elections, nil)
}

func(a Admins) ElectionsNew(c *gin.Context) {
	Render(c, a.Templates.ElectionsNew, nil)
}

func (a Admins) ProcessElectionsNew(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	startDate := c.PostForm("start_date")
	endDate := c.PostForm("end_date")
	published := c.PostForm("published") == "true"

	names := c.PostFormArray("candidates[]")

	form, _:= c.MultipartForm()
	files := form.File["images[]"]

	if len(files) != len(names) {
		c.String(http.StatusInternalServerError, "names != files")
		return
	}

	layout := "2006-01-02T15:04"

	str, err := time.Parse(layout, startDate)
	if err != nil {
		c.String(http.StatusInternalServerError, "could not conver the time")
		return
	}


	end, err := time.Parse(layout, endDate)
	if err != nil {
		c.String(http.StatusInternalServerError, "could not conver the time")
		return
	}


	election, err := a.ElectionService.Create(title, description, str, end, published)
	if err != nil {
		c.String(http.StatusInternalServerError, "could not create the election")
		return
	}


	for i := range names {
		name := names[i]
		file := files[i]

		dst := "./assets/uploads/" + file.Filename

		fmt.Println(dst)

		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			c.String(http.StatusInternalServerError, "could not save the candidate photo")
			return
		}

		

		_, err := a.CandidateService.CreateCandidate(name, dst, election.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, "could not create a candidate")
			return
		}

	}
	
	c.Redirect(http.StatusFound, "/admin/elections")
}

