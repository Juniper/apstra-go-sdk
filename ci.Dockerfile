# Used for testing

FROM golang:1.18

RUN apt-get install -y git \
 && go install honnef.co/go/tools/cmd/staticcheck@latest \
 && go install github.com/mgechev/revive@latest \
 && rm -rf /var/lib/apt/lists/*
