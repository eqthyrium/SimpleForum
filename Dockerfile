# Stage 1: Builder
FROM golang:1.20.1-alpine3.16 AS base

RUN apk add --no-cache build-base

WORKDIR /SimpleForum

COPY . .
RUN go mod download

RUN go build -o SimpleForum cmd/main.go

# Stage 2: Final image
FROM alpine:3.16

WORKDIR /SimpleForum

# Copy all necessary files and folders (including the binary) from the builder stage
COPY --from=base /SimpleForum/SimpleForum .
COPY --from=base /SimpleForum/ui /SimpleForum/ui
COPY --from=base /SimpleForum/uploads /SimpleForum/uploads
COPY --from=base /SimpleForum/logging /SimpleForum/logging
COPY --from=base /SimpleForum/tls /SimpleForum/tls
COPY --from=base /SimpleForum/.env /SimpleForum/.env
COPY --from=base /SimpleForum/migration /SimpleForum/migration
COPY --from=base /SimpleForum/mydb.db /SimpleForum/mydb.db





# Logs directory
VOLUME ["/logging"]
VOLUME ["/uploads"]

CMD ["./SimpleForum"]
