package api

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"ngini.com/test-api/internal/dao"
	"ngini.com/test-api/internal/model"
	"strings"
)
import "github.com/go-chi/render"

type OrdersEndpoint struct {
	dao dao.DAO
}

func NewOrdersEndpoint(dao dao.DAO) *OrdersEndpoint {
	return &OrdersEndpoint{dao}
}

func (o *OrdersEndpoint) ListOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var list []*model.Order
	list, err = o.dao.GetOrders(ctx)
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
	if err := render.RenderList(w, r, NewOrderListResponse(list)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

// OrderCtx middleware is used to load an Order object from
// the URL parameters passed through as the request. In case
// the Order could not be found, we stop here and return a 404.
func (o *OrdersEndpoint) OrderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var order *model.Order
		var err error

		if orderID := chi.URLParam(r, "orderID"); orderID != "" {
			order, err = o.dao.GetOrder(ctx, orderID)
		} else if orderSlug := chi.URLParam(r, "orderSlug"); orderSlug != "" {
			order, err = o.dao.GetOrderBySlug(ctx, orderSlug)
		} else {
			_ = render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			_ = render.Render(w, r, ErrNotFound)
			return
		}

		if order == nil {
			_ = render.Render(w, r, ErrNotFound)
			return
		}

		ctx = context.WithValue(r.Context(), "order", order)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateOrder persists the posted Order and returns it
// back to the client as an acknowledgement.
func (o *OrdersEndpoint) CreateOrder(w http.ResponseWriter, r *http.Request) {
	data := &OrderRequest{}
	ctx := r.Context()

	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	order := data.Order
	returnedOrder, err := o.dao.AddOrder(ctx, *order)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewOrderResponse(returnedOrder))
}

// GetOrder returns the specific Order. You'll notice it just
// fetches the Order right off the context, as its understood that
// if we made it this far, the Order must be on the context. In case
// it's not due to a bug, then it will panic, and our Recoverer will save us.
func (o *OrdersEndpoint) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the order
	// context because this handler is a child of the OrderCtx
	// middleware. The worst case, the recoverer middleware will save us.
	order := r.Context().Value("order").(*model.Order)

	if err := render.Render(w, r, NewOrderResponse(order)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdateOrder updates an existing Order in our persistent store.
func (o *OrdersEndpoint) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	order := r.Context().Value("order").(*model.Order)
	ctx := r.Context()
	data := &OrderRequest{Order: order}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	order = data.Order
	returnedOrder, err := o.dao.UpdateOrder(ctx, order.ID, order)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_ = render.Render(w, r, NewOrderResponse(returnedOrder))
}

// DeleteOrder removes an existing Order from our persistent store.
func (o *OrdersEndpoint) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	// Assume if we've reach this far, we can access the order
	// context because this handler is a child of the OrderCtx
	// middleware. The worst case, the recoverer middleware will save us.
	order := r.Context().Value("order").(*model.Order)

	order, err = o.dao.DeleteOrder(ctx, order.ID)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_ = render.Render(w, r, NewOrderResponse(order))
}

// OrderRequest is the request payload for Order data model.
//
// NOTE: It's good practice to have well-defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type OrderRequest struct {
	*model.Order

	//ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *OrderRequest) Bind(r *http.Request) error {
	// a.Order is nil if no Order fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Order == nil {
		return errors.New("missing required Order fields")
	}

	if len(a.Order.ID) == 0 || len(a.Order.Name) == 0 {
		return errors.New("missing required order fields")
	}

	// just a post-process after a decode..
	//a.ProtectedID = ""                           // unset the protected ID
	a.Order.Name = strings.ToLower(a.Order.Name) // as an example, we down-case
	return nil
}

// OrderResponse is the response payload for the Article data model.
// See NOTE above in OrderRequest as well.
//
// In the ArticleResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type OrderResponse struct {
	*model.Order

	// We add another field to the response here, such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

func NewOrderResponse(order *model.Order) *OrderResponse {
	resp := &OrderResponse{Order: order}

	return resp
}

func NewOrderListResponse(orders []*model.Order) []render.Renderer {
	list := []render.Renderer{}
	for _, order := range orders {
		list = append(list, NewOrderResponse(order))
	}
	return list
}

func (rd *OrderResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

// Paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub: some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
