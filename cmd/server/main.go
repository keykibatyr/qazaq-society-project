package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project/internal/controllers"
	"github.com/keykibatyr/qazaq-society-project/internal/middleware"
	"github.com/keykibatyr/qazaq-society-project/internal/models"
	"github.com/keykibatyr/qazaq-society-project/internal/views"
)

var cfg models.PostgresConfig

func main() {
	r := gin.Default()

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

	filesUser := []string{
		"internal/views/layout.tmpl",
		"internal/views/navbar.tmpl",
		"internal/views/footer.tmpl",
	}

	filesAdmin := []string{
		"internal/views/admin/layout.tmpl",
		"internal/views/admin/navbar.tmpl",
		"internal/views/admin/footer.tmpl",
	}

	r.Static("/assets", "./assets")

	tplAbout := views.Must(views.ParseFileSys(append(filesUser, "internal/views/about.tmpl")))

	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB:            db,
		BytesPerToken: 32,
	}

	eventService := &models.EventService{
		DB: db,
	}

	adminC := controllers.Admins{
		UserService:     userService,
		SerssionService: sessionService,
		EventService:    eventService,
	}

	userC := controllers.Users{
		UserService:    userService,
		SessionService: sessionService,
		EventService:    eventService,
	}

	userC.Templates.New = views.Must(views.ParseFileSys(append(filesUser, "internal/views/auth/signup.tmpl")))

	userC.Templates.SignIn = views.Must(views.ParseFileSys(append(filesUser, "internal/views/auth/signin.tmpl")))

	userC.Templates.Events = views.Must(views.ParseFileSys(append(filesUser, "internal/views/events/index.tmpl")))

	userC.Templates.Home = views.Must(views.ParseFileSys(append(filesUser, "internal/views/home.tmpl")))

	UserMW := middleware.UserMiddleware{
		SessionService: sessionService,
		SignInPage:     "/signin",
		CookieName:     "session",
	}

	r.Use(UserMW.SetUser())

	auth := r.Group("/")

	auth.Use(UserMW.RequireUser())
	auth.GET("/users/me", userC.CurrentUserController)

	r.GET("/", userC.Home)

	r.GET("/about", func(c *gin.Context) {
		data := gin.H{
			"Title": "About Page",
		}
		controllers.Render(c, tplAbout, data)
	})

	r.GET("/events", userC.Events)

	r.GET("/signup", userC.New)
	r.POST("/signup", userC.Create)

	r.GET("/signin", userC.SignIn)
	r.POST("/signin", userC.ProcessSignIn)

	r.POST("/signout", userC.ProcessSignOut)

	tplTest := views.Must(views.ParseFileSys(append(filesAdmin, "internal/views/admin/test.tmpl")))

	adminC.Templates.Events = views.Must(views.ParseFileSys(append(filesAdmin, "internal/views/admin/events/index.tmpl")))
	adminC.Templates.EventsNew = views.Must(views.ParseFileSys(append(filesAdmin, "internal/views/admin/events/new.tmpl")))

	admin := r.Group("/admin")
	admin.Use(UserMW.RequireUser())
	admin.Use(middleware.RequireAdmin())

	admin.GET("/test", func(c *gin.Context) {
		data := gin.H{
			"Title": "Test",
		}
		controllers.Render(c, tplTest, data)
	})

	admin.GET("/events", adminC.Events)
	admin.GET("/events/new", adminC.EventsNew)
	admin.POST("/events/new", adminC.AddEvents)
	admin.POST("/events/:id/publish", adminC.PublishEvent)
	admin.POST("/events/:id/unpublish", adminC.UnPublishEvent)
	admin.POST("/events/:id/delete", adminC.DeleteEvent)

	r.Run(":8080")
}
