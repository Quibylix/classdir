FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 appuser
COPY --from=builder /app/server /server
USER appuser
EXPOSE 8080
CMD ["/server"]
