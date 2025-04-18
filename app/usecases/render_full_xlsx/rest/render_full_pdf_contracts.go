package rest

type RenderFullXlsxRequest struct {
	DateFrom int32 `json:"date_from" binding:"required"`
	DateTo   int32 `json:"date_to" binding:"required"`
}

type RenderFullXlsxResponse struct {
	Payload []byte `json:"payload" binding:"required"`
}
