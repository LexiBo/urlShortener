package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inspectorvitya/shortenerURL/internal/app"
	"github.com/inspectorvitya/shortenerURL/internal/storage"
	"go.uber.org/zap"
)

type request struct {
	FullURL string `json:"fullURL"`
}
type response struct {
	ShortURL string `json:"shortUrl"`
}

func (s *Server) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	req := &request{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.FullURL == "" {
		newErrorResponse(w, http.StatusBadRequest, "empty url")
		return
	}
	shortURL, err := s.App.CreateShortURL(context.TODO(), req.FullURL)
	if err != nil {
		if errors.Is(err, app.ErrInvalidURL) || errors.Is(err, app.ErrExistURL) {
			zap.L().Info("handler create url: ", zap.Error(err))
			newErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			zap.L().Error("handler create url: ", zap.Error(err))
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	resp := response{ShortURL: r.Host + "/" + shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		zap.L().Error("handler create url: ", zap.Error(err))
	}
}

func (s *Server) GetFullURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fullURL, err := s.App.GetFullURL(context.TODO(), vars["shortURL"])
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			zap.L().Info("handler create url: ", zap.Error(err))
			newErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			zap.L().Error("handler create url: ", zap.Error(err))
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	//http.Redirect(w, r, fullURL, http.StatusFound) //
	_, _ = io.WriteString(w, fullURL)
}
