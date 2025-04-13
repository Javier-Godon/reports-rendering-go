package rest

type RenderFullXlsxRequest struct {
	DateFrom int32 `json:"date_from" binding:"required"`
	DateTo   int32 `json:"date_to" binding:"required"`
}

type RenderFullXlsxResponse struct {
	Payload string `json:"payload" binding:"required"`
}
