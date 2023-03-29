# apstra-go-sdk

This project aims to be a simple-to-consume client library for Apstra.

It was initially developed to collect metric/event/anomaly/statistics kinds of
things, but could eventually support the whole Apstra API.

It's only ever been tested against AOS 4.1.0, and 4.1.1, and will complain when
asked to connect to unsupported versions of AOS.

It has three major features: Client, TwoStageL3ClosClient, and StreamTarget.

### Client
The `Client{}` object has methods closely related to Apstra platform API
endpoints, and returns data structures which closely resemble the JSON returned
by the Apstra API.

`Client` forces Apstra into "full asynchronous" mode by default to avoid issues
related to "optimistic" object ID assignment, where we've got the new object's
ID before it's ready for use. A behind-the-scenes polling function keeps track
of all outstanding tasks and returns final API results only when they're
complete. It should be safe to run many client tasks concurrently.

Logins are handled automatically.

### TwoStageL3ClosClient
The `TwoStageL3ClosClient{}` object is intended for interaction with a single
*blueprint* of the **Datacenter** reference design type. `TwoStageL3ClosClient`
has both a `Client` and a single blueprint ID embedded within.

### StreamTarget

`StreamTarget` is a listener/decoder for Apstra's "Streaming Receiver" feature.

It has Start/Stop (listening) methods and Register/Unregister methods which
add Streaming Receiver configurations via the Apstra API.

Messages and Errors are returned to the consuming code via channels.

The proto file `streaming-telemetry.proto` came from an AOS server. The easiest way to grab
one is probably via the web UI:

Click `platform -> developers` then `Rest API Documentation`.

Scroll down to `streaming-telemetry-schema-proto`, click `GET`, `Try it out` and `Execute`

Render the go code by running the following in the main project directory
```shell
protoc --go_out=.      --go_opt=Mapstra/streaming-telemetry.proto=./apstra \
       apstra/streaming-telemetry.proto
```

### Using this library

```go
package main
import "github.com/Juniper/apstra-go-sdk"
func main() {
  clientCfg := &apstra.ClientCfg{
	Url: "https://apstra-hostname",
    User:      "admin",
    Pass:      "password",
    TlsConfig: &tls.Config{InsecureSkipVerify: true},
  }
  client, _ := apstra.NewClient(clientCfg) //error ignored
  blueprintIds, _ := client.GetAllBlueprintIds(context.TODO()) //error ignored
}
```

There's an example program at `cmd/example_streaming/main.go` which implements
the streaming capability.

### Development

1. Copy `pre-push` script to `.git/hooks` to run fast validations on `git push`;
2. Use `make` targets to build, run tests or static analysis. This requires that
   your environment has all the necessary tools, if not - use `ci.Dockerfile`
   docker image.
