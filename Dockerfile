FROM --platform=linux/amd64 golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

COPY storage ./storage

COPY external ./external

COPY service ./service

RUN go build -o ./bin/filing-extractor

CMD ["/app/bin/filing-extractor"]
