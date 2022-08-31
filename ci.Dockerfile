# Used for testing

FROM golang:1.18

COPY . /go/src

ENV GOPATH=/go
WORKDIR /go/src

RUN go install honnef.co/go/tools/cmd/staticcheck@latest \
 && go install github.com/mgechev/revive@latest
