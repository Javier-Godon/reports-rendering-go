package rest

type RenderFullPdfRequest struct {
	DateFrom int32 `json:"date_from" binding:"required"`
	DateTo   int32 `json:"date_to" binding:"required"`
}

type RenderFullPdfResponse struct {
	Payload string `json:"payload" binding:"required"`
}
