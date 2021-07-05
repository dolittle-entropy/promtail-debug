package main

import (
	"promtail-debug/loki"
	"promtail-debug/promtail"
)

func main() {
	loki.StartMockLokiServer()
	promtail.RedirectStdinToPromtail()
}
