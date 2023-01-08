package main

import (
	"cycir/internal/models"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type response struct {
	Error          bool        `json:"error"`
	HostData       [][9]string `json:"host_data"`
	Message        string      `json:"message"`
	Redirect       bool        `json:"redirect"`
	RedirectStatus int         `json:"status"`
	Route          string      `json:"route"`
}

func (app *application) PostExcel(w http.ResponseWriter, r *http.Request) {
	var resp response
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("input-excel")
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	defer file.Close()

	f, err := os.OpenFile("./excel/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	// toggle monitoring off
	app.setSystemPref("monitoring_live", "0")
	app.toggleMonitoring("0")

	// add hosts
	hostsData, err := readHostsExcel("./excel/" + handler.Filename)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var hosts []models.Host
	for _, row := range hostsData {
		var h models.Host

		h.ID, _ = strconv.Atoi(row[0])
		h.HostName = row[1]
		h.CanonicalName = row[2]
		h.URL = row[3]
		h.IP = row[4]
		h.IPV6 = row[5]
		h.Location = row[6]
		h.OS = row[7]
		h.Active, _ = strconv.Atoi(row[8])
		hosts = append(hosts, h)
	}

	err = app.repo.BulkInsertHost(hosts)
	if err != nil {
		log.Println(err)
		return
	}
	// app.pushHostChangedEvent(hosts)
	resp.Message = "Hosts imported successfully, reset to see changes, as i am lazy"
	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) GetExcel(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	filename := r.FormValue("excel-name")
	
	// var colMap = []string {"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	hostsData, err := readHostsExcel("./excel/" + filename + ".xlsx")
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp response
	resp.HostData = hostsData
	resp.Message = "Hosts downloaded."
	app.writeJSON(w, http.StatusOK, resp)
}

func readHostsExcel(filename string) ([][9]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Println(err)
			return
		}
	}()

	var data [][9]string // len 9 = len(['id,	host_name,	canonical_name,	url,	ip,	ipv6,	location,	os,	active'])

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for i, rows := range rows {
		var row [9]string
		if i == 0 {
			continue // ignore id
		}
		for j, colCells := range rows {
			row[j] = colCells
		}
		data = append(data, row)
	}

	return data, nil
}
