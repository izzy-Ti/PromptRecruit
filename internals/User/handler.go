package user

import (
	"net/http"
	"time"

	"github.com/izzy-Ti/PromptRecruit/internals/Utils"
)

type AuthHandler struct {
	svc *UserService
}

func NewAuthHandler(svc *UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	Utils.ParseJSON(r, &req)
	user, token, err := h.svc.LoginService(req.Email, req.Password)
	if err != nil {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})
	Utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user_id": user.ID,
	})
}
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	Utils.ParseJSON(r, &req)
	ok, err := h.svc.RegisterService(req.Email, req.Password, req.Name)
	if ok {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err,
		})
	}
	Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
		"success": false,
		"message": "Registration success",
	})
}
