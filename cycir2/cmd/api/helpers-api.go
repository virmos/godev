package main

import (
	"cycir/internal/elastics"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// writeJSON writes aribtrary data out as JSON
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

// readJSON reads json from request body into data. We only accept a single json value in the body
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // max one megabyte in request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// we only allow one entry in the json file
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

// badRequest sends a JSON response with status http.StatusBadRequest, describing the error
func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = err.Error()

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(out)
	return nil
}

func (app *application) invalidCredentials(w http.ResponseWriter) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = "invalid authentication credentials"

	err := app.writeJSON(w, http.StatusUnauthorized, payload)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) inActiveAccount(w http.ResponseWriter) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = "inactive account"

	err := app.writeJSON(w, http.StatusUnauthorized, payload)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) passwordMatches(hash []byte, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func parseUptimeReports(uptimeReports, reports map[string]elastics.Report) (map[string]elastics.Report, error) {
	results :=  make(map[string]elastics.Report)

	for key, report := range reports {
		var result elastics.Report
		result.Histogram = make([]string, 24)
		result.Count = make([]string, 24)
		result.Host = key

		_, found := uptimeReports[key]
		for j := 0; j < 24; j++ { // number of hours in a day
			result.Count[j] = strconv.Itoa(report.HoursHistogram[j])
			if (found) {
				uptime := uptimeReports[key].HoursHistogram[j]	
				reportTime := report.HoursHistogram[j]	
				
				if (reportTime == 0) {
					result.Histogram[j] = "0%"
				} else {
					percent := int(float32(uptime) / float32(reportTime) * 100)
					result.Histogram[j] = strconv.Itoa(percent) + "%"
				}
			}
		}
		results[key] = result
	}
	return results, nil
}

func parseUptimeRangeReports(uptimeReports, reports map[string]elastics.Report) (map[string]elastics.Report, error) {
	results :=  make(map[string]elastics.Report)
	
	for key, report := range reports {
		var result elastics.Report
		result.Histogram = make([]string, 31)
		result.Count = make([]string, 31)
		result.Host = key

		_, found := uptimeReports[key]
		for j := 0; j < 31; j++ { // number of days in a month
			result.Count[j] = strconv.Itoa(report.DaysHistogram[j])
			result.Histogram[j] = "0%"
			if (found) {
				uptime := uptimeReports[key].DaysHistogram[j]	
				reportTime := report.DaysHistogram[j]	
				
				if (reportTime != 0) {
					percent := int(float32(uptime) / float32(reportTime) * 100)
					result.Histogram[j] = strconv.Itoa(percent) + "%"
				} 
			}
		}
		results[key] = result
	}
	return results, nil
}

// SetSystemPref sets a given system preference to supplied value, and returns JSON response
func (app *application) setSystemPref(prefName, prefValue string) error {
	err := app.repo.UpdateSystemPref(prefName, prefValue)
	if err != nil {
		log.Println(err)
		return err
	}

	app.PreferenceMap[prefName] = prefValue
	return nil
}

// ToggleMonitoring turns monitoring on and off
func (app *application) toggleMonitoring(enabled string) error {
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
			log.Println(err)
			return err
		}
	}
	return nil
}
