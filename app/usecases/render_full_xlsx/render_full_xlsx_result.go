package render_full_xlsx

type RenderFullXlsxResult struct {
	Payload string `json:"payload" binding:"required"`
}