package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type RespError struct {
	Body string `json:"message"`
}

type RegisterErrorList struct {
	List []RegisterError `json:"errors"`
}

type RegisterError struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

func NewRespError(w http.ResponseWriter, text string, statusCode int, log *log.Logger) {
	res, err := json.Marshal(RespError{Body: text})
	w.WriteHeader(statusCode)
	if err != nil {
		if log != nil {
			log.Println("json error marshal")
		}
		return
	}
	if _, err = w.Write(res); err != nil {
		if log != nil {
			log.Println("write error")
		}
		return
	}
	if log != nil {
		log.Println(text)
	}
}

func NewRegisterError(w http.ResponseWriter, list RegisterErrorList, statusCode int) {
	res, err := json.Marshal(list)
	w.WriteHeader(statusCode)
	if err != nil {
		return
	}
	if _, err = w.Write(res); err != nil {
		return
	}
}
