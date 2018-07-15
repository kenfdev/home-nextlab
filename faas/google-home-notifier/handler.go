package function

import (
	"sync"

	"github.com/evalphobia/google-home-client-go/googlehome"
)

// Handle a serverless request
func Handle(req []byte, wg *sync.WaitGroup) string {
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: "192.168.100.4",
		Lang:     "ja",
		Accent:   "GB",
	})
	if err != nil {
		panic(err)
	}
	defer cli.QuitApp()

	cli.Notify("こんにちは、グーグル。", "ja")

	return "success"
}
