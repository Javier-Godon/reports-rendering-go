package render_full_pdf

type RenderFullPdfResult struct {
	Payload string `json:"payload" binding:"required"`
}