package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"sync"

	"github.com/nlopes/slack"
)

func getAPISecret(secretName string) (secretBytes []byte, err error) {
	// read from the openfaas secrets folder
	secretBytes, err = ioutil.ReadFile("/var/openfaas/secrets/" + secretName)
	if err != nil {
		// read from the original location for backwards compatibility with openfaas <= 0.8.2
		secretBytes, err = ioutil.ReadFile("/run/secrets/" + secretName)
	}

	return secretBytes, err
}

type CloudEvent struct {
	EventType          string            `json:"eventType"`
	EventID            string            `json:"eventID"`
	CloudEventsVersion string            `json:"cloudEventsVersion"`
	Source             string            `json:"source"`
	EventTime          string            `json:"eventTime"`
	Data               *HTTPRequestEvent `json:"data"`
	ContentType        string            `json:"contentType"`
}

type HTTPRequestEvent struct {
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Host    string            `json:"host"`
	Query   map[string]string `json:"query"`
	Params  map[string]string `json:"params"`
	Body    string            `json:"body"`
}

func SlashCommandParse(str string) (sc slack.SlashCommand, err error) {
	vals, err := url.ParseQuery(str)
	if err != nil {
		return sc, err
	}
	sc.Token = vals.Get("token")
	sc.TeamID = vals.Get("team_id")
	sc.TeamDomain = vals.Get("team_domain")
	sc.EnterpriseID = vals.Get("enterprise_id")
	sc.EnterpriseName = vals.Get("enterprise_name")
	sc.ChannelID = vals.Get("channel_id")
	sc.ChannelName = vals.Get("channel_name")
	sc.UserID = vals.Get("user_id")
	sc.UserName = vals.Get("user_name")
	sc.Command = vals.Get("command")
	sc.Text = vals.Get("text")
	sc.ResponseURL = vals.Get("response_url")
	sc.TriggerID = vals.Get("trigger_id")
	return sc, nil
}

// Handle a serverless request
func Handle(req []byte, wg *sync.WaitGroup) string {
	var event CloudEvent
	err := json.Unmarshal(req, &event)
	if err != nil {
		panic(err)
	}

	// token, _ := getAPISecret("bot-user-oauth-access-token")

	wg.Add(1)
	go func() {

		defer wg.Done()

		s, err := SlashCommandParse(event.Data.Body)
		if err != nil {
			panic(err)
		}

		// if !s.ValidateToken(verificationToken) {
		// 	panic()
		// }

		fmt.Printf("SlashCommand: %+v\n", s)
		// switch s.Command {
		// case "/echo":
		// 	params := &slack.Msg{Text: s.Text}
		// 	b, err := json.Marshal(params)
		// 	if err != nil {
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		return
		// 	}
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.Write(b)
		// default:
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// if e != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		// }

	}()
	params := &slack.Msg{Text: "Hi"}
	b, _ := json.Marshal(params)
	return string(b)
}
