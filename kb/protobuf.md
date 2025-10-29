### protobuf

The proto file `streaming-telemetry.proto` came from an AOS server. The easiest way to grab
one is probably via the web UI:

Click `platform -> developers` then `Rest API Documentation`.

Scroll down to `streaming-telemetry-schema-proto`, click `GET`, `Try it out` and `Execute`

Render the go code by running the following in the main project directory
```shell
protoc --go_out=.      --go_opt=Mapstra/streaming-telemetry.proto=./apstra \
       apstra/streaming-telemetry.proto
```
