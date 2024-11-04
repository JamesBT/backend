package main

import (
	"TemplateProject/route"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func createFileLog(namFile string) {

	// deteksi apakah file sudah ada
	var _, err = os.Stat(namFile)

	// buat file baru jika belum ada
	if os.IsNotExist(err) {
		file, err := os.Create(namFile)

		if err != nil {
			return
		}
		defer file.Close()
	}

}

func InsertLogs(string2 string) error {

	now := time.Now()
	namaFile := now.Format("20060102")
	ts := now.Format("2006-01-02 15:04:05")
	path := "log/server/" + namaFile + ".txt"
	createFileLog(path)
	format := "<" + ts + ">||" + string2
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	defer f.Close()
	//Error:<nil>
	if strings.Contains(format, "Error:<nil>") == false {
		InsertLogsError(format)
	}
	if _, err = f.WriteString(format); err != nil {

		log.Println(err)
		return err
	}

	return nil
}
func InsertLogsError(string2 string) error {

	now := time.Now()
	namaFile := now.Format("20060102")

	path := "log/error/" + namaFile + ".txt"
	createFileLog(path)
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	defer f.Close()
	//Error:<nil>
	if _, err = f.WriteString(string2); err != nil {

		log.Println(err)
		return err
	}

	return nil
}

func main() {
	e := route.Init()
	//e.Use(middleware.Recover())
	e.HideBanner = true
	//e.Use(middleware.Logger())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:        true,
		LogURI:           true,
		LogRemoteIP:      true,
		LogContentLength: true,
		LogUserAgent:     true,
		LogHost:          true,
		LogError:         true,
		LogLatency:       true,
		LogMethod:        true,
		LogProtocol:      true,
		LogRequestID:     true,
		LogResponseSize:  true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(int)
			logs := fmt.Sprintf(" {REQUEST: uri: %v, status: %v, ip:%s, Host:%v, ContentLength:%v, User-Agent:%v, Error:%v, Latency:%v, method:%v ,Protocol:%v, RequestId:%v, ResponseSize:%d, custom-value: %v}\n",
				v.URI, v.Status, v.RemoteIP, v.Host, v.ContentLength, v.UserAgent, v.Error, v.Latency, v.Method, v.Protocol, v.RequestID, v.ResponseSize, value)
			//log.Println(logs)
			InsertLogs(logs)
			return nil
		},
	}))
	//autoTLSManager := autocert.Manager{
	//	Prompt: autocert.AcceptTOS,
	//	// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
	//	Cache:      autocert.DirCache("/var/www/.cache"),
	//	HostPolicy: autocert.HostWhitelist("https://psc.gesit119.id"),
	//}
	//s := http.Server{
	//	Addr:    ":2681",
	//	Handler: e, // set Echo as handler
	//	TLSConfig: &tls.Config{
	//		//Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
	//		GetCertificate: autoTLSManager.GetCertificate,
	//		NextProtos:     []string{acme.ALPNProto},
	//	},
	//	//ReadTimeout: 30 * time.Second, // use custom timeouts
	//}
	//if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
	//	e.Logger.Fatal(err)
	//}
	e.Logger.Fatal(e.Start(":1762"))
}
