package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"ngini.com/test-api/internal/api"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	r.Route("/orders", func(r chi.Router) {
		r.With(api.Paginate).Get("/", api.ListOrders)
		r.Post("/", api.CreateOrder) // POST /orders
		//r.Get("/search", SearchArticles) // GET /orders/search

		r.Route("/{orderID}", func(r chi.Router) {
			r.Use(api.OrderCtx)            // Load the *Order on the request context
			r.Get("/", api.GetOrder)       // GET /orders/123
			r.Put("/", api.UpdateOrder)    // PUT /orders/123
			r.Delete("/", api.DeleteOrder) // DELETE /orders/123
		})

		// GET /articles/whats-up
		r.With(api.OrderCtx).Get("/{orderSlug:[a-z-]+}", api.GetOrder)
	})

	http.ListenAndServe(":80", r)
}
