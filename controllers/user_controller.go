package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"sheinko.tk/copy_project/repository"
	"sheinko.tk/copy_project/utils/responses"
)

func (handler Handler) handleUserGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := repository.NewUserRepository(handler.DB)

	user, err := db.FindById(uint(i))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func (handler Handler) handleUserPostsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	posts, err := db.FindPostsByUserId(uint(i))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		log.Println(err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}
