FROM golang:1.12.7 as dev

WORKDIR $GOPATH/src/github.com/metegol-project
ADD . $GOPATH/src/github.com/metegol-project

RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/go-playground/validator

RUN govendor init
RUN govendor sync

RUN CGO_ENABLED=0 go build

FROM alpine:3.7 as prod

ENV PROJECT_DIR=/go/src/github.com/metegol-project

WORKDIR $PROJECT_DIR

COPY --from=dev $PROJECT_DIR/metegol-project .

ENTRYPOINT ./metegol-project