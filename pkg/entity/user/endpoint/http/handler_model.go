package http

type CalculateTripInput struct {
	Dates string `json:"dates"`
}

type UpdateLangInput struct {
	Lang string `json:"lang"`
}

type UpdateCountryInput struct {
	Code          string `json:"code"`
	Flag          string `json:"flag"`
	DaysContinual int    `json:"daysContinual"`
	DaysLimit     int    `json:"daysLimit"`
	ResetInterval int    `json:"resetInterval"`
}

type ErrorUserResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	Response string `json:"response"`
}
