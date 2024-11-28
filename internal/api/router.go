package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
	"ngini.com/test-api/internal/dao"
)

func SetUpRouter(dao dao.DAO) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("root."))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	o := NewOrdersEndpoint(dao)
	r.Route("/orders", func(r chi.Router) {
		r.With(Paginate).Get("/", o.ListOrders)
		r.Post("/", o.CreateOrder) // POST /orders
		//r.Get("/search", SearchArticles) // GET /orders/search

		r.Route("/{orderID}", func(r chi.Router) {
			r.Use(o.OrderCtx)            // Load the *Order on the request context
			r.Get("/", o.GetOrder)       // GET /orders/123
			r.Put("/", o.UpdateOrder)    // PUT /orders/123
			r.Delete("/", o.DeleteOrder) // DELETE /orders/123
		})

		// GET /articles/whats-up
		r.With(o.OrderCtx).Get("/{orderSlug:[a-z-]+}", o.GetOrder)
	})
	return r
}
