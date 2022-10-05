package main
import (
	"io"
	"net/http"
	"log"
	"os"
	"github.com/xuri/excelize/v2"
)

type response struct {
	Error          bool   `json:"error"`
	Message        string `json:"message"`
	Redirect       bool   `json:"redirect"`
	RedirectStatus int    `json:"status"`
	Route          string `json:"route"`
}
func (app *application) PostExcel(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("input-excel")
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	defer file.Close()
	
	f, err := os.OpenFile("./excel/"+ handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	excelize.OpenFile("Hosts.xlsx")

	var resp response
	app.writeJSON(w, http.StatusOK, resp)
}
