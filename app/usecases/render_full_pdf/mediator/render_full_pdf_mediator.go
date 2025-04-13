package mediator

import (
	"log"
	"github.com/Javier-Godon/reports-rendering-go/framework"
	"github.com/Javier-Godon/reports-rendering-go/usecases/render_full_pdf"
)

func init() {
	err := framework.Register[render_full_pdf.RenderFullPdfQuery, render_full_pdf.RenderFullPdfResult](render_full_pdf.NewRenderFullPdfHandler())
	if err != nil {
		return
	}
}

func Send(command render_full_pdf.RenderFullPdfQuery) render_full_pdf.RenderFullPdfResult {
	RenderFullPdfResult, err := framework.Send[render_full_pdf.RenderFullPdfQuery, render_full_pdf.RenderFullPdfResult](command)
	if err != nil {
		log.Fatalf("Could not execute: %v", command)
	}
	return RenderFullPdfResult
}