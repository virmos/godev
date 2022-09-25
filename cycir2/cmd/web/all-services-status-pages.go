package main

import (
	"github.com/CloudyKit/jet/v6"
	"log"
	"net/http"
)

// AllHealthyServices lists all healthy services
func (app *application) AllHealthyServices(w http.ResponseWriter, r *http.Request) {
	// get all host services (with host info) for status pending
	services, err := app.DB.GetServicesByStatus("healthy")
	if err != nil {
		log.Println(err)
		return
	}
	vars := make(jet.VarMap)
	vars.Set("services", services)
	err = app.RenderPage(w, r, "healthy", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllWarningServices lists all warning services
func (app *application) AllWarningServices(w http.ResponseWriter, r *http.Request) {
	// get all host services (with host info) for status pending
	services, err := app.DB.GetServicesByStatus("warning")
	if err != nil {
		log.Println(err)
		return
	}
	vars := make(jet.VarMap)
	vars.Set("services", services)
	err = app.RenderPage(w, r, "warning", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllProblemServices lists all problem services
func (app *application) AllProblemServices(w http.ResponseWriter, r *http.Request) {
	// get all host services (with host info) for status pending
	services, err := app.DB.GetServicesByStatus("problem")
	if err != nil {
		log.Println(err)
		return
	}
	vars := make(jet.VarMap)
	vars.Set("services", services)
	err = app.RenderPage(w, r, "problems", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllPendingServices lists all pending services
func (app *application) AllPendingServices(w http.ResponseWriter, r *http.Request) {
	// get all host services (with host info) for status pending
	services, err := app.DB.GetServicesByStatus("pending")
	if err != nil {
		log.Println(err)
		return
	}
	vars := make(jet.VarMap)
	vars.Set("services", services)

	err = app.RenderPage(w, r, "pending", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}