package httpbody

import "net/http"

func missingClose() {
	resp, _ := http.Get("https://example.com") // want `HTTP response "resp" is not closed`
	_ = resp
}

func hasClose() {
	resp, _ := http.Get("https://example.com")
	defer resp.Body.Close()
}
