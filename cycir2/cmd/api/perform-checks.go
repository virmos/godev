package main

import (
	"cycir/internal/certificateutils"
	"cycir/internal/channeldata"
	"cycir/internal/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	// HTTP is the unencrypted web service check
	HTTP = 1
	// HTTPS is the encrypted web service check
	HTTPS = 2
	// SSLCertificate is ssl certificate check
	SSLCertificate = 3
)

// jsonResp describes the JSON response sent back to client
type jsonResp struct {
	OK            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	LastCheck     time.Time `json:"last_check"`
}

// ScheduledCheck performs a scheduled check on a host service by id
func (app *application) ScheduledCheck(hostServiceID int) {
	hs, err := app.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	h, err := app.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		return
	}

	// tests the service
	newStatus, msg := app.testServiceForHost(h, hs)

	if newStatus != hs.Status {
		app.updateHostServiceStatusCount(h, hs, newStatus, msg)
	}
}

func (app *application) updateHostServiceStatusCount(h models.Host, hs models.HostService, newStatus, msg string) {
	// update host service record in db with status and last check
	hs.Status = newStatus
	hs.LastMessage = msg
	hs.LastCheck = time.Now()
	err := app.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		return
	}

	pending, healthy, warning, problem, err := app.DB.GetAllServiceStatusCounts()
	if err != nil {
		log.Println(err)
		return
	}

	data := make(map[string]string)
	data["healthy_count"] = strconv.Itoa(healthy)
	data["pending_count"] = strconv.Itoa(pending)
	data["problem_count"] = strconv.Itoa(problem)
	data["warning_count"] = strconv.Itoa(warning)
	app.broadcastMessage("public-channel", "host-service-count-changed", data)
}

func (app *application) broadcastMessage(channel, messageType string, data map[string]string) {
	err := app.WsClient.Trigger(channel, messageType, data)
	if err != nil {
		log.Println(err)
	}
}

// TestCheck manually tests a host service and sends JSON response
func (app *application) TestCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	oldStatus := chi.URLParam(r, "oldStatus")
	okay := true

	// get host service
	hs, err := app.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		okay = false
	}

	// get host
	h, err := app.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		okay = false
	}

	// test the service
	newStatus, msg := app.testServiceForHost(h, hs)

	// save event
	event := models.Event{
		EventType:     newStatus,
		HostServiceID: hs.ID,
		HostID:        h.ID,
		ServiceName:   hs.Service.ServiceName,
		HostName:      hs.HostName,
		Message:       msg,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err = app.DB.InsertEvent(event)
	if err != nil {
		log.Println(err)
	}

	// broadcast service status changed event
	if newStatus != hs.Status {
		app.pushStatusChangedEvent(h, hs, newStatus)
	}

	// update the host service in the database with status (if changed) and last check
	hs.Status = newStatus
	hs.LastMessage = msg
	hs.LastCheck = time.Now()
	hs.UpdatedAt = time.Now()

	err = app.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		okay = false
	}

	app.pushScheduleChangedEvent(hs, newStatus)
	
	var resp jsonResp

	// create json
	if okay {
		resp = jsonResp{
			OK:            true,
			Message:       msg,
			ServiceID:     hs.ServiceID,
			HostServiceID: hs.ID,
			HostID:        hs.HostID,
			OldStatus:     oldStatus,
			NewStatus:     newStatus,
			LastCheck:     time.Now(),
		}
	} else {
		resp.OK = false
		resp.Message = "Something went wrong"
	}

	// send json to client
	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// testServiceForHost tests a service for a host
