package amadeus

type AmadeusError struct {
	Errors []struct {
		Code   int    `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Status int    `json:"status"`
	} `json:"errors"`
}
