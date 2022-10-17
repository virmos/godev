package main

import (
	"cycir/internal/channeldata"
	"fmt"
	"strings"
	"html/template"
)

// ScheduledCheck performs a scheduled check on a host service by id
func (app *application) ScheduleReport() {
	mm := channeldata.MailData{
		ToName:    app.PreferenceMap["notify_name"],
		ToAddress: app.PreferenceMap["notify_email"],
	}

	reports, _ := app.esrepo.GetYesterdayReport(app.config.esIndex)
	uptimeReports, _ := app.esrepo.GetYesterdayUptimeReport(app.config.esIndex)
	results, _ := parseUptimeReports(uptimeReports, reports)
	
	msgBuilder := ""
	for key, report := range results {
		msgBuilder = msgBuilder + fmt.Sprintf(`<div style="font-size: 2rem;"> Yesterday of %s uptime percentage report in 24 hours: </div>`, key)
		msgBuilder = msgBuilder + `<p> First Hour ----> Percentage: `
		msgBuilder = msgBuilder + strings.Join(report.Histogram, ", ")
		msgBuilder = msgBuilder + ` <---- Final Hour </p>`  

		msgBuilder = msgBuilder + `<p> First Hour ----> Total reqs: ` 
		msgBuilder = msgBuilder + strings.Join(report.Count, ", ")
		msgBuilder = msgBuilder + `<---- Final Hour </p>`  
	}
	mm.Subject = "Yesterday uptime report"
	mm.Content = template.HTML(msgBuilder)
	app.SendEmail(mm)
}