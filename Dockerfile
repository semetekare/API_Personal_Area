FROM golang:1.21.3-alpine as builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . ./

RUN go build -o ./bin/cmd cmd/main.go

FROM alpine:latest AS runner

COPY --from=builder /app/bin/cmd ./
COPY .env /

CMD ["./cmd"]