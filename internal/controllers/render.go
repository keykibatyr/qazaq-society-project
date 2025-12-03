package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/keykibatyr/qazaq-society-project/internal/middleware"
)

func Render(c *gin.Context, tmpl Template, data gin.H) {
	if data == nil {
		data = gin.H{}
	}

	data["CurrentUser"] = middleware.CurrentUser(c)
	
    tmpl.ExecuteTemplate(c.Writer, c.Request, data)
}
