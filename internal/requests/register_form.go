package requests

type RegisterForm struct {
	Email string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required,min=8"`
	Name string `form:"first_name" binding:"required"`
	Surname string `form:"second_name" binding:"required"`
}