package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Drinnn/go-expert-api/internal/dto"
	"github.com/Drinnn/go-expert-api/internal/entity"
	"github.com/Drinnn/go-expert-api/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB            database.UserInterface
	Jwt               *jwtauth.JWTAuth
	JwtExpirationTime int
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpirationTime int) *UserHandler {
	return &UserHandler{
		UserDB:            db,
		Jwt:               jwt,
		JwtExpirationTime: jwtExpirationTime,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, req *http.Request) {
	var inputDto dto.LoginInput
	err := json.NewDecoder(req.Body).Decode(&inputDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserDB.FindByEmail(inputDto.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !user.ValidatePassword(inputDto.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := h.Jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * time.Duration(h.JwtExpirationTime)).Unix(),
	})
	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	var userInput dto.CreateUserInput
	err := json.NewDecoder(req.Body).Decode(&userInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := entity.NewUser(userInput.Name, userInput.Email, userInput.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
