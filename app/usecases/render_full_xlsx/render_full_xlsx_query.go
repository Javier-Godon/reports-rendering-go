package render_full_xlsx

type RenderFullXlsxQuery struct {
	DateFrom int32 `json:"date_from" binding:"required"`
	DateTo   int32 `json:"date_to" binding:"required"`
}