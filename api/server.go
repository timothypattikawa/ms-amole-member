package api

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/timothypattikawa/ms-kamoro-costumer/pkg/exception"
	"log"
)

func RunServer(memberHandler func(*echo.Echo), v *viper.Viper) {
	e := echo.New()

	e.HTTPErrorHandler = exception.CostumeEchoError

	e.HideBanner = true
	memberHandler(e)

	port := v.GetString("service.port")
	err := e.Start(":" + port)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
