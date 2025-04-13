package main

import (
	"log"
	"github.com/Javier-Godon/reports-rendering-go/framework"
	rendeRFullPdf "github.com/Javier-Godon/reports-rendering-go/usecases/render_full_pdf/rest"
	renderFullXlsx "github.com/Javier-Godon/reports-rendering-go/usecases/render_full_xlsx/rest"

	"github.com/gin-gonic/gin"
)

func main() {
	framework.ReadConfig()
	serverPort := framework.AppConfig.ServerPort.PORT

	router := gin.Default()
	rendeRFullPdf.RouteRenderFullPdf(router)
	renderFullXlsx.RouteRenderFullXlsx(router)

	err := router.Run("0.0.0.0:" + serverPort)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
