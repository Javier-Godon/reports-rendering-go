package render_full_pdf

type RenderFullPdfQuery struct {
	DateFrom int32 `json:"date_from" binding:"required"`
	DateTo   int32 `json:"date_to" binding:"required"`
}