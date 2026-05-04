FROM golang:1.23 AS builder
WORKDIR /src
COPY . .
RUN go mod init ddtrace-sender || true
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/ddtrace-sender main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/ddtrace-sender /ddtrace-sender
ENTRYPOINT ["/ddtrace-sender"]
