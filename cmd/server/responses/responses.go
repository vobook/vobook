package responses

type Error struct {
	Error string `json:"error"`
}

type Success struct {
	Message string `json:"message"`
}

func OK(message string) Success {
	return Success{
		Message: message,
	}
}
