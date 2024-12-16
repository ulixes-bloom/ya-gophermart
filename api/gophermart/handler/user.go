package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
	"github.com/ulixes-bloom/ya-gophermart/internal/security"
)

//	@Summary	Регистрация пользователя
//	@ID			RegisterUser
//	@Produce	json
//	@Success	200	"пользователь успешно зарегистрирован и аутентифицирован"
//	@Failure	400	"неверный формат запроса"
//	@Failure	409	"логин уже занят"
//	@Failure	500	"внутренняя ошибка сервера"
//	@Router		/api/user/register [post]
//	@Param		user	body	models.User	true	"User Registration Information"
func (h *HTTPHandler) RegisterUser(rw http.ResponseWriter, req *http.Request) {
	user := &models.User{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(user); err != nil {
		h.Error(rw, err, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Login == "" || user.Password == "" {
		h.Error(rw,
			appErrors.ErrUserLoginAndPasswordRequired,
			appErrors.ErrUserLoginAndPasswordRequired.Error(),
			http.StatusBadRequest)
		return
	}

	createdUserID, err := h.app.RegisterUser(user)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserLoginAlreadyExists) {
			h.Error(rw, err, appErrors.ErrUserLoginAlreadyExists.Error(), http.StatusConflict)
			return
		}
		h.Error(rw, err, err.Error(), http.StatusInternalServerError)
		return
	}

	jwtToken, err := security.BuildJWTString(createdUserID, h.conf.TokenSecretKey, h.conf.TokenLifetime)
	if err != nil {
		h.Error(rw, err, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	rw.WriteHeader(http.StatusOK)
}

//	@Summary	Аутентификация пользователя
//	@ID			AuthUser
//	@Produce	json
//	@Success	200	"пользователь успешно аутентифицирован"
//	@Failure	400	"неверный формат запроса"
//	@Failure	401	"неверная пара логин/пароль"
//	@Failure	500	"внутренняя ошибка сервера"
//	@Router		/api/user/login [post]
//	@Param		user	body	models.User	true	"User Registration Information"
func (h *HTTPHandler) AuthUser(rw http.ResponseWriter, req *http.Request) {
	user := &models.User{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(user); err != nil {
		h.Error(rw, err, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Login == "" || user.Password == "" {
		h.Error(rw,
			appErrors.ErrUserLoginAndPasswordRequired,
			appErrors.ErrUserLoginAndPasswordRequired.Error(),
			http.StatusBadRequest)
		return
	}

	dbUser, err := h.app.ValidateUser(user)
	if err != nil {
		h.Error(rw, err, appErrors.ErrInvalidUserLoginOrPassword.Error(), http.StatusUnauthorized)
		return
	}

	jwtToken, err := security.BuildJWTString(dbUser.ID, h.conf.TokenSecretKey, h.conf.TokenLifetime)
	if err != nil {
		h.Error(rw, err, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	rw.WriteHeader(http.StatusOK)
}
