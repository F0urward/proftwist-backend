FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main cmd/${SERVICE_NAME}/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/seed db/seed/*.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin .
COPY --from=builder /app/docker/.env .
COPY --from=builder /app/config ./config/
COPY --from=builder /app/services/roadmap/repository/prompts/roadmap_example.json ./prompts/
COPY --from=builder /app/db/seed/data/roadmaps.json ./data/

EXPOSE 8080

CMD ["./main"]
