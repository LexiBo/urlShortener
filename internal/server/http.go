package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/inspectorvitya/shortenerURL/internal/app"
)

type Server struct {
	App        *app.App
	router     *mux.Router
	HTTPServer *http.Server
}

func New(port string, app *app.App) *Server {
	router := mux.NewRouter()

	server := &Server{
		HTTPServer: &http.Server{
			Addr:         net.JoinHostPort("", port),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      router,
		},
		router: router,
		App:    app,
	}

	return server
}

func (s *Server) Start() error {
	s.router.HandleFunc("/", s.CreateShortURL).Methods(http.MethodPost)
	s.router.HandleFunc("/{shortURL}", s.GetFullURL).Methods(http.MethodGet)
	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.HTTPServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
