package api

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
)

type FaceUpdateRequest struct {
	Count int `json:"count"`
}

func (s Server) Faces(ctx echo.Context) error {
	r := new(FaceUpdateRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("Face Count [%d]\n", r.Count)

	return ctx.String(http.StatusOK, "")
}
