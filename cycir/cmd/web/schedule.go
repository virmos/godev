package main

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"cycir/internal/models"
	"log"
	"net/http"
	"sort"
)

// ByHost allows us to sort by host
type ByHost []models.Schedule

// Len is used to sort by host
func (a ByHost) Len() int { return len(a) }

// Less is used to sort by host
func (a ByHost) Less(i, j int) bool { return a[i].Host < a[j].Host }

// Swap is used to sort by host
func (a ByHost) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// ListEntries lists schedule entries
func (app *application) ListEntries(w http.ResponseWriter, r *http.Request) {
	var items []models.Schedule
	servicesToMonitor, err := app.repo.GetServicesToMonitor()
	if err != nil {
		log.Println(err)
		return
	}	

	for _, x := range servicesToMonitor {
		var item models.Schedule
		hs, err := app.repo.GetHostServiceByID(x.ID)
		if err != nil {
			log.Println(err)
			return
		}
		item.ScheduleText = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)
		item.LastRunFromHS = hs.LastCheck
		item.Host = hs.HostName
		item.Service = hs.Service.ServiceName
		items = append(items, item)
	}

	// sort the slice
	sort.Sort(ByHost(items))

	data := make(jet.VarMap)
	data.Set("items", items)

	err = app.RenderPage(w, r, "schedule", data, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}
