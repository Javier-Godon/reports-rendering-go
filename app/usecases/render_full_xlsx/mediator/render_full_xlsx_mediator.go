package mediator

import (
	"log"
	"github.com/Javier-Godon/reports-rendering-go/framework"
	"github.com/Javier-Godon/reports-rendering-go/usecases/render_full_xlsx"
)

func init() {
	err := framework.Register[render_full_xlsx.RenderFullXlsxQuery, render_full_xlsx.RenderFullXlsxResult](render_full_xlsx.NewRenderFullXlsxHandler())
	if err != nil {
		return
	}
}

func Send(query render_full_xlsx.RenderFullXlsxQuery) render_full_xlsx.RenderFullXlsxResult {
	RenderFullXlsxResult, err := framework.Send[render_full_xlsx.RenderFullXlsxQuery, render_full_xlsx.RenderFullXlsxResult](query)
	if err != nil {
		log.Fatalf("Could not execute: %v", query)
	}
	return RenderFullXlsxResult
}