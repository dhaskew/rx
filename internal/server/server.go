package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dhaskew/rx/internal/films"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type RouterFunc func() *chi.Mux

type Server struct {
	Logger         *zap.Logger
	Router         *chi.Mux
	MoviesDB       *sql.DB
	FilmRepository films.FilmRepository
	*http.Server
}

func NewServer(options ...func(*Server) *Server) *Server {
	server := &Server{
		Server: &http.Server{Addr: ":" + "80",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second},
	}

	for _, o := range options {
		o(server)
	}

	return server
}

func WithFilmRepository(rep *films.FilmRepository) func(*Server) *Server {
	return func(s *Server) *Server {
		s.FilmRepository = *rep
		return s
	}
}

func WithPort(port string) func(*Server) *Server {
	return func(s *Server) *Server {
		s.Addr = ":" + port
		return s
	}
}

func WithRouterFunc(routerFunc RouterFunc) func(*Server) *Server {
	return func(s *Server) *Server {
		s.Router = routerFunc()
		s.Server.Handler = s.Router
		return s
	}
}

func WithLogger(log *zap.Logger) func(*Server) *Server {
	return func(s *Server) *Server {
		s.Logger = log
		return s
	}
}

type ApiVersion struct{}

// provide a context with the api version for handler flexing
func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ApiVersion{}, version))
			next.ServeHTTP(w, r)
		})
	}
}

func (s Server) SetupRoutes() {
	//global middleware - all routes
	s.Router.Use(middleware.RequestID)
	s.Router.Use(ZapRequestLogger(s.Logger))
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Heartbeat("/ping"))
	s.Router.Use(middleware.Timeout(60 * time.Second))

	// if you want to use other buckets than the default (300, 1200, 5000) you can run:
	// m := negroniprometheus.NewMiddleware("serviceName", 400, 1600, 700)
	//m := NewMiddleware("serviceName")
	//s.Router.Use(m)

	// Utility Routes
	// /ping route provide by chi heartbeat middleware
	//s.Router.Get("/healthcheck", s.HealthcheckHandler())
	//s.Router.Get("/metrics", s.MetricsHandler())

	// API version 1.
	s.Router.Route("/v1", func(v1 chi.Router) {
		v1.Use(apiVersionCtx("v1"))
		v1.Mount("/films", func() http.Handler {
			v1Routes := chi.NewRouter()
			v1Routes.Get("/", s.filmsHandler())
			//v1Routes.With(EnsureJSONContentType).Post("/", s.CreateFilmCommentHandler())
			v1Routes.Get("/{filmID}", s.getFilmHandler())
			return v1Routes
		}())
	})

	s.Logger.Info("Done setting up routing")
}

func (s Server) PrintRoutes() {
	s.Logger.Info("Printing routes ...")
	_ = chi.Walk(s.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		s.Logger.Info("Route registered: " + "(" + method + "): " + route)
		return nil
	})
	s.Logger.Info("Done printing routes")
}

func (s Server) Start() {
	s.Logger.Info("Starting server ...")
	defer func() {
		_ = s.Logger.Sync()
	}()

	s.SetupRoutes()
	s.PrintRoutes()
	//s.SetupTracing()

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
			s.Logger.Fatal("Could not listen on", zap.String("addr", s.Addr), zap.Error(err))
		}
	}()

	s.Logger.Info("Server is ready to handle requests", zap.String("addr", s.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	s.Logger.Info("Server is shutting down", zap.String("reason", sig.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Fatal("Could not gracefuly shutdown the server", zap.Error(err))
	}
	s.Logger.Info("Server stopped")
}

func (s Server) filterByRatingIfPresent(w http.ResponseWriter, r *http.Request) bool {
	rating := r.URL.Query().Get("rating")
	s.Logger.Info("Rating: " + rating)
	if rating != "" {
		films, err := s.FilmRepository.GetAllByRating(r.Context(), rating)
		s.Logger.Info("Films: " + strconv.Itoa(len(films)))
		if err != nil {
			s.Logger.Error("Error getting films", zap.Error(err))
			http.Error(w, "Error getting films", http.StatusInternalServerError)
			return true
		}

		res, err := json.MarshalIndent(films, "", "\t")

		if err != nil {
			s.Logger.Error("Error marshalling films", zap.Error(err))
			http.Error(w, "Error marshalling films", http.StatusInternalServerError)
			return true
		}

		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write(res)

		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func (s Server) filterByCategoryIfPresent(w http.ResponseWriter, r *http.Request) bool {
	category := r.URL.Query().Get("category")
	s.Logger.Info("Category: " + category)
	if category != "" {
		films, err := s.FilmRepository.GetAllByCategory(r.Context(), category)
		s.Logger.Info("Films: " + strconv.Itoa(len(films)))
		if err != nil {
			s.Logger.Error("Error getting films", zap.Error(err))
			http.Error(w, "Error getting films", http.StatusInternalServerError)
			return true
		}

		res, err := json.MarshalIndent(films, "", "\t")

		if err != nil {
			s.Logger.Error("Error marshalling films", zap.Error(err))
			http.Error(w, "Error marshalling films", http.StatusInternalServerError)
			return true
		}

		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write(res)

		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func (s Server) filmsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.filterByRatingIfPresent(w, r) {
			return
		}

		if s.filterByCategoryIfPresent(w, r) {
			return
		}

		films, err := s.FilmRepository.GetAll(r.Context())
		if err != nil {
			s.Logger.Error("Error getting films", zap.Error(err))
			http.Error(w, "Error getting films", http.StatusInternalServerError)
			return
		}

		res, err := json.MarshalIndent(films, "", "\t")

		if err != nil {
			s.Logger.Error("Error marshalling films", zap.Error(err))
			http.Error(w, "Error marshalling films", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write(res)

		w.WriteHeader(http.StatusOK)
	}
}

func (s Server) getFilmHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filmIDparam := chi.URLParam(r, "filmID")

		filmID, err := strconv.Atoi(filmIDparam)
		if err != nil {
			s.Logger.Error("Error converting filmID to int", zap.Error(err))
			http.Error(w, "Error converting filmID to int", http.StatusInternalServerError)
			return
		}

		film, err := s.FilmRepository.GetByID(r.Context(), filmID)
		if err == films.ErrNotFound {
			http.Error(w, "Film Not Found", http.StatusNotFound)
			return
		} else if err != nil {
			s.Logger.Error("Error getting film", zap.Error(err))
			http.Error(w, "Error getting film", http.StatusInternalServerError)
			return
		}

		res, err := json.MarshalIndent(film, "", "\t")

		if err != nil {
			s.Logger.Error("Error marshalling film", zap.Error(err))
			http.Error(w, "Error marshalling film", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write(res)

		w.WriteHeader(http.StatusOK)
	}
}
