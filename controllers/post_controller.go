package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"sheinko.tk/copy_project/models"
	"sheinko.tk/copy_project/repository"
	"sheinko.tk/copy_project/utils/auth"
	"sheinko.tk/copy_project/utils/responses"
)

func (handler Handler) handleMyPosts(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db := repository.NewPostRepository(handler.DB)
	posts, err := db.FindMyPosts(uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func (handler *Handler) handlePostCreate(w http.ResponseWriter, r *http.Request) {
	var postDTO models.PostDTO
	if err := json.NewDecoder(r.Body).Decode(&postDTO); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.DTOToPost(postDTO)

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	post.AuthorID = &uid

	db := repository.NewPostRepository(handler.DB)

	if err := db.Save(&post); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func (handler *Handler) handlePostGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	post, err := db.FindById(uint(i))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id "+id+" could not found"))
		} else {
			responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
			log.Println(err)
		}
		return
	}

	if post.IsPublished == false {
		uid, err := auth.ExtractTokenID(r)
		if err != nil {
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id "+id+" could not found"))
			return
		}

		if uid != post.Author.ID {
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id "+id+" could not found"))
			return
		}
	}

	responses.JSON(w, http.StatusOK, post)
}

func (handler Handler) handlePostUpdate(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	post, err := db.FindById(uint(pid))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.ERROR(w, http.StatusNotFound, errors.New("the post with id "+vars["id"]+" could not found"))
		} else {
			responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
			log.Println(err)
		}
		return
	}

	if post.Author.ID != uid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("you can not update the post who belongs to someone else"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var postUpdate models.PostDTO

	if err = json.Unmarshal(body, &postUpdate); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	newPost := models.DTOToPost(postUpdate)

	if err = db.UpdateById(&post, newPost); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func (handler *Handler) handlePostDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := repository.NewPostRepository(handler.DB)

	post, err := db.FindById(uint(i))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("the post with id "+id+" could not found"))
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if post.Author.ID != uid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("you can not delete the post who belongs to someone else"))
		return
	}

	if err := db.DeleteById(uint(i)); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("something went wrong"))
		log.Println(err)
		return
	}

	responses.JSON(w, http.StatusNoContent, "")
}

func (handler Handler) handlePostGetMany(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	limitStr := keys.Get("limit")

	var limit int
	var err error
	if len(limitStr) != 0 {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
	}

	db := repository.NewPostRepository(handler.DB)

	posts, err := db.FindMany(limit)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}
