package function

import (
	"fmt"
)

// Handle a serverless request
func Handle(req []byte, wg *sync.WaitGroup) string {
	return fmt.Sprintf("Hello, Go. You said: %s", string(req))
}
