FROM golang:1.24 AS build-stage
WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main /app/cmd/main.go

FROM alpine:latest AS release-stage
WORKDIR /app

COPY --from=build-stage /app/main .

EXPOSE 8080

ENTRYPOINT ["./main"]
