package requests

type LoginForm struct {
	Email string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required,min=8"`
}