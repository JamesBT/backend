package route

import (
	"TemplateProject/controler"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/getPhoto", controler.GetFoto)
	e.POST("/uploadFoto", controler.UploadFoto)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "welcome here")
	})

	e.POST("/user", controler.SignUp)
	e.POST("/user/auth", controler.Login)
	e.GET("/user/:id", controler.GetUserById)
	e.PUT("/user", controler.UpdateUser)
	e.DELETE("/user/:id", controler.DeleteUserById)
	return e
}
