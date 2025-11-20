package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageInfo struct {
	Title   string
}

func main() {
	r := gin.Default()

	files := []string{
		"internal/views/layout.tmpl",
		"internal/views/navbar.tmpl",
		"internal/views/footer.tmpl",
		"internal/views/home.tmpl",
		"internal/views/about.tmpl",
		"internal/views/events/index.tmpl",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("could not parse the template: %v", err)
	}

	for _, t := range tpl.Templates() {
		log.Printf("template %q was loaded", t.Name())
	}

	r.SetHTMLTemplate(tpl)

	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		data := PageInfo{
			Title: "Qazaq Society In Belgium",
		}

		c.HTML(http.StatusOK, "home.tmpl", data)
	})

	r.GET("/about", func(c *gin.Context) {
		data := PageInfo{
			Title: "About Us",
		}

		c.HTML(http.StatusOK, "about.tmpl", data)
	})

	r.GET("/events", func(c *gin.Context) {
		data := PageInfo{
			Title: "Events",
		}

		// the parsed template for the events file registers under its base name
		// (index.tmpl) so use that when executing.
		c.HTML(http.StatusOK, "/events/index.tmpl", data)
	})


	log.Print("listening to :8080...")
	r.Run(":8080")
}
