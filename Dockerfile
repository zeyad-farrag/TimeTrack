FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /team-app-server ./cmd/server

FROM gcr.io/distroless/static

COPY --from=builder /team-app-server /team-app-server
EXPOSE 8080
ENTRYPOINT ["/team-app-server"]
