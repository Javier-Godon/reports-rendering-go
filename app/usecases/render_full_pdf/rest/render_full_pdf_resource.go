package rest

import (
	"net/http"
	"github.com/Javier-Godon/reports-rendering-go/usecases/render_full_pdf"
	"github.com/Javier-Godon/reports-rendering-go/usecases/render_full_pdf/mediator"
	"github.com/gin-gonic/gin"
) 

func RouteRenderFullPdf(route *gin.Engine)(routes gin.IRoutes){
	RenderFullPdfRoute := route.POST("/render/pdf/", func(ctx *gin.Context) {
		var request RenderFullPdfRequest
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		RenderFullPdfResult := mediator.Send(buildRenderFullPdfQuery(request))
		ctx.JSON(http.StatusOK, fromRenderFullPdfResultToResponse(RenderFullPdfResult))
	} )
	return RenderFullPdfRoute

}

func fromRenderFullPdfResultToResponse(result render_full_pdf.RenderFullPdfResult) RenderFullPdfResponse {
	return RenderFullPdfResponse{
		Payload: result.Payload,
	}
}

func buildRenderFullPdfQuery(request RenderFullPdfRequest) render_full_pdf.RenderFullPdfQuery {
	return render_full_pdf.RenderFullPdfQuery{
		DateFrom:        request.DateFrom,
		DateTo: request.DateTo,
	}
}

//https://stackoverflow.com/questions/42967235/golang-gin-gonic-split-routes-into-multiple-files
//https://www.youtube.com/watch?v=BkAoT2XZM24