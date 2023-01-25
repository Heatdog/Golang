package middleware

import (
	"context"
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
	"net/http"
	"strings"
)

func (mid *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		list := strings.Split(request.Header.Get("Authorization"), " ")
		if len(list) != 2 {
			utils.NewRespError(writer, "not enough arg in token", 400, mid.logger)
			return
		}
		header := list[1]
		if header == "" {
			utils.NewRespError(writer, "empty token", 400, mid.logger)
			return
		}
		id, err := mid.session.Check(header)
		if err != nil {
			utils.NewRespError(writer, err.Error(), 400, mid.logger)
			return
		}
		ctx := context.WithValue(request.Context(), "id", id)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
