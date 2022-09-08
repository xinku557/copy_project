package controllers

import (
	"log"

	"github.com/gorilla/mux"
	"sheinko.tk/copy_project/middlewares"
)

func (handler *Handler) initializeRoutes() {
	log.Println("Initializing the routes...")

	handler.Router = mux.NewRouter()
	handler.Router.Use(middlewares.SetLoggingMiddleware)
	handler.Router.Use(middlewares.SetMiddlewareJSON)

	apiRoute := handler.Router.PathPrefix("/api/v1").Subrouter()
	apiRoute.HandleFunc("/register", handler.handleRegister).Methods("POST")
	apiRoute.HandleFunc("/login", handler.handleLogin).Methods("POST")
	apiRoute.HandleFunc("/profile", middlewares.SetMiddlewareAuthentication(handler.handleMe)).Methods("GET")
	apiRoute.HandleFunc("/profile", middlewares.SetMiddlewareAuthentication(handler.handleUpdateMe)).Methods("PUT")
	apiRoute.HandleFunc("/profile/posts", middlewares.SetMiddlewareAuthentication(handler.handleMyPosts)).Methods("GET")

	userRoute := apiRoute.PathPrefix("/users/{id}").Subrouter()
	userRoute.HandleFunc("", handler.handleUserGet).Methods("GET")
	userRoute.HandleFunc("/posts", handler.handleUserPostsGet).Methods("GET")

	postRoute := apiRoute.PathPrefix("/posts").Subrouter()
	postRoute.HandleFunc("/{id}", handler.handlePostGet).Methods("GET")
	postRoute.HandleFunc("", middlewares.SetMiddlewareAuthentication(handler.handlePostCreate)).Methods("POST")
	postRoute.HandleFunc("/{id}", middlewares.SetMiddlewareAuthentication(handler.handlePostUpdate)).Methods("PUT")
	postRoute.HandleFunc("/{id}", middlewares.SetMiddlewareAuthentication(handler.handlePostDelete)).Methods("DELETE")
	postRoute.HandleFunc("", handler.handlePostGetMany).Methods("GET")

}
