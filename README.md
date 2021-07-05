# Promtail Debugging tool
This repository contains the source for a little tool written in Go to simplify testing of [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) configuration. It builds upon the [log piping mode](https://grafana.com/docs/loki/latest/clients/promtail/troubleshooting/#pipe-data-to-promtail) of Promtail, and adds a couple of other nice features to make it work without a Loki instance. The Go code does two things:

1. Keeps Promtail running after the standard input ends, to make it easier to work with metrics scraped from logs; and
2. Exposes a [Loki Push API endpoint](https://grafana.com/docs/loki/latest/api/#post-lokiapiv1push) that prints the received log lines to standard out in a JSON structure per line, which makes it easy to pipe to e.g. [JQ](https://stedolan.github.io/jq/).

## How to use
The tools is built and packaged as a Docker image - and released with the tag `dolittle/promtail-debug`, the tags track the versions of Promtail it uses. The basic usage entails running the Docker image, and piping in some logs:
```shell
$ docker run -i dolittle/promtail-debug < some.log
level=info ts=2021-07-05T08:19:44.332931977Z caller=server.go:225 http=[::]:9007 grpc=[::]:9095 msg="server listening on addresses"
level=info ts=2021-07-05T08:19:44.333546482Z caller=main.go:110 msg="Starting Promtail" version="(version=2.1.0, branch=HEAD, revision=1b79df375)"
{"timestamp":"2021-07-05T06:58:04.139463736Z","message":"MIT License","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141358581Z","message":"Copyright (c) 2021 dolittle-entropy","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141379737Z","message":"Permission is hereby granted, free of charge, to any person obtaining a copy","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141384497Z","message":"of this software and associated documentation files (the \"Software\"), to deal","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141390944Z","message":"in the Software without restriction, including without limitation the rights","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.14139806Z","message":"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141404809Z","message":"copies of the Software, and to permit persons to whom the Software is","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141408136Z","message":"furnished to do so, subject to the following conditions:","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141413686Z","message":"The above copyright notice and this permission notice shall be included in all","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141463719Z","message":"copies or substantial portions of the Software.","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141480901Z","message":"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141484825Z","message":"IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141490041Z","message":"FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141493376Z","message":"AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141498724Z","message":"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141501881Z","message":"OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE","labels":{"job":"debugging"}}
{"timestamp":"2021-07-05T06:58:04.141507207Z","message":"SOFTWARE.","labels":{"job":"debugging"}}
```

If you want a little more color and life, try piping it through JQ:
```shell
$ docker run -i dolittle/promtail-debug < some.log | jq
level=info ts=2021-07-05T08:19:44.332931977Z caller=server.go:225 http=[::]:9007 grpc=[::]:9095 msg="server listening on addresses"
level=info ts=2021-07-05T08:19:44.333546482Z caller=main.go:110 msg="Starting Promtail" version="(version=2.1.0, branch=HEAD, revision=1b79df375)"
{
  "timestamp": "2021-07-05T06:59:30.946869576Z",
  "message": "MIT License",
  "labels": {
    "job": "debugging"
  }
}
{
  "timestamp": "2021-07-05T06:59:30.947057188Z",
  "message": "Copyright (c) 2021 dolittle-entropy",
  "labels": {
    "job": "debugging"
  }
}
{
  "timestamp": "2021-07-05T06:59:30.947073921Z",
  "message": "Permission is hereby granted, free of charge, to any person obtaining a copy",
  "labels": {
    "job": "debugging"
  }
}
```

### Configuring
Out of the box; the Docker image doesn't really do a lot. It just prints out whatever lines you pipe into it with some fancy JSON wrapper. To really get something out of it - you need to provide a more interesting configuration. The easiest way to do that is to start by extending the [configuration file](Docker/config.yaml) and mounting it in the `promtail-debug` container at start, like this:
```shell
$ docker run -i -v "<path-to-config-on-your-machine>:/config.yaml" dolittle/promtail-debug < some.log
```

If you e.g. add some static labels in the config
```yaml
scrape_configs:
  - job_name: "debugging"
    static_configs:
    - labels:
        job: "debugging"
        hello: "world"
```
and run the Docker image from the directory where the config is saved as `config.yaml`, the result should be
```shell
$ docker run -i -v "$PWD/config.yaml:/config.yaml" dolittle/promtail-debug < some.logs | jq
level=info ts=2021-07-05T08:19:44.332931977Z caller=server.go:225 http=[::]:9007 grpc=[::]:9095 msg="server listening on addresses"
level=info ts=2021-07-05T08:19:44.333546482Z caller=main.go:110 msg="Starting Promtail" version="(version=2.1.0, branch=HEAD, revision=1b79df375)"
{
  "timestamp": "2021-07-05T07:47:42.888328017Z",
  "message": "MIT License",
  "labels": {
    "hello": "world",
    "job": "debugging"
  }
}
{
  "timestamp": "2021-07-05T07:47:42.888743469Z",
  "message": "Copyright (c) 2021 dolittle-entropy",
  "labels": {
    "hello": "world",
    "job": "debugging"
  }
}
```

> NOTE: Since Promtail is running in a special log piping mode, it only supports a single scrape config, and the Service Discovery mechanism (with relabel configs) does not do anything. This means that this debugging tool is really only useful for testing out pipeline stages for parsing/rewriting log formats and scraping metrics from the logs.

### Metrics
To test the creation of metrics from logs, setup a metric scraping pipeline in the config and expose the container port `9007` on localhost, like this:

```yaml
server:
    http_listen_port: 9007

scrape_configs:
  - job_name: "debugging"
    static_configs:
    - labels:
        job: "debugging"
    pipeline_stages:
      - match:
            selector: '{job="debugging"} |~ `dolittle`'
            stages:
                - metrics:
                    dolittle_mentions_total:
                        type: Counter
                        description: The number of log lines with 'dolittle' in them
                        config:
                            match_all: true
                            action: inc
```

```shell
$ docker run -i -v "$PWD/config.yaml:/config.yaml" -p 9007:9007 dolittle/promtail-debug < some.logs
```

Opening [http://localhost:9007/metrics](http://localhost:9007/metrics) should then include a metric called `promtail_custom_dolittle_mentions_total` that looks something like this:
```
# HELP promtail_custom_dolittle_mentions_total The number of log lines with 'dolittle' in them
# TYPE promtail_custom_dolittle_mentions_total counter
promtail_custom_dolittle_mentions_total{job="debugging"} 1
```