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
		msgBuilder = fmt.Sprintf(`<h2> Yesterday of %s uptime percentage report in 24 hours: </h2>`, key) + msgBuilder
		msgBuilder = `<p>` + msgBuilder 
		msgBuilder = msgBuilder + strings.Join(report.Histogram, ", ")
		msgBuilder = msgBuilder + `</p>`  

		msgBuilder = `<p>` + msgBuilder 
		msgBuilder = msgBuilder + strings.Join(report.Count, ", ")
		msgBuilder = msgBuilder + `</p>`  
	}
	mm.Subject = "Yesterday uptime report"
	mm.Content = template.HTML(msgBuilder)
	app.SendEmail(mm)
}