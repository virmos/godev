package main

import (
	"fmt"
	_ "github.com/pusher/pusher-http-go"
	_ "io"
	_ "log"
	"net/http"
	_ "strconv"
)

// PusherAuth authenticates the user to our pusher server
func (app *application) PusherAuth(w http.ResponseWriter, r *http.Request) {
	// userID := app.Session.GetInt(r.Context(), "userID")

	// u, _ := app.DB.GetUserById(userID)

	// params, _ := io.ReadAll(r.Body)

	// presenceData := pusher.MemberData{
	// 	UserID: strconv.Itoa(userID),
	// 	UserInfo: map[string]string{
	// 		"name": u.FirstName,
	// 		"id":   strconv.Itoa(userID),
	// 	},
	// }

	// response, err := app.WsClient.AuthenticatePresenceChannel(params, presenceData)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// _, _ = w.Write(response)
}

// SendPrivateMessage is sample code for sending to private channel
func (app *application) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	id := r.URL.Query().Get("id")

	data := make(map[string]string)
	data["message"] = msg

	_ = app.WsClient.Trigger(fmt.Sprintf("private-channel-%s", id), "private-message", data)
}
