package http

type GetCountryInput struct {
	Lang string `json:"lang"`
}

type ErrorCountryResponse struct {
	Error string `json:"error"`
}

type CountryResponse struct {
	Response string `json:"response"`
}
