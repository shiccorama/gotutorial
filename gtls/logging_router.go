package main

import (
	"log"
	"net/http"
)

type LoggingRouter struct {
	adminHandler *AdminHandler
	linksHandler *LinksHandler
}

type BaseHandler struct {
	constructorCanary bool
	linksDB           map[string]string
}

func MakeBaseHandler(linksDB map[string]string) BaseHandler {
	return BaseHandler{
		constructorCanary: true,
		linksDB:           linksDB,
	}
}

func NewLoggingRouter() *LoggingRouter {
	lr := new(LoggingRouter)
	linksDB := map[string]string{}
	lr.adminHandler = NewAdminHandler(linksDB)
	lr.linksHandler = NewLinksHandler(linksDB)
	return lr
}

type SimpleRouteHandler interface {
	Respond(req *http.Request) (statusCode int, headers map[string]string, responseBytes []byte)
}

// At a high level, a router inspects a request and routes it to an appropriate subcomponent for handling.
// Here, we just look for a simple prefix
func (router *LoggingRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	url := req.URL.Path
	var code int
	responseBytes := []byte{}
	extra_headers := map[string]string{}

	switch {
	case url == "/":
		code = http.StatusOK
		responseBytes = []byte("Welcome to gtls")
	case url == "/admin/add":
		// use the admin handler
		code, extra_headers, responseBytes = router.adminHandler.AddShortlinkFormResponse(req)
	case url == "/admin/post":
		// use the admin handler
		code, extra_headers, responseBytes = router.adminHandler.PostResponse(req)
	default:
		// use the shortlink handler
		code, extra_headers, responseBytes = router.linksHandler.Respond(req)
	}

	for k, v := range extra_headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(code)
	writtenCount, err := w.Write(responseBytes)
	if err != nil {
		log.Println("error writing response", req, err)
	}
	log.Printf("%s", CommonLogFormat(req, code, writtenCount))
}
