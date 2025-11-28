package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project.git/internal/controllers"
	"github.com/keykibatyr/qazaq-society-project.git/internal/middleware"
	"github.com/keykibatyr/qazaq-society-project.git/internal/models"
	"github.com/keykibatyr/qazaq-society-project.git/internal/views"
)

type PageInfo struct {
	Title string
}

var cfg models.PostgresConfig

func main() {
	r := gin.Default()

	fmt.Print("bug")

	cfg = models.DefaultConfig()

	db, err := models.Open(cfg)
	if err != nil {
		log.Fatalf("could not open the DB: %v", err)
	}

	defer db.Close()

	err = models.Migrate(db, "./migrations/sql")
	if err != nil {
		log.Fatalf("could not migrate: %v", err)
	}

	files := []string{
		"internal/views/layout.tmpl",
		"internal/views/navbar.tmpl",
		"internal/views/footer.tmpl",
	}

	r.Static("/assets", "./assets")

	tplHome := views.Must(views.ParseFileSys(append(files, "internal/views/home.tmpl")))

	tplAbout := views.Must(views.ParseFileSys(append(files, "internal/views/about.tmpl")))

	tplEvents := views.Must(views.ParseFileSys(append(files, "internal/views/events/index.tmpl")))

	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
		BytesPerToken: 32,
	}

	userC := controllers.Users{
		UserService: userService,
		SessionService: sessionService,
	}

	userC.Templates.New = views.Must(views.ParseFileSys(append(files, "internal/views/auth/signup.tmpl")))

	userC.Templates.SignIn = views.Must(views.ParseFileSys(append(files, "internal/views/auth/signin.tmpl")))

	UserMW := middleware.UserMiddleware{
		SessionService: sessionService,
		SignInPage: "/signin",
		CookieName: "session",
	}

	r.Use(UserMW.SetUser())

	auth := r.Group("/")
	
	auth.Use(UserMW.RequireUser())
	{
		r.GET("/users/me", userC.CurrentUserController)
	}

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

	r.GET("/signup", userC.New)

	r.POST("/signup", userC.Create)

	r.GET("/signin", userC.SignIn)

	r.POST("/signin", userC.ProcessSignIn)


	r.Run(":8080")
}
