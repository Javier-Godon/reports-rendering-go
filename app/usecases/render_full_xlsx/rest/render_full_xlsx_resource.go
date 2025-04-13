package rest

import (
	"net/http"
	"github.com/Javier-Godon/reports-rendering-go/usecases/render_full_xlsx"
	"github.com/Javier-Godon/reports-rendering-go/usecases/render_full_xlsx/mediator"
	"github.com/gin-gonic/gin"
) 

func RouteRenderFullXlsx(route *gin.Engine)(routes gin.IRoutes){
	RenderFullXlsxRoute := route.POST("/render/xlsx/", func(ctx *gin.Context) {
		var request RenderFullXlsxRequest
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		RenderFullXlsxResult := mediator.Send(buildRenderFullXlsxQuery(request))
		ctx.JSON(http.StatusOK, fromRenderFullXlsxResultToResponse(RenderFullXlsxResult))
	} )
	return RenderFullXlsxRoute

}

func fromRenderFullXlsxResultToResponse(result render_full_xlsx.RenderFullXlsxResult) RenderFullXlsxResponse {
	return RenderFullXlsxResponse{
		Payload: result.Payload,
	}
}

func buildRenderFullXlsxQuery(request RenderFullXlsxRequest) render_full_xlsx.RenderFullXlsxQuery {
	return render_full_xlsx.RenderFullXlsxQuery{
		DateFrom:        request.DateFrom,
		DateTo: request.DateTo,
	}
}

//https://stackoverflow.com/questions/42967235/golang-gin-gonic-split-routes-into-multiple-files
//https://www.youtube.com/watch?v=BkAoT2XZM24