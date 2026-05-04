FROM golang:1.26.2-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /team-app-server ./cmd/server

FROM gcr.io/distroless/static:nonroot@sha256:e3f945647ffb95b5839c07038d64f9811adf17308b9121d8a2b87b6a22a80a39

COPY --from=builder /team-app-server /team-app-server
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/team-app-server"]
