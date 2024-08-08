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
	e.PUT("/user/auth", controler.UpdateUserFull)
	e.DELETE("/user/:id", controler.DeleteUserById)
	e.GET("/user/ktp/:id", controler.GetUserKTP)
	e.GET("/user/foto/:id", controler.GetUserKTP)

	e.GET("/asset", controler.GetAllAsset)
	e.POST("/asset", controler.TambahAsset)

	return e
}
