package function

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/evalphobia/google-home-client-go/googlehome"
)

type Request struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

// Handle a serverless request
func Handle(req []byte, wg *sync.WaitGroup) string {
	host := os.Getenv("GOOGLE_HOME_HOST")
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: host,
	})
	if err != nil {
		panic(err)
	}

	var r Request
	err = json.Unmarshal(req, &r)
	if err != nil {
		panic(err)
	}

	err = cli.Notify(r.Text, r.Language)
	if err != nil {
		panic(err)
	}

	return "success"
}
