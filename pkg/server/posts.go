package server

import (
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	itemdata "gitlab.com/vk-go/lectures-2022-2/pkg/repository/itemData"
	"gitlab.com/vk-go/lectures-2022-2/pkg/utils"
	"io"
	"net/http"
)

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(string)
	if !ok {
		utils.NewRespError(w, "invalid user id", 400, s.log)
		return
	}
	var buf []byte
	post := itemdata.CreatePost{}
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	if err = r.Body.Close(); err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	if err = utils.UnmarshalCreatPost(buf, &post); err != nil {
		utils.NewRespError(w, "invalid json input", 400, s.log)
		return
	}
	if _, err = govalidator.ValidateStruct(post); err != nil {
		utils.NewRespError(w, "invalid struct fields", 400, s.log)
		return
	}
	resPost, err := s.service.CreatePost(post, id)
	if err != nil {
		utils.NewRespError(w, err.Error(), 400, nil)
		return
	}
	resp, err := utils.MarshalPost(resPost)
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
	s.log.Printf("Successful post creating | userID %s | postID %s \n", id, resPost.ID)
}

func (s *Server) GetPosts(w http.ResponseWriter, _ *http.Request) {
	posts, err := s.service.GetPosts()
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	resp, err := utils.MarshalSlice(posts)
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
}

func (s *Server) GetCategory(w http.ResponseWriter, r *http.Request) {
	cat := mux.Vars(r)["category"]
	if !govalidator.IsIn(cat, "music", "funny", "videos", "programming", "news", "fashion", "all") {
		utils.NewRespError(w, "invalid category", 400, s.log)
		return
	}
	posts, err := s.service.GetCategory(cat)
	if err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
	resp, err := utils.MarshalSlice(posts)
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
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	usr := mux.Vars(r)["user_login"]
	posts, err := s.service.GetName(usr)
	if err != nil {
		utils.NewRespError(w, err.Error(), 400, s.log)
		return
	}
	resp, err := utils.MarshalSlice(posts)
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
}

func (s *Server) GetPostID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["post_id"]
	post, err := s.service.GetPostID(id)
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
	w.WriteHeader(200)
	if _, err = w.Write(resp); err != nil {
		utils.NewRespError(w, err.Error(), 500, s.log)
		return
	}
}

func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("id").(string)
	if !ok {
		utils.NewRespError(w, "invalid user id", 400, s.log)
		return
	}
	postID := mux.Vars(r)["post_id"]
	err := s.service.DeletePost(postID, userID)
	if err != nil {
		utils.NewRespError(w, err.Error(), 400, s.log)
		return
	}
	w.Header().Add("content-type", "application/json; charset=utf-8")
	utils.NewRespError(w, "success", 200, s.log)
	s.log.Printf("Successful post deleting | userID %s | postID %s", userID, postID)
}
