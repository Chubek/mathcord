package discord

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)



type Interaction struct {
	ID        	string          `json:"id"`
	Type     	int 			`json:"type"`
	Data      	int 			`json:"data"`
	GuildID   	string          `json:"guild_id"`
	ChannelID 	string          `json:"channel_id"`
	Message 	interface{} 	`json:"message"`
	Member 		interface{} 	`json:"member"`	
	User 		interface{} 	`json:"user"`
	Locale 		string 			`json:"locale"`	
	GuildLocale string 			`json:"guild_locale"`
	Token   	string 			`json:"token"`
	Version 	int    			`json:"version"`
}
func handler(w http.ResponseWriter, r *http.Request) {
    var interaction Interaction
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &interaction)
	
	if interaction.Type == 1 {
		w.Header().Add("content-type", "application/json")
		w.Write(json.RawMessage(`{"type": 1}`))
		
	}
}

func RunServer() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8280", nil))
}