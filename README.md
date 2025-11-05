# apstra-go-sdk

This project aims to be a simple-to-consume client library for Apstra.

Currently supports Apstra 4.2.0 - 6.0.0. It will complain when
asked to connect to unsupported versions of Apstra.

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

### FreeformClient
The `FreeformClient{}` object is intended for interaction with a single
*blueprint* of the **Freeform** reference design type. `FreeformClient`
has both a `Client` and a single blueprint ID embedded within.


### StreamTarget
`StreamTarget` is a listener/decoder for Apstra's "Streaming Receiver" feature.

It has Start/Stop (listening) methods and Register/Unregister methods which
add Streaming Receiver configurations via the Apstra API.

Messages and Errors are returned to the consuming code via channels.

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
