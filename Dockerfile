FROM golang:latest

WORKDIR /axxiom

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /usr/local/bin/axxiom ./cmd/bot

CMD ["axxiom"]
