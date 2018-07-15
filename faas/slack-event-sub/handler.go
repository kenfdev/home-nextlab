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
	Path    string                 `json:"path"`
	Method  string                 `json:"method"`
	Headers map[string]string      `json:"headers"`
	Host    string                 `json:"host"`
	Query   map[string]string      `json:"query"`
	Params  map[string]string      `json:"params"`
	Body    map[string]interface{} `json:"body"`
}

// Handle a serverless request
func Handle(req []byte, wg *sync.WaitGroup) string {
	oauthToken, _ := getAPISecret("bot-user-oauth-access-token")
	verifyToken, _ := getAPISecret("slack-verify-token")
	var api = slack.New(string(oauthToken))

	var event CloudEvent
	err := json.Unmarshal(req, &event)
	if err != nil {
		panic(err)
	}

	body, _ := json.Marshal(event.Data.Body)
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: string(verifyToken)}))
	if e != nil {
		panic(e)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			panic(err)
		}

		return fmt.Sprintf(`{"headers": {"Content-Type": "text/plain"}, "body": "%s"}`, r.Challenge)
	}

	wg.Add(1)
	go func() {

		defer wg.Done()

		fmt.Printf("eventAPIEvent: %+v\n", eventsAPIEvent)
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			postParams := slack.PostMessageParameters{}
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				api.PostMessage(ev.Channel, "Yes, hello.", postParams)
			}
		}

	}()
	params := &slack.Msg{Text: "Hi"}
	b, _ := json.Marshal(params)
	return string(b)
}
