FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /grocademy ./cmd/grocademy/main.go

FROM golang:1.24.2-alpine3.21 AS migrator
WORKDIR /migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /grocademy .
COPY --from=migrator /go/bin/migrate .
COPY --from=builder /app/migrations ./migrations/
COPY --from=builder /app/web/templates ./web/templates/
COPY --from=builder /app/web/static ./web/static/
EXPOSE 8080
CMD ["./grocademy"]