func (app *application) testServiceForHost(h models.Host, hs models.HostService) (string, string) {
	var msg, newStatus string

	switch hs.ServiceID {
	case HTTP:
		msg, newStatus = testHTTPForHost(h.URL)
		break

	case HTTPS:
		msg, newStatus = testHTTPSForHost(h.URL)
		break

	case SSLCertificate:
		msg, newStatus = testSSLForHost(h.URL)
		break
	}

	// broadcast to clients if appropriate
	if hs.Status != newStatus {
		app.pushStatusChangedEvent(h, hs, newStatus)

		// save event
		event := models.Event{
			EventType:     newStatus,
			HostServiceID: hs.ID,
			HostID:        h.ID,
			ServiceName:   hs.Service.ServiceName,
			HostName:      hs.HostName,
			Message:       msg,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		err := app.DB.InsertEvent(event)
		if err != nil {
			log.Println(err)
		}

		// send email if appropriate
		if app.PreferenceMap["notify_via_email"] == "1" {
			if hs.Status != "pending" {
				mm := channeldata.MailData{
					ToName:    app.PreferenceMap["notify_name"],
					ToAddress: app.PreferenceMap["notify_email"],
				}

				if newStatus == "healthy" {
					mm.Subject = fmt.Sprintf("HEALTHY: service %s on %s", hs.Service.ServiceName, hs.HostName)
					mm.Content = template.HTML(fmt.Sprintf(`<p>Service %s on %s apprted healthy status</p>
						<p><strong>Message received: %s</strong>/p>`, hs.Service.ServiceName, hs.HostName, msg))
				} else if newStatus == "problem" {
					mm.Subject = fmt.Sprintf("PROBLEM: service %s on %s", hs.Service.ServiceName, hs.HostName)
					mm.Content = template.HTML(fmt.Sprintf(`<p>Service %s on %s apprted problem</p>
						<p><strong>Message received: %s</strong></p>`, hs.Service.ServiceName, hs.HostName, msg))
				} else if newStatus == "warning" {

				}

				SendEmail(mm)
			}
		}

		// // send sms if appropriate
		// if app.PreferenceMap["notify_via_sms"] == "1" {
		// 	to := app.PreferenceMap["sms_notify_number"]
		// 	smsMessage := ""

		// 	if newStatus == "healthy" {
		// 		smsMessage = fmt.Sprintf("Service %s on %s is healthy", hs.Service.ServiceName, hs.HostName)
		// 	} else if newStatus == "problem" {
		// 		smsMessage = fmt.Sprintf("Service %s on %s apprts a problem: %s", hs.Service.ServiceName, hs.HostName, msg)
		// 	} else if newStatus == "warning" {
		// 		smsMessage = fmt.Sprintf("Service %s on %s apprts a warning: %s", hs.Service.ServiceName, hs.HostName, msg)
		// 	}

		// 	err := sms.SendTextTwilio(to, smsMessage, app)
		// 	if err != nil {
		// 		log.Println("Error sending sms in peform-checks.go", err)
		// 	}
		// }
	}

	app.pushScheduleChangedEvent(hs, newStatus)

	return newStatus, msg
}

func (app *application) pushStatusChangedEvent(h models.Host, hs models.HostService, newStatus string) {
	data := make(map[string]string)
	data["host_id"] = strconv.Itoa(hs.HostID)
	data["host_service_id"] = strconv.Itoa(hs.ID)
	data["host_name"] = h.HostName
	data["service_name"] = hs.Service.ServiceName
	data["icon"] = hs.Service.Icon
	data["status"] = newStatus
	data["message"] = fmt.Sprintf("%s on %s apprts %s", hs.Service.ServiceName, h.HostName, newStatus)
	data["last_check"] = time.Now().Format("2006-01-02 3:04:06 PM")

	app.broadcastMessage("public-channel", "host-service-status-changed", data)
}

func (app *application) pushScheduleChangedEvent(hs models.HostService, newStatus string) {
	// broadcast schedule-changed-event
	yearOne := time.Date(0001, 1, 1, 0, 0, 0, 1, time.UTC)
	data := make(map[string]string)

	data["host_service_id"] = strconv.Itoa(hs.ID)
	data["service_id"] = strconv.Itoa(hs.ServiceID)
	data["host_id"] = strconv.Itoa(hs.HostID)

	if app.Scheduler.Entry(app.MonitorMap[hs.ID]).Next.After(yearOne) {
		data["next_run"] = app.Scheduler.Entry(app.MonitorMap[hs.ID]).Next.Format("2006-01-02 3:04:05 PM")
	} else {
		data["next_run"] = "Pending..."
	}
	data["last_run"] = time.Now().Format("2006-01-02 3:04:05 PM")
	data["host"] = hs.HostName
	data["service"] = hs.Service.ServiceName
	data["schedule"] = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)
	data["status"] = newStatus
	data["host_name"] = hs.HostName
	data["icon"] = hs.Service.Icon

	app.broadcastMessage("public-channel", "schedule-changed-event", data)
}

