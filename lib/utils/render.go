package utils

import (
	"net/http"

	"gopkg.in/unrolled/render.v1"
)

// Reply takes in a render instance, an http responsewriter, a message string, http status code
// and additional data to return to the user. This needs a "return" call after being run.
func Reply(r *render.Render, w http.ResponseWriter, message string, code int, data ...interface{}) {
	r.JSON(w, code, map[string]interface{}{
		"message": message,
		"code":    code,
		"data":    data,
	})
}
