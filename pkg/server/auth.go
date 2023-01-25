package server

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	userdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/userData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
	"io"
	"net/http"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var bd []byte
	bd, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	if err = r.Body.Close(); err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	var elem userdata.User
	if err = json.Unmarshal(bd, &elem); err != nil {
		utils.NewRespError(w, "invalid json input", 400, s.log)
		return
	}
	if ok, err := govalidator.ValidateStruct(elem); !ok || err != nil {
		utils.NewRespError(w, "invalid struct fields", 400, s.log)
		return
	}
	login := elem.Login
	elem, err = s.service.CreateUser(elem)
	if err != nil {
		utils.NewRegisterError(w, utils.RegisterErrorList{
			List: []utils.RegisterError{
				{
					Location: "body",
					Param:    "username",
					Value:    login,
					Msg:      err.Error(),
				},
			},
		}, 422)
		return
	}
	token, err := s.service.GenerateToken(elem.Login, elem.Password)
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	resp, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(resp); err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	s.log.Printf("Successful registration | login %s \n", elem.Login)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var bd []byte
	bd, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	if err = r.Body.Close(); err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	var elem userdata.User
	if err = json.Unmarshal(bd, &elem); err != nil {
		utils.NewRespError(w, "invalid json input", 400, s.log)
		return
	}
	if ok, err := govalidator.ValidateStruct(elem); !ok || err != nil {
		utils.NewRespError(w, "invalid struct fields", 400, s.log)
		return
	}
	token, err := s.service.GenerateToken(elem.Login, s.service.GetHash(elem.Password))
	if err != nil {
		utils.NewRespError(w, err.Error(), 401, s.log)
		return
	}
	resp, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	if _, err = w.Write(resp); err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	s.log.Printf("Successful login | login %s \n", elem.Login)
}
