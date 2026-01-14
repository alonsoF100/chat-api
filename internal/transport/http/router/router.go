package router

import (
	"github.com/alonsoF100/chat-api/internal/transport/http/handlers"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Handlers *handlers.Handler
}

func New(handlers *handlers.Handler) *Router {
	return &Router{
		Handlers: handlers,
	}
}

func (rt Router) Setup() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/chats", func(r chi.Router) {
		r.Post("/", rt.Handlers.CreateChat)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", rt.Handlers.GetMessages)
			r.Delete("/", rt.Handlers.DeleteChat)
			r.Post("/messages/", rt.Handlers.CreateMessage)
		})
	})

	return r
}
