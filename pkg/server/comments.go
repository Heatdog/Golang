package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
	"io"
	"net/http"
)

func (s *Server) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post_id"]
	userID, ok := r.Context().Value("id").(string)
	if !ok {
		utils.NewRespError(w, "invalid user id", 400, s.log)
		return
	}
	var data []byte
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	if err = r.Body.Close(); err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	text := struct {
		Comment string `json:"comment"`
	}{}
	if err = json.Unmarshal(data, &text); err != nil {
		utils.NewRespError(w, "invalid json input", 400, s.log)
		return
	}
	post, err := s.service.CreateComm(postID, userID, text.Comment)
	if err != nil {
		utils.NewRespError(w, err.Error(), 400, nil)
		return
	}
	resp, err := utils.MarshalPost(post)
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	if _, err = w.Write(resp); err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	s.log.Printf("Successful comment creating | userID %s | postID %s \n", userID, postID)
}

func (s *Server) DeleteComment(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["post_id"]
	commID := mux.Vars(r)["comment_id"]
	userID, ok := r.Context().Value("id").(string)
	if !ok {
		utils.NewRespError(w, "invalid user id", 400, s.log)
		return
	}
	post, err := s.service.DeleteComm(postID, userID, commID)
	if err != nil {
		utils.NewRespError(w, err.Error(), 400, s.log)
		return
	}
	resp, err := utils.MarshalPost(post)
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
	s.log.Printf("Successful comment deleting | userID %s | postID %s | commentID %s \n", userID, postID, commID)
}

func (s *Server) Upvote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(string)
	if !ok {
		utils.NewRespError(w, "invalid user id", 400, s.log)
		return
	}
	postID := mux.Vars(r)["post_id"]
	post, err := s.service.Upvote(postID, userID)
	if err != nil {
		utils.NewRespError(w, err.Error(), 400, s.log)
		return
	}
	resp, err := utils.MarshalPost(post)
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
	s.log.Printf("Successful post upvote | userID %s | postID %s \n", userID, postID)
}

func (s *Server) Downvote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(string)
	if !ok {
		utils.NewRespError(w, "invalid user id", 500, s.log)
		return
	}
	postID := mux.Vars(r)["post_id"]
	post, err := s.service.Downvote(postID, userID)
	if err != nil {
		utils.NewRespError(w, err.Error(), 400, s.log)
		return
	}
	resp, err := utils.MarshalPost(post)
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
	s.log.Printf("Successful post downvote | userID %s | postID %s \n", userID, postID)
}

func (s *Server) Unvote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(string)
	if !ok {
		utils.NewRespError(w, "invalid user id", 400, s.log)
		return
	}
	postID := mux.Vars(r)["post_id"]
	post, err := s.service.Unvote(postID, userID)
	if err != nil {
		utils.NewRespError(w, err.Error(), 400, s.log)
		return
	}
	resp, err := utils.MarshalPost(post)
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
	s.log.Printf("Successful post unvote | userID %s | postID %s \n", userID, postID)
}
