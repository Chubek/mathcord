package discord

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"mathcord/ed25519"
	"mathcord/utils"
	"net/http"
	"os"
)

type Interaction struct {
	ID          string      `json:"id"`
	Type        int         `json:"type"`
	Data        int         `json:"data"`
	GuildID     string      `json:"guild_id"`
	ChannelID   string      `json:"channel_id"`
	Message     interface{} `json:"message"`
	Member      interface{} `json:"member"`
	User        interface{} `json:"user"`
	Locale      string      `json:"locale"`
	GuildLocale string      `json:"guild_locale"`
	Token       string      `json:"token"`
	Version     int         `json:"version"`
}

var PK string

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Print(err)
	}

	PK = os.Getenv("APP_PK")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print(r)

	w.Header().Add("content-type", "application/json")

	var interaction Interaction
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(body, &interaction)

	if err != nil {
		log.Print(err)
	}

	if interaction.Type == 1 {
		_, err = w.Write(json.RawMessage(`{"type": 1}`))

		if err != nil {
			log.Print(err)
		}

	}

	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")

	tsBytes := []byte(timestamp)

	var m []byte

	m = append(m, tsBytes...)
	m = append(m, body...)

	verified := ed25519.CheckValidBytes(m, utils.HexToByte(signature), utils.HexToByte(PK))

	if !verified {
		w.WriteHeader(401)
	}

	fmt.Println(interaction)

}

func RunServer() {
	http.HandleFunc("/", handler)
	log.Print(http.ListenAndServe(":8050", nil))
}
