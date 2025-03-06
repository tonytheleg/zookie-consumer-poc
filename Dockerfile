FROM golang:1.23

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go get github.com/confluentinc/confluent-kafka-go/kafka && go build -o /consumer
CMD ["/consumer"]