// testHTTPForHost tests HTTP service
func testHTTPForHost(url string) (string, string) {
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}

	url = strings.Replace(url, "https://", "http://", -1)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}

// testHTTPSForHost tests HTTPS service
func testHTTPSForHost(url string) (string, string) {
	log.Println("Testing HTTPS")
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}

	url = strings.Replace(url, "http://", "https://", -1)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("HTTPS error 1")
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("HTTPS error 2", resp.StatusCode)
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}

// scanHost gets cert details from an internet host
func scanHost(hostname string, certDetailsChannel chan certificateutils.CertificateDetails, errorsChannel chan error) {
	res, err := certificateutils.GetCertificateDetails(hostname, 10)
	if err != nil {
		errorsChannel <- err
	} else {
		certDetailsChannel <- res
	}
}

// testSSLForHost tests an ssl certificate for a host
func testSSLForHost(url string) (string, string) {
	if strings.HasPrefix(url, "https://") {
		url = strings.Replace(url, "https://", "", -1)
	}
	if strings.HasPrefix(url, "http://") {
		url = strings.Replace(url, "http://", "", -1)
	}
	var certDetailsChannel chan certificateutils.CertificateDetails
	var errorsChannel chan error
	certDetailsChannel = make(chan certificateutils.CertificateDetails, 1)
	errorsChannel = make(chan error, 1)

	var msg string
	var newStatus string

	scanHost(url, certDetailsChannel, errorsChannel)

	for i, certDetailsInQueue := 0, len(certDetailsChannel); i < certDetailsInQueue; i++ {
		certDetails := <-certDetailsChannel
		certificateutils.CheckExpirationStatus(&certDetails, 30)

		if certDetails.Expired {
			// cert expired
			msg = certDetails.Hostname + " has expired!"

		} else if certDetails.ExpiringSoon {
			// cert expiring sono
			if certDetails.DaysUntilExpiration < 7 {
				msg = certDetails.Hostname + " expiring in " + strconv.Itoa(certDetails.DaysUntilExpiration) + " days"
				newStatus = "problem"
			} else {
				msg = certDetails.Hostname + " expiring in " + strconv.Itoa(certDetails.DaysUntilExpiration) + " days"
				newStatus = "warning"
			}
		} else {
			// cert okay
			msg = certDetails.Hostname + " expiring in " + strconv.Itoa(certDetails.DaysUntilExpiration) + " days"
			newStatus = "healthy"
		}
	}

	if len(errorsChannel) > 0 {
		fmt.Printf("There were %d error(s):\n", len(errorsChannel))
		for i, errorsInChannel := 0, len(errorsChannel); i < errorsInChannel; i++ {
			msg = fmt.Sprintf("%s\n", <-errorsChannel)
		}
		fmt.Printf("\n")
		newStatus = "problem"
	}

	return msg, newStatus
}

func (app *application) addToMonitorMap(hs models.HostService) {
	if app.PreferenceMap["monitoring_live"] == "1" {
		var j job
		j.HostServiceID = hs.ID
		scheduleID, err := app.Scheduler.AddJob(fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit), j)
		if err != nil {
			log.Println(err)
			return
		}

		app.MonitorMap[hs.ID] = scheduleID
		data := make(map[string]string)
		data["message"] = "scheduling"
		data["host_service_id"] = strconv.Itoa(hs.ID)
		data["next_run"] = "Pending..."
		data["service"] = hs.Service.ServiceName
		data["host"] = hs.HostName
		data["last_run"] = hs.LastCheck.Format("2006-01-02 3:04:05 PM")
		data["schedule"] = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)

		app.broadcastMessage("public-channel", "schedule-changed-event", data)
	}
}

func (app *application) removeFromMonitorMap(hs models.HostService) {
	if app.PreferenceMap["monitoring_live"] == "1" {
		app.Scheduler.Remove(app.MonitorMap[hs.ID])
		data := make(map[string]string)
		data["host_service_id"] = strconv.Itoa(hs.ID)
		app.broadcastMessage("public-channel", "schedule-item-removed-event", data)
	}
}