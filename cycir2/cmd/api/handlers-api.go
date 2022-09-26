package main

import (
	_ "bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"cycir/internal/models"
	"log"

	"github.com/go-chi/chi/v5"
	_ "golang.org/x/crypto/bcrypt"
)

// CreateAuthToken creates and sends an auth token, if user supplies valid information
func (app *application) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// get the user from the database by email; send error if invalid email
	user, err := app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	// validate the password; send error if invalid password
	_, _, err = app.DB.Authenticate(userInput.Email, userInput.Password)
	if err == models.ErrInvalidCredentials {
		app.invalidCredentials(w)
		return
	} else if err == models.ErrInactiveAccount {
		app.inActiveAccount(w)
		return
	}

	// generate the token
	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	log.Println(token)

	// save to database
	err = app.DB.InsertToken(token, user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// send response

	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Token   *models.Token `json:"authentication_token"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("token for %s created", userInput.Email)
	payload.Token = token

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// authenticateToken checks an auth token for validity
func (app *application) authenticateToken(r *http.Request) (*models.User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authentication token wrong size")
	}

	// get the user from the tokens table
	user, err := app.DB.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

// CheckAuthentication checks auth status
func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	// validate the token, and get associated user
	user, err := app.authenticateToken(r)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	// valid user
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("authenticated user %s", user.Email)
	app.writeJSON(w, http.StatusOK, payload)
}

// DeleteUser deletes a user, and all associated tokens, from the database
func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	err := app.DB.DeleteUser(userID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	app.writeJSON(w, http.StatusOK, resp)
}

// PostSettings saves site settings
func (app *application) PostSettings(w http.ResponseWriter, r *http.Request) {
	prefMap := make(map[string]string)

	prefMap["site_url"] = r.Form.Get("site_url")
	prefMap["notify_name"] = r.Form.Get("notify_name")
	prefMap["notify_email"] = r.Form.Get("notify_email")
	prefMap["smtp_server"] = r.Form.Get("smtp_server")
	prefMap["smtp_port"] = r.Form.Get("smtp_port")
	prefMap["smtp_user"] = r.Form.Get("smtp_user")
	prefMap["smtp_password"] = r.Form.Get("smtp_password")
	prefMap["sms_enabled"] = r.Form.Get("sms_enabled")
	prefMap["sms_provider"] = r.Form.Get("sms_provider")
	prefMap["twilio_phone_number"] = r.Form.Get("twilio_phone_number")
	prefMap["twilio_sid"] = r.Form.Get("twilio_sid")
	prefMap["twilio_auth_token"] = r.Form.Get("twilio_auth_token")
	prefMap["smtp_from_email"] = r.Form.Get("smtp_from_email")
	prefMap["smtp_from_name"] = r.Form.Get("smtp_from_name")
	prefMap["notify_via_sms"] = r.Form.Get("notify_via_sms")
	prefMap["notify_via_email"] = r.Form.Get("notify_via_email")
	prefMap["sms_notify_number"] = r.Form.Get("sms_notify_number")

	if r.Form.Get("sms_enabled") == "0" {
		prefMap["notify_via_sms"] = "0"
	}

	err := app.DB.InsertOrUpdateSitePreferences(prefMap)
	if err != nil {
		log.Println(err)
		return
	}

	// update app config
	for k, v := range prefMap {
		app.PreferenceMap[k] = v
	}

	if r.Form.Get("action") == "1" {
		http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
	}
}

// PostHost handles posting of host form
func (app *application) PostHost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var h models.Host

	if id > 0 {
		// get the host from the database
		host, err := app.DB.GetHostByID(id)
		if err != nil {
			log.Println(err)
			return
		}
		h = host
	}

	h.HostName = r.Form.Get("host_name")
	h.CanonicalName = r.Form.Get("canonical_name")
	h.URL = r.Form.Get("url")
	h.IP = r.Form.Get("ip")
	h.IPV6 = r.Form.Get("ipv6")
	h.Location = r.Form.Get("location")
	h.OS = r.Form.Get("os")
	active, _ := strconv.Atoi(r.Form.Get("active"))
	h.Active = active

	if id > 0 {
		err := app.DB.UpdateHost(h)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		newID, err := app.DB.InsertHost(h)
		if err != nil {
			log.Println(err)
			return
		}
		h.ID = newID
	}

	if r.Form.Get("action") == "1" {
		http.Redirect(w, r, "/admin/host/all", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/host/%d", h.ID), http.StatusSeeOther)
	}
}

// PostOneUser adds/edits a user
func (app *application) PostOneUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	var u models.User

	if id > 0 {
		u, _ = app.DB.GetUserById(id)
		u.FirstName = r.Form.Get("first_name")
		u.LastName = r.Form.Get("last_name")
		u.Email = r.Form.Get("email")
		u.UserActive, _ = strconv.Atoi(r.Form.Get("user_active"))
		err := app.DB.UpdateUser(u)
		if err != nil {
			log.Println(err)
			return
		}

		if len(r.Form.Get("password")) > 0 {
			// changing password
			err := app.DB.UpdatePassword(id, r.Form.Get("password"))
			if err != nil {
				log.Println(err)
				return
			}
		}
	} else {
		u.FirstName = r.Form.Get("first_name")
		u.LastName = r.Form.Get("last_name")
		u.Email = r.Form.Get("email")
		u.UserActive, _ = strconv.Atoi(r.Form.Get("user_active"))
		u.Password = []byte(r.Form.Get("password"))
		u.AccessLevel = 3

		_, err := app.DB.InsertUser(u)
		if err != nil {
			log.Println(err)
			return
		}
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

type serviceJSON struct {
	OK bool `json:"ok"`
}

// ToggleServiceForHost turns a host service on or off (active or inactive)
func (app *application) ToggleServiceForHost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	var resp serviceJSON
	resp.OK = true

	hostID, _ := strconv.Atoi(r.Form.Get("host_id"))
	serviceID, _ := strconv.Atoi(r.Form.Get("service_id"))
	active, _ := strconv.Atoi(r.Form.Get("active"))

	err = app.DB.UpdateHostServiceStatus(hostID, serviceID, active)
	if err != nil {
		log.Println(err)
		resp.OK = false
	}

	// broadcast
	hs, _ := app.DB.GetHostServiceByHostIDServiceID(hostID, serviceID)
	h, _ := app.DB.GetHostByID(hostID)

	// add or remove from schedule
	if active == 1 {
		app.pushScheduleChangedEvent(hs, "pending")
		app.pushStatusChangedEvent(h, hs, "pending")
		app.addToMonitorMap(hs)
	} else {
		app.removeFromMonitorMap(hs)
	}

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// SetSystemPref sets a given system preference to supplied value, and returns JSON response
func (app *application) SetSystemPref(w http.ResponseWriter, r *http.Request) {
	prefName := r.PostForm.Get("pref_name")
	prefValue := r.PostForm.Get("pref_value")

	var resp jsonResp
	resp.OK = true
	resp.Message = ""

	err := app.DB.UpdateSystemPref(prefName, prefValue)
	if err != nil {
		resp.OK = false
		resp.Message = err.Error()
	}

	app.PreferenceMap["monitoring_live"] = prefValue

	out, _ := json.MarshalIndent(resp, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// ToggleMonitoring turns monitoring on and off
func (app *application) ToggleMonitoring(w http.ResponseWriter, r *http.Request) {
	enabled := r.PostForm.Get("enabled")

	if enabled == "1" {
		// start monitoring
		app.PreferenceMap["monitoring_live"] = "1"
		app.StartMonitoring()
		app.Scheduler.Start()
	} else {
		// stop monitoring
		app.PreferenceMap["monitoring_live"] = "0"

		// remove all items in map from schedule
		for _, x := range app.MonitorMap {
			app.Scheduler.Remove(x)
		}

		// empty the map
		for k := range app.MonitorMap {
			delete(app.MonitorMap, k)
		}

		// delete all entries from schedule, to be sure
		for _, i := range app.Scheduler.Entries() {
			app.Scheduler.Remove(i.ID)
		}

		app.Scheduler.Stop()

		data := make(map[string]string)
		data["message"] = "Monitoring is off!"
		err := app.WsClient.Trigger("public-channel", "app-stopping", data)
		if err != nil {
			log.Println(err)
		}

	}

	var resp jsonResp
	resp.OK = true

	out, _ := json.MarshalIndent(resp, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
