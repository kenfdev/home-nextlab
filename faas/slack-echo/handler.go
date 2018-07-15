package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
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

// Handle a serverless request
func Handle(req []byte, wg *sync.WaitGroup) string {
	var event CloudEvent
	err := json.Unmarshal(req, &event)
	if err != nil {
		panic(err)
	}

	token, _ := getAPISecret("bot-user-oauth-access-token")

	api := slack.New(string(token))

	wg.Add(1)
	go func() {

		defer wg.Done()

		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage([]byte(event.Data.Body)), slackevents.OptionVerifyToken(&slackevents.TokenComparator{string(token)}))
		fmt.Printf("eventsAPIEvent: %+v\n", eventsAPIEvent)
		if e != nil {
			// w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(req), &r)
			if err != nil {
				// w.WriteHeader(http.StatusInternalServerError)
			}
			// w.Header().Set("Content-Type", "text")
			// w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			postParams := slack.PostMessageParameters{}
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				api.PostMessage(ev.Channel, "Yes, hello.", postParams)
			}
		}

	}()
	return fmt.Sprintf("Hello, Go. You said: %s", string(req))
}
