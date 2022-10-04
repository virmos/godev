package main

import (
	"cycir/internal/channeldata"
	"fmt"
	"strings"
	"html/template"
)

// ScheduledCheck performs a scheduled check on a host service by id
func (app *application) ScheduleReport() {
	// mm := channeldata.MailData{
	// 	ToName:    app.PreferenceMap["notify_name"],
	// 	ToAddress: app.PreferenceMap["notify_email"],
	// }
	// startDate := time.Date(0001, 11, 17, 20, 34, 58, 65138737, time.UTC)
	// endDate := time.Date(2023, 11, 17, 20, 34, 58, 65138737, time.UTC)

	// reports, _ := app.esrepo.GetRangeReport(app.config.esIndex, startDate.Format("2006-01-02T15:04:05Z07:00"), endDate.Format("2006-01-02T15:04:05Z07:00"))
	// uptimeReports, _ := app.esrepo.GetRangeUptimeReport(app.config.esIndex, startDate.Format("2006-01-02T15:04:05Z07:00"), endDate.Format("2006-01-02T15:04:05Z07:00"))
	// results, _ := parseUptimeRangeReports(uptimeReports, reports)
	// msgBuilder := ""
	// for key, report := range results {
	// 	msgBuilder = fmt.Sprintf(`Host %s 31 days uptime percentage report: `, key) + msgBuilder
	// 	msgBuilder = `<p>` + msgBuilder 
	// 	msgBuilder = msgBuilder + strings.Join(report.Histogram, ", ")
	// 	msgBuilder = msgBuilder + `</p>`  

	// 	msgBuilder = `<p>` + msgBuilder 
	// 	msgBuilder = msgBuilder + strings.Join(report.Count, ", ")
	// 	msgBuilder = msgBuilder + `</p>`  
	// }
	// mm.Subject = fmt.Sprintf("Range uptime report")
	// mm.Content = template.HTML(msgBuilder)
	// app.SendEmail(mm)

	mm := channeldata.MailData{
		ToName:    app.PreferenceMap["notify_name"],
		ToAddress: app.PreferenceMap["notify_email"],
	}

	reports, _ := app.esrepo.GetYesterdayReport(app.config.esIndex)
	uptimeReports, _ := app.esrepo.GetYesterdayUptimeReport(app.config.esIndex)
	results, _ := parseUptimeReports(uptimeReports, reports)
	
	msgBuilder := ""
	for key, report := range results {
		msgBuilder = fmt.Sprintf(`<h2> Yesterday %s 24 hours uptime percentage report: </h2>`, key) + msgBuilder
		msgBuilder = `<p>` + msgBuilder 
		msgBuilder = msgBuilder + strings.Join(report.Histogram, ", ")
		msgBuilder = msgBuilder + `</p>`  

		msgBuilder = `<p>` + msgBuilder 
		msgBuilder = msgBuilder + strings.Join(report.Count, ", ")
		msgBuilder = msgBuilder + `</p>`  
	}
	mm.Subject = "Range uptime report"
	mm.Content = template.HTML(msgBuilder)
	app.SendEmail(mm)
}