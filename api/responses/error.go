package responses

import (
	"github.com/go-chi/render"
	"net/http"
)

// Common response error
// swagger:response errorResponse
type ErrorResponse_ struct {
	// in: body
	Body ErrorResponse
}

type ErrorResponse struct {
	HTTPCode int    `json:"-"`
	Error    string `json:"error"`
	Debug    string `json:"debug"`
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPCode)
	return nil
}

func ErrInternal(err error) render.Renderer {
	return ErrorResponse{
		HTTPCode: http.StatusInternalServerError,
		Error:    "internal error",
		Debug:    err.Error(),
	}
}

func ErrBadRequest(err error) render.Renderer {
	return ErrorResponse{
		HTTPCode: http.StatusBadRequest,
		Error:    "bad request",
		Debug:    err.Error(),
	}
}
