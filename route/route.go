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

	e.GET("/kirimAngka/:angka", controler.KirimAngka)
	e.GET("/barang", controler.GetAllBarang)
	e.GET("/barang/:id", controler.GetBarangById)

	e.POST("/barang", controler.InsertBarang)

	e.POST("/user", controler.Login)

	return e
}
