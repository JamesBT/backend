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

	e.GET("/user", controler.GetAllUser)
	e.POST("/user", controler.SignUp)
	e.POST("/user/auth", controler.Login)
	e.GET("/user/:id", controler.GetUserById)
	e.PUT("/user", controler.UpdateUser)
	e.PUT("/user/auth", controler.UpdateUserFull)
	e.DELETE("/user/:id", controler.DeleteUserById)
	e.GET("/user/ktp/:id", controler.GetUserKTP)
	e.GET("/user/foto/:id", controler.GetUserKTP)
	e.GET("/user/unverified", controler.GetAllUserUnverified)
	e.GET("/user/perusahaan/:id", controler.GetAllUserByPerusahaanId)

	e.GET("/asset", controler.GetAllAsset)
	// e.GET("/asset/:nama", controler.GetAssetByName)
	e.GET("/asset/:id", controler.GetAssetById)
	e.PUT("/asset/:id", controler.UbahVisibilitasAset)
	e.POST("/asset", controler.TambahAsset)
	e.POST("/asset/child", controler.TambahAssetChild)
	e.GET("/asset/detail/:id", controler.GetAssetDetailedById)

	e.GET("/surveyor", controler.GetAllSurveyor)
	e.GET("/surveyor/user/:id", controler.GetSurveyorByUserId)
	e.PUT("/surveyor/surv", controler.UpdateUserBySurveyorId)
	e.PUT("/surveyor/user", controler.UpdateSurveyorByUserId)
	e.GET("/surveyor/:nama", controler.GetSurveyorByName)
	e.PUT("/surveyor", controler.UpdateSurveyorById)
	e.GET("/surveyor/detail", controler.GetAllSurveyorDetailed)
	e.POST("/surveyor", controler.SignUpSurveyor)
	e.POST("/surveyor/auth", controler.LoginSurveyor)

	e.GET("/survey_req/detail", controler.GetAllSurveyReqDetailed)
	e.GET("/survey_req/ongoing/:id", controler.GetAllOngoingSurveyReqByUserId)
	e.GET("/survey_req/finished/:id", controler.GetAllFinishedSurveyReqByUserId)
	e.GET("/survey_req/:id", controler.GetSurveyReqById)
	e.POST("/survey_req", controler.CreateSurveyReq)
	e.GET("/survey_req/aset/:id", controler.GetSurveyReqByAsetId)
	e.GET("/survey_req/user/:id", controler.GetAllSurveyReqByUserId)
	e.GET("/survey_req/aset/nama/:nama", controler.GetSurveyReqByAsetNama)

	e.POST("/perusahaan", controler.TambahPerusahaan)
	e.GET("/perusahaan/detail", controler.GetAllPerusahaanDetailed)
	e.GET("/perusahaan/detail/:id", controler.GetPerusahaanDetailById)
	e.GET("/perusahaan/home/:id", controler.HomeUserPerusahaan)
	e.GET("/perusahaan/user/:id", controler.GetPerusahaanByUserId)
	e.GET("/perusahaan/unverified", controler.GetAllPerusahaanUnverified)

	e.PUT("/verify/otp", controler.VerifyOTP)
	e.PUT("/verify/user/accept", controler.VerifyUserAccept)
	e.PUT("/verify/user/decline", controler.VerifyUserDecline)
	e.PUT("/verify/perusahaan/accept", controler.VerifyPerusahaanAccept)
	e.PUT("/verify/perusahaan/decline", controler.VerifyPerusahaanDecline)
	e.PUT("/verify/asset/accept", controler.VerifyAssetAccept)
	e.PUT("/verify/asset/reassign", controler.VerifyAssetReassign)

	e.GET("/tranreq/user/:id", controler.GetTranReqByUserId)
	e.GET("/tranreq/perusahaan/:id", controler.GetTranReqByPerusahaanId)
	return e
}
