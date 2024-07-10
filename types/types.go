package types

type Chirp struct {
	Body string `json:"body"`
}

type ReturnVals struct {
	Id          int    `json:"id"`
	Error       string `json:"error"`
	CleanedBody string `json:"cleaned_body"`
}
