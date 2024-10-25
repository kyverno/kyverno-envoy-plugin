package handlers

import (
	"context"
	"net/http"
)

func HttpError(ctx context.Context, writer http.ResponseWriter, request *http.Request, err error, code int) {
	// logger.Error(err, "an error has occurred", "url", request.URL.String())
	http.Error(writer, err.Error(), code)
}
