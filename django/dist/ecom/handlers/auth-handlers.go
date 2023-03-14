package handlers

import (
	"ecom/data"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	// "github.com/justinas/nosurf"
)

// UserRegister displays the register page
func (h *Handlers) UserRegister(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "register", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// PostUserRegister creates a new user
func (h *Handlers) PostUserRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var isAdmin bool
	if r.Form.Get("email") == "admin@example.com" {
		isAdmin = true
	}

	u := data.User{
		Name:     r.Form.Get("name"),
		Email:    r.Form.Get("email"),
		IsAdmin:  isAdmin,
		Password: r.Form.Get("password"),
	}

	id, err := h.Models.Users.Insert(u)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var userResp struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}
	userResp.ID = id
	userResp.Message = "User created successfully"
	_ = h.App.WriteJSON(w, http.StatusCreated, userResp)
}

// PostUserLogin attempts to log a user in
func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userFormInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		CSRF     string `json:"csrf_token"`
	}

	err := h.App.ReadJSON(w, r, &userFormInput)
	if err != nil {
		h.App.Error500(w, r)
		return
	}
	// if !nosurf.VerifyToken(nosurf.Token(r), userFormInput.CSRF) {
	// 	h.App.Error500(w, r)
	// 	return
	// }

	user, err := h.Models.Users.GetByEmail(userFormInput.Email)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	matches, err := user.PasswordMatches(userFormInput.Password)
	if err != nil {
		w.Write([]byte("Error validating password"))
		return
	}

	if !matches {
		w.Write([]byte("Invalid password!"))
		return
	}
	var userResp struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		IsAdmin bool   `json:"is_admin"`
	}

	userResp.Email = user.Email
	userResp.ID = user.ID
	userResp.Name = user.Name
	userResp.IsAdmin = user.IsAdmin
	_ = h.App.WriteJSON(w, http.StatusOK, userResp)
}

// GetAllUsers gets all users
func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var userFormInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		CSRF     string `json:"csrf_token"`
	}

	err := h.App.ReadJSON(w, r, &userFormInput)
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	// if !nosurf.VerifyToken(nosurf.Token(r), userFormInput.CSRF) {
	// 	h.App.Error500(w, r)
	// 	return
	// }

	users, err := h.Models.Users.GetAll()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var usersResp struct {
		Users []*data.User `json:"users"`
	}
	usersResp.Users = users
	_ = h.App.WriteJSON(w, http.StatusOK, usersResp)
}

// GetUserByEmail gets a user by email
func (h *Handlers) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	var userFormInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		CSRF     string `json:"csrf_token"`
	}

	err := h.App.ReadJSON(w, r, &userFormInput)
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	// if !nosurf.VerifyToken(nosurf.Token(r), userFormInput.CSRF) {
	// 	h.App.Error500(w, r)
	// 	return
	// }

	user, err := h.Models.Users.GetByEmail(userFormInput.Email)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var userResp struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		IsAdmin bool   `json:"is_admin"`
	}
	userResp.Email = user.Email
	userResp.ID = user.ID
	userResp.Name = user.Name
	userResp.IsAdmin = user.IsAdmin
	_ = h.App.WriteJSON(w, http.StatusOK, userResp)
}

// GetUserById gets a user by id
func (h *Handlers) GetUserById(w http.ResponseWriter, r *http.Request) {
	var userFormInput struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		CSRF     string `json:"csrf_token"`
	}

	err := h.App.ReadJSON(w, r, &userFormInput)
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	// if !nosurf.VerifyToken(nosurf.Token(r), userFormInput.CSRF) {
	// 	h.App.Error500(w, r)
	// 	return
	// }

	user, err := h.Models.Users.Get(userFormInput.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var userResp struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		IsAdmin bool   `json:"is_admin"`
	}
	userResp.Email = user.Email
	userResp.ID = user.ID
	userResp.Name = user.Name
	userResp.IsAdmin = user.IsAdmin
	_ = h.App.WriteJSON(w, http.StatusOK, userResp)
}

// UpdateUser updates a user
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userFormInput struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		CSRF     string `json:"csrf_token"`
	}

	err := h.App.ReadJSON(w, r, &userFormInput)
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	// if !nosurf.VerifyToken(nosurf.Token(r), userFormInput.CSRF) {
	// 	h.App.Error500(w, r)
	// 	return
	// }

	u, err := h.Models.Users.Get(userFormInput.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(userFormInput.Password), 12)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	
	u.Name = userFormInput.Name
	u.Email = userFormInput.Email
	u.Password = string(newHash)
	validator := h.App.Validator(nil)

	u.Validate(validator)

	if !validator.Valid() {
		w.Write([]byte("Failed validation!"))
		return
	}
	err = u.Update(*u)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var userResp struct {
		Message string `json:"message"`
	}
	userResp.Message = "User updated successfully"
	_ = h.App.WriteJSON(w, http.StatusOK, userResp)
}

// DeleteUser deletes a user by id
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userFormInput struct {
		ID       int    `json:"id"`
	}

	err := h.App.ReadJSON(w, r, &userFormInput)
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	// if !nosurf.VerifyToken(nosurf.Token(r), userFormInput.CSRF) {
	// 	h.App.Error500(w, r)
	// 	return
	// }

	u, err := h.Models.Users.Get(userFormInput.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = u.Delete(userFormInput.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var userResp struct {
		Message string `json:"message"`
	}
	userResp.Message = "User deleted successfully"
	_ = h.App.WriteJSON(w, http.StatusOK, userResp)
}