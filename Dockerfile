FROM golang:1.23 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-mall main.go

FROM golang:1.23-alpine

ENV ENV=dev

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/go-mall /usr/local/bin/go-mall

EXPOSE 8080

CMD ["go-mall"]
