FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal internal
COPY pkg pkg
COPY gen gen
COPY cmd cmd

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /character-service ./cmd/character-service

FROM scratch

COPY --from=builder /character-service /character-service

ENV PORT="8080"
ENV MONGO_URL="mongodb://mongo:27017"
ENV MONGO_DATABASE="character-service-db"

EXPOSE 8080

CMD ["/character-service"]