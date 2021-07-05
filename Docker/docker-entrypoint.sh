#!/usr/bin/env bash
set -e

/usr/bin/mock-loki &

mkfifo /keep-open
cat /dev/stdin /keep-open | /usr/bin/promtail --stdin --config.file=/etc/promtail/config.yaml --client.url http://127.0.0.1:3100/loki/api/v1/push 1>&2

echo "DONE"

wait -n
