# build stage
FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/api ./cmd/server

# runtime stage
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=build /tmp/api .
USER nonroot:nonroot
CMD ["/app/api"]
