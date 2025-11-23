package main

import (
	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project.git/internal/views"
)

type PageInfo struct {
	Title string
}

func main() {
	r := gin.Default()

	files := []string{
		"internal/views/layout.tmpl",
		"internal/views/navbar.tmpl",
		"internal/views/footer.tmpl",
	}
	
	r.Static("/assets", "./assets")

	tplHome := views.Must(views.ParseFileSys(append(files, "internal/views/home.tmpl")))

	tplAbout:= views.Must(views.ParseFileSys(append(files, "internal/views/about.tmpl")))

	tplEvents:= views.Must(views.ParseFileSys(append(files, "internal/views/events/index.tmpl")))


	r.GET("/", func(c *gin.Context) {
		tplHome.ExecuteTemplate(c.Writer, c.Request, PageInfo{
        Title: "Home Page",
    	})
	})

	r.GET("/about", func(c *gin.Context) {
		tplAbout.ExecuteTemplate(c.Writer, c.Request, PageInfo{
			Title: "About Page",
		})
	})

	r.GET("/events", func(c *gin.Context) {
		tplEvents.ExecuteTemplate(c.Writer, c.Request, PageInfo{
			Title: "Events Page",
		})
	})

	r.Run(":8080")
}
