package promtail

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func RedirectStdinToPromtail() {
	cmd := exec.Command("/usr/bin/promtail", "--stdin", "--config.file=/config.yaml", "--client.url", "http://127.0.0.1:3100/loki/api/v1/push")

	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		_, err = io.Copy(stdin, os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
