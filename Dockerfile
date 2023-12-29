FROM --platform=linux/amd64 golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/extractor

CMD ["/app/bin/extractor"]
