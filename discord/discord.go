package discord

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"mathcord/ed25519"
	"mathcord/parser"
	"mathcord/utils"
	"net/http"
	"os"
)

type ResponseContent struct {
	Content string `json:"content"`
}

type ResponseData struct {
	Type int             `json:"type"`
	Data ResponseContent `json:"data"`
}

type Options struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Value string `json:"value"`
}

type Data struct {
	ID      int     `json:"id"`
	Name    string  `json:"string"`
	Options Options `json:"options"`
}

type Interaction struct {
	ID          string      `json:"id"`
	Type        int         `json:"type"`
	Data        interface{} `json:"data"`
	GuildID     string      `json:"guild_id"`
	ChannelID   string      `json:"channel_id"`
	Message     interface{} `json:"message"`
	Member      interface{} `json:"member"`
	User        interface{} `json:"user"`
	Locale      *string     `json:"locale"`
	GuildLocale string      `json:"guild_locale"`
	Token       string      `json:"token"`
	Version     int         `json:"version"`
}

var PK string

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Println(err)
	}

	PK = os.Getenv("APP_PK")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("------------ > START < ---------------")

	log.Println(r)

	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	tsBytes := []byte(timestamp)

	var m []byte

	m = append(m, tsBytes...)
	m = append(m, body...)

	verified := ed25519.CheckValidBytes(m, utils.HexToByte(signature), utils.HexToByte(PK))

	if !verified {
		log.Println("Failed to verify")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("invalid request signature"))

		if err != nil {
			log.Println(err)
		}

	} else {
		w.WriteHeader(http.StatusOK)

		log.Println("Status set to OK")

	}

	var interaction Interaction

	err = json.Unmarshal(body, &interaction)

	if err != nil {
		log.Println(err)
	}

	log.Println("Intersection is: \n", interaction)

	if interaction.Type == 1 {
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(json.RawMessage(`{"type": 1}`))

		if err != nil {
			log.Println(err)
		}

	}

	if interaction.Data != nil {
		data := interaction.Data.(map[string]interface{})
		val, _ := data["options"].([]interface{})
		options, _ := val[0].(map[string]interface{})
		exp_, _ := options["value"]
		exp := exp_.(string)

		result := parser.ShuntingYard(exp)
		log.Println("Result is ", result)

		respContent := ResponseContent{Content: "Result is: " + result}
		respData := ResponseData{
			Data: respContent,
			Type: 4,
		}

		jData, err := json.Marshal(respData)
		if err != nil {
			// handle error
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jData)

		if err != nil {
			return
		}

	}

	log.Println("------------ > END < ---------------")

}

func RunServer() {
	http.HandleFunc("/", handler)
	log.Println(http.ListenAndServe(":8050", nil))
}
