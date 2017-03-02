package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type HandlerFunc http.HandlerFunc

// type HandlerFunc func(w http.ResponseWriter, r *http.Request, c Context) *Error

type Middleware negroni.HandlerFunc

// type Middleware func(w http.ResponseWriter, r *http.Request, next HandlerFunc)

type Route struct {
	*mux.Route
}

type Router struct {
	*mux.Router
	log     *Logger
	handler *negroni.Negroni
}

func NewRouter() *Router {
	return newRouterRoute(mux.NewRouter())
}

func newRouterRoute(router *mux.Router) *Router {
	return &Router{
		Router:  router,
		log:     NewLogger(),
		handler: negroni.New(),
	}
}

func (r *Router) HandleAssets(path, dir string) {
	p := fmt.Sprintf("%s/", path)
	r.PathPrefix(p).Handler(http.StripPrefix(p, http.FileServer(http.Dir(dir))))
}

func (r *Router) HandleRedirect(method, path, to string) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, 302)
	}).Methods(method)
}

// func (rt *Router) HandleFunc(path string, fn HandlerFunc) *Route {
//   return &Route{
//     Route: rt.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
//       for i := 0; i < len(rt.middleware)-1; i++ {
//         rt.middleware(w, r, rt.middleware[i+1])
//       }
//       last := fn
//       for i := len(rt.middleware) - 1; i >= 0; i-- {
//         rt.middleware[i](w, r, last)
//         last = rt.middleware[i]
//       }
//       fmt.Printf("stack = %+v\n", stack)
//       fmt.Println("fn")
//       fmt.Printf("rt = %+v\n", rt)
//       fmt.Printf("w = %+v\n", w)
//       fmt.Printf("r = %+v\n", r)
//       stack(w, r)
//     }),
//   }
// }

func (r *Router) HandleText(method, path, text string) {
	r.HandleFunc(path, textHandler(text)).Methods(method)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServeHTTP")
	fmt.Printf("rt = %+v\n", rt)
	fmt.Printf("rt.handler = %+v\n", rt.handler)
	rt.handler.UseHandler(rt.Router)
	rt.handler.ServeHTTP(w, r)
}

func (rt *Router) Subrouter() *Router {
	return newRouterRoute(rt.Router.PathPrefix("/").Subrouter())
}

func (rt *Router) Use(m Middleware) {
	rt.handler.UseFunc(m)
}

func (rt *Router) UseHandlerFunc(fn HandlerFunc) {
	rt.handler.UseHandlerFunc(fn)
}

func textHandler(s string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(s))
	}
}
