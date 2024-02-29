FROM golang:1.22-alpine as builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /rinha cmd/main.go

FROM alpine
WORKDIR /app
USER 1000
COPY --from=builder /rinha /app/

EXPOSE 8080

CMD ["/app/rinha"]