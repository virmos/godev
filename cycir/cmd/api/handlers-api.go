package main

import (
	_ "bytes"
	"cycir/internal/channeldata"
	"cycir/internal/models"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	user, err := app.repo.GetUserByEmail(userInput.Email)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	// validate the password; send error if invalid password
	_, _, err = app.repo.Authenticate(userInput.Email, userInput.Password)
	if err == models.ErrInvalidCredentials {
		app.invalidCredentials(w)
		return
	} else if err == models.ErrInactiveAccount {
		app.inActiveAccount(w)
		return
	}

	// generate the token
	token, err := app.repo.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// save to database
	err = app.repo.InsertToken(token, user)
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
	user, err := app.repo.GetUserForToken(token)
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

	err := app.repo.DeleteUser(userID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	resp.Error = false
	resp.Redirect = true
	resp.Route = "/admin/users"
	resp.RedirectStatus = http.StatusSeeOther

	app.writeJSON(w, http.StatusOK, resp)
}

// PostSettings saves site settings
func (app *application) PostSettings(w http.ResponseWriter, r *http.Request) {
	var setting models.Setting

	err := app.readJSON(w, r, &setting)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	prefMap := make(map[string]string)

	prefMap["site_url"] = setting.SiteURL
	prefMap["notify_name"] = setting.NotifyName
	prefMap["notify_email"] = setting.NotifyEmail
	prefMap["smtp_server"] = setting.SMTP_Server
	prefMap["smtp_port"] = setting.SMTP_Port
	prefMap["smtp_user"] = setting.SMTP_User
	prefMap["smtp_password"] = setting.SMTP_Password
	prefMap["sms_enabled"] = setting.SMS_Enabled
	prefMap["sms_provider"] = setting.SMS_Provider
	prefMap["twilio_phone_number"] = setting.TWILIO_PhoneNumber
	prefMap["twilio_sid"] = setting.TWILIO_SID
	prefMap["twilio_auth_token"] = setting.TWILIO_AuthToken
	prefMap["smtp_from_email"] = setting.SMTP_FromEmail
	prefMap["smtp_from_name"] = setting.SMTP_FromName
	prefMap["notify_via_sms"] = setting.NotifyViaSMS
	prefMap["notify_via_email"] = setting.NotifyViaEmail
	prefMap["sms_notify_number"] = setting.SMS_NotifyNumber

	if setting.SMS_Enabled == "0" {
		prefMap["notify_via_sms"] = "0"
	}

	err = app.repo.InsertOrUpdateSitePreferences(prefMap)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// update app config
	for k, v := range prefMap {
		app.PreferenceMap[k] = v
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	resp.Error = false

	if setting.Action == "1" {
		resp.Redirect = true
		resp.Route = "/admin/overview"
		resp.RedirectStatus = http.StatusSeeOther
	} else {
		resp.Redirect = true
		resp.Route = "/admin/settings"
		resp.RedirectStatus = http.StatusSeeOther
	}

	app.writeJSON(w, http.StatusOK, resp)
}

// PostHost handles posting of host form
func (app *application) PostHost(w http.ResponseWriter, r *http.Request) {
	var newHost models.Host

	err := app.readJSON(w, r, &newHost)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var h models.Host

	if id > 0 {
		// get the host from the database
		host, err := app.repo.GetHostByID(id)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		h = host
	}

	h.HostName = newHost.HostName
	h.CanonicalName = newHost.CanonicalName
	h.URL = newHost.URL
	h.IP = newHost.IP
	h.IPV6 = newHost.IPV6
	h.Location = newHost.Location
	h.OS = newHost.OS
	h.Active = newHost.Active

	if id > 0 {
		err := app.repo.UpdateHost(h)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
	} else {
		newID, err := app.repo.InsertHost(h)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		h.ID = newID
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	resp.Error = false

	if newHost.Action == "1" {
		resp.Redirect = true
		resp.Route = "/admin/host/all"
		resp.RedirectStatus = http.StatusSeeOther
	} else {
		resp.Redirect = true
		resp.Route = fmt.Sprintf("/admin/host/%d", h.ID)
		resp.RedirectStatus = http.StatusSeeOther
	}

	app.writeJSON(w, http.StatusOK, resp)
}

// PostHost handles posting of host form
func (app *application) DeleteHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	hostID, _ := strconv.Atoi(id)

	err := app.repo.DeleteHost(hostID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	resp.Error = false
	resp.Message = "Host deleted"
	app.pushHostRemovedEvent(id)
 	app.writeJSON(w, http.StatusOK, resp)
}

// PostOneUser adds/edits a user
func (app *application) PostOneUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorLog.Println(err)
	}

	var newUser models.User
	err = app.readJSON(w, r, &newUser)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var u models.User

	if id > 0 {
		u, _ = app.repo.GetUserById(id)
		u.FirstName = newUser.FirstName
		u.LastName = newUser.LastName
		u.Email = newUser.Email
		u.UserActive = newUser.UserActive

		err := app.repo.UpdateUser(u)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		if len(newUser.Password) > 0 {
			// changing password
			err := app.repo.UpdatePassword(id, string(newUser.Password))
			if err != nil {
				app.errorLog.Println(err)
				return
			}
		}
	} else {
		u.FirstName = newUser.FirstName
		u.LastName = newUser.LastName
		u.Email = newUser.Email
		u.UserActive = newUser.UserActive
		u.Password = newUser.Password
		u.AccessLevel = 3

		_, err := app.repo.InsertUser(u)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	resp.Error = false
	resp.Redirect = true
	resp.Route = "/admin/users"
	resp.RedirectStatus = http.StatusSeeOther

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) SendRangeUptimeReport(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		HostName  string `json:"host_name"`
		CSRF      string `json:"csrf_token"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	mm := channeldata.MailData{
		ToName:    app.PreferenceMap["notify_name"],
		ToAddress: app.PreferenceMap["notify_email"],
	}
	startDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	endDate, err := time.Parse("2006-01-02", payload.EndDate)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	numDays := int(endDate.Sub(startDate).Hours()/24) + 1
	var endReportedDate string
	if numDays < 31 {
		endReportedDate = startDate.AddDate(0, 0, numDays).Format("2006-01-02")
	}
	if numDays == 31 {
		endReportedDate = startDate.AddDate(0, 0, 31).Format("2006-01-02")
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		HostName       string `json:"host_name"`
		Histogram      string `json:"histogram"`
		Count          string `json:"count"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	if numDays <= 31 {
		if numDays == 1 { 
			resp.Error = true
			resp.Message = "Please specify range of days"
		} else {
			reports, _ := app.esrepo.GetRangeReport(app.config.esIndex, payload.HostName, startDate.Format("2006-01-02T15:04:05Z07:00"), endDate.Format("2006-01-02T15:04:05Z07:00"))
			uptimeReports, _ := app.esrepo.GetRangeUptimeReport(app.config.esIndex, payload.HostName, startDate.Format("2006-01-02T15:04:05Z07:00"), endDate.Format("2006-01-02T15:04:05Z07:00"))
			results, _ := parseUptimeRangeReports(uptimeReports, reports)
	
			if len(results) == 0 {
				resp.Error = true
				resp.Message = "No report was made between specified dates"
			} else {
				msgBuilder := ""
				for key, report := range results {
					histogramString := strings.Join(report.Histogram, ", ")
					countString := strings.Join(report.Count, ", ")
					resp.Histogram = histogramString
					resp.Count = countString
					resp.HostName = key
	
					msgBuilder = msgBuilder + fmt.Sprintf(`<h2> Host %s %s days uptime percentage report. </h2>`, key, strconv.Itoa(numDays))
					msgBuilder = msgBuilder + fmt.Sprintf(`<p> Data between to dates: %s to %s are reported: </p>`, payload.StartDate, payload.EndDate)
					msgBuilder = msgBuilder + fmt.Sprintf(`<p> Start Day: %s ----> Percentage: `, payload.StartDate)
					msgBuilder = msgBuilder + histogramString
					msgBuilder = msgBuilder + fmt.Sprintf(` <---- End Day: %s</p>`, endReportedDate)
	
					msgBuilder = msgBuilder + fmt.Sprintf(`<p> Start Day: %s ----> Total reqs: `, payload.StartDate)
					msgBuilder = msgBuilder + countString
					msgBuilder = msgBuilder + fmt.Sprintf(` <---- End Day: %s</p>`, endReportedDate)
				}
				mm.Subject = "Range uptime report"
				mm.Content = template.HTML(msgBuilder)
				app.SendEmail(mm)
	
				resp.Error = false
				resp.Message = "Sent report to mail"
			}
		}
	} else {
		resp.Error = true
		resp.Message = "Date range too big (>31 days)"
	}
	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) SendRangeUptimeReportCached(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		HostName  string `json:"host_name"`
		Histogram string `json:"histogram"`
		Count     string `json:"count"`
		CSRF      string `json:"csrf_token"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	mm := channeldata.MailData{
		ToName:    app.PreferenceMap["notify_name"],
		ToAddress: app.PreferenceMap["notify_email"],
	}
	startDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	endDate, err := time.Parse("2006-01-02", payload.EndDate)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	numDays := int(endDate.Sub(startDate).Hours()/24) + 1

	var endReportedDate string
	if numDays < 31 {
		endReportedDate = startDate.AddDate(0, 0, numDays).Format("2006-01-02")
	}
	if numDays == 31 {
		endReportedDate = startDate.AddDate(0, 0, 31).Format("2006-01-02")
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	if numDays <= 31 {
		msgBuilder := ""
		histogramString := payload.Histogram
		countString := payload.Count

		msgBuilder = msgBuilder + fmt.Sprintf(`<h2> Host %s %s days uptime percentage report. </h2>`, payload.HostName, strconv.Itoa(numDays))
		msgBuilder = msgBuilder + fmt.Sprintf(`<p> Start Day: %s ----> Percentage: `, payload.StartDate)
		msgBuilder = msgBuilder + histogramString
		msgBuilder = msgBuilder + fmt.Sprintf(` <---- End Day: %s</p>`, endReportedDate)

		msgBuilder = msgBuilder + fmt.Sprintf(`<p> Start Day: %s ----> Total reqs: `, payload.StartDate)
		msgBuilder = msgBuilder + countString
		msgBuilder = msgBuilder + fmt.Sprintf(` <---- End Day: %s</p>`, endReportedDate)

		mm.Subject = "Range uptime report"
		mm.Content = template.HTML(msgBuilder)
		app.SendEmail(mm)

		resp.Error = false
		resp.Message = "Sent report to mail"
	} else {
		resp.Error = true
		resp.Message = "Date range too big (>31 days)"
	}

	resp.Redirect = false
	resp.Route = "/admin/host/all"
	resp.RedirectStatus = http.StatusSeeOther

	app.writeJSON(w, http.StatusOK, resp)
}

// ToggleServiceForHost turns a host service on or off (active or inactive)
func (app *application) ToggleServiceForHost(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		HostID    int `json:"host_id"`
		ServiceID int `json:"service_id"`
		Active    int `json:"active"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	hostID := payload.HostID
	serviceID := payload.ServiceID
	active := payload.Active

	err = app.repo.UpdateHostServiceStatus(hostID, serviceID, active)
	if err != nil {
		app.errorLog.Println(err)
	}

	// broadcast
	hs, _ := app.repo.GetHostServiceByHostIDServiceID(hostID, serviceID)
	h, _ := app.repo.GetHostByID(hostID)

	// add or remove from schedule
	if active == 1 {
		app.pushScheduleChangedEvent(hs, "pending")
		app.pushStatusChangedEvent(h, hs, "pending")
		app.addToMonitorMap(hs)
	} else {
		app.removeFromMonitorMap(hs)
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	resp.Error = false
	resp.Redirect = false

	app.writeJSON(w, http.StatusOK, resp)
}

// SetSystemPref sets a given system preference to supplied value, and returns JSON response
func (app *application) SetSystemPref(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}

	var payload struct {
		PrefName  string `json:"pref_name"`
		PrefValue string `json:"pref_value"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	prefName := payload.PrefName
	prefValue := payload.PrefValue

	err = app.repo.UpdateSystemPref(prefName, prefValue)
	if err != nil {
		resp.Error = true
		resp.Message = err.Error()
	}

	app.PreferenceMap["monitoring_live"] = prefValue

	resp.Error = false
	resp.Redirect = false

	app.writeJSON(w, http.StatusOK, resp)
}

// ToggleMonitoring turns monitoring on and off
func (app *application) ToggleMonitoring(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Enabled string `json:"enabled"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	enabled := payload.Enabled

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

		// empty the monitor map
		for k := range app.MonitorMap {
			delete(app.MonitorMap, k)
		}

		// empty the schedule map
		for k := range app.MonitorMap {
			delete(app.FunctionMap, k)
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
			app.errorLog.Println(err)
		}
	}

	var resp struct {
		Error          bool   `json:"error"`
		Message        string `json:"message"`
		Redirect       bool   `json:"redirect"`
		RedirectStatus int    `json:"status"`
		Route          string `json:"route"`
	}
	resp.Error = false
	resp.Redirect = false

	app.writeJSON(w, http.StatusOK, resp)
}
