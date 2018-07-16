FROM golang

WORKDIR /go/src/keep
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV DOCKER_API_VERSION 1.35

CMD ["keep"]
