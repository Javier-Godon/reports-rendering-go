package render_full_xlsx

type RenderFullXlsxResult struct {
	Payload []byte `json:"payload" binding:"required"`
}