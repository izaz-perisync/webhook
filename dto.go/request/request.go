package request

import (
	"net/http"

	"github.com/gorilla/schema"
)

func BindBody(r *http.Request, data any) error {

	if r.ContentLength == 0 {
		return nil
	}
	// err := r.ParseMultipartForm(32 << 20) // 32 MB
	// if err != nil {
	// 	return service.Generic{
	// 		Code:   2,
	// 		Msg:    "unable to parse form-data",
	// 		Source: err.Error(),
	// 	}
	// }

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return err
	}

	// decoder := schema.NewDecoder()
	return schema.NewDecoder().Decode(data, r.Form)

}
