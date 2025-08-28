package response

type ResponseModel struct {
	Object     interface{} `json:"data"`
	Pagination interface{} `json:"data,omitempty"`
	StatusHttp int         `json:"status"`
}
