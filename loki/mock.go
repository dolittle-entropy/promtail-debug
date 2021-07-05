package loki

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang/snappy"
	"github.com/grafana/loki/pkg/logproto"
)

var output = json.NewEncoder(os.Stdout)

type logEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Message   string            `json:"message"`
	Labels    map[string]string `json:"labels"`
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	defer w.WriteHeader(200)

	encoded, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	decoded, err := snappy.Decode(nil, encoded)
	if err != nil {
		return
	}

	request := &logproto.PushRequest{}
	err = request.Unmarshal(decoded)
	if err != nil {
		return
	}

	for _, stream := range request.Streams {
		parser := NewLabelParser(stream.Labels)
		labels, _ := parser.Parse()

		for _, entry := range stream.Entries {
			output.Encode(logEntry{
				Timestamp: entry.Timestamp,
				Message:   entry.Line,
				Labels:    labels,
			})
		}
	}
}

func StartMockLokiServer() {
	http.HandleFunc("/loki/api/v1/push", pushHandler)
	go http.ListenAndServe(":3100", nil)
}
