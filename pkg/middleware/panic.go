package middleware

import (
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
	"net/http"
)

func (mid Middleware) Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				utils.NewRespError(writer, "internal server error", 500, nil)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
