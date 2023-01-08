package main

import (
	"net/http"
	"encoding/json"
)

type serviceJSON struct {
	Error   bool   `json:"error"`
	Message string `json:"msg"`
	Value   string `json:"value"`
}

func (app *application) SaveInCache(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}

	key := r.Form.Get("key")
	value := r.Form.Get("value")
	_ = r.Form.Get("csrf_token")

	// if !nosurf.VerifyToken(nosurf.Token(r), userInput.CSRF) {
	// 	Error500(w, r)
	// 	return
	// }

	var resp serviceJSON
	resp.Error = false

	err = app.Cache.Set(key, value)
	if err != nil {
		resp.Error = true
		ClientError(w, r, http.StatusInternalServerError)
		return
	}

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) GetFromCache(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}

	key := r.Form.Get("key")
	_ = r.Form.Get("csrf_token")

	var resp serviceJSON
	resp.Error = false

	fromCache, err := app.Cache.Get(key)
	inCache := true

	if err != nil {
		resp.Message = "Not found in cache!"
		inCache = false
	}

	if inCache {
		resp.Error = false
		resp.Message = "Success"
		resp.Value = fromCache.(string)
	} else {
		resp.Error = true
	}

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) DeleteFromCache(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}

	key := r.Form.Get("key")
	_ = r.Form.Get("csrf_token")

	var resp serviceJSON
	resp.Error = false	

	err = app.Cache.Forget(key)
	if err != nil {
		ClientError(w, r, http.StatusInternalServerError)
		return
	}

	resp.Error = false
	resp.Message = "Deleted from cache (if it existed)"

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) EmptyCache(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}

	_ = r.Form.Get("csrf_token")

	var resp serviceJSON
	resp.Error = false	

	err = app.Cache.Empty()
	if err != nil {
		ClientError(w, r, http.StatusInternalServerError)
		return
	}

	resp.Error = false
	resp.Message = "Emptied cache"

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
