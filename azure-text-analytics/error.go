package function

import (
	"encoding/json"
)

func createErrorResponseBody(m string) []byte {
	context := &ErrorContext{
		Message: m,
	}

	e := &ErrorResponse{
		Error: context,
	}

	b, _ := json.Marshal(e)

	return b
}
