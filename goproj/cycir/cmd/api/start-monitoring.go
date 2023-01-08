package main

import (
	"fmt"
	"strconv"
	"time"
)

// job is the unit of work to be performed

var a *application

type job struct {
	HostServiceID int
}

func NewScheduler(app *application) {
	a = app
}

// Run runs the scheduled job
func (j job) Run() {
	a.ScheduledCheck(j.HostServiceID)
}

// StartMonitoring starts the monitoring process
func (app *application) StartMonitoring() {
	if app.PreferenceMap["monitoring_live"] == "1" {
		// trigger a message to broadcast to all clients that app is starting to monitor
		data := make(map[string]string)
		data["message"] = "Monitoring is starting..."
		err := app.WsClient.Trigger("public-channel", "app-starting", data)
		if err != nil {
			app.errorLog.Println(err)
		}

		// get all of the services that we want to monitor
		servicesToMonitor, err := app.repo.GetServicesToMonitor()
		if err != nil {
			app.errorLog.Println(err)
		}

		// range through the services
		for _, x := range servicesToMonitor {
			// get the schedule unit and number
			var sch string
			if x.ScheduleUnit == "d" {
				sch = fmt.Sprintf("@every %d%s", x.ScheduleNumber*24, "h")
			} else {
				// sch = fmt.Sprintf("@every %d%s", x.ScheduleNumber, x.ScheduleUnit) // @every 3m, scheduled check
				sch = fmt.Sprintf("@every %d%s", 10, "s")
			}

			// create a job
			var j job
			j.HostServiceID = x.ID
			scheduleID, err := app.Scheduler.AddJob(sch, j)
			if err != nil {
				app.errorLog.Println(err)
			}
			// save the id of the job so we can start/stop it
			app.MonitorMap[x.ID] = scheduleID

			// create a schedule
			funcID, err := app.Scheduler.AddFunc("0 5 * * ?", app.ScheduleReport) // 5 am, send yesterday uptime report
			// funcID, err := app.Scheduler.AddFunc("@every 0h0m3s", app.ScheduleReport)
			if err != nil {
				app.errorLog.Println(err)
			}
			// save the id of the function so we can start/stop it
			app.FunctionMap[x.ID] = funcID

			// broadcast over websockets the fact that the service is scheduled
			payload := make(map[string]string)
			payload["message"] = "scheduling"
			payload["host_service_id"] = strconv.Itoa(x.ID)
			yearOne := time.Date(0001, 11, 17, 20, 34, 58, 65138737, time.UTC)
			if app.Scheduler.Entry(app.MonitorMap[x.ID]).Next.After(yearOne) {
				payload["next_run"] = app.Scheduler.Entry(app.MonitorMap[x.ID]).Next.Format("2006-01-02 3:04:05 PM")
			} else {
				payload["next_run"] = "Pending...."
			}
			payload["host"] = x.HostName
			payload["service"] = x.Service.ServiceName
			if x.LastCheck.After(yearOne) {
				payload["last_run"] = x.LastCheck.Format("2006-01-02 3:04:05 PM")
			} else {
				payload["last_run"] = "Pending..."
			}
			payload["schedule"] = fmt.Sprintf("@every %d%s", x.ScheduleNumber, x.ScheduleUnit)

			err = app.WsClient.Trigger("public-channel", "next-run-event", payload)
			if err != nil {
				app.errorLog.Println(err)
			}

			err = app.WsClient.Trigger("public-channel", "schedule-changed-event", payload)
			if err != nil {
				app.errorLog.Println(err)
			}

		}
	}
}
