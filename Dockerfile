ARG GOLANG_VERSION=1.22.2
FROM golang:${GOLANG_VERSION}-alpine as builder

WORKDIR /app

# User
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN mkdir bin
RUN go build -o bin/main .

# For each folder in ./microservices, cd into it and build the binary
RUN for d in ./microservices/*; do (cd "$d" && go build -o ../../bin/$(basename $d) .); done

FROM alpine:latest as production

WORKDIR /app

COPY --from=builder /etc/passwd /etc/passwd

RUN mkdir data && \
    chown -R appuser /app && \
    chmod -R 755 /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin ./bin

# Copy the ./microservices/gtfs-parser/sources.csv
COPY --from=builder /app/microservices/gtfs-parser/sources.csv ./microservices/gtfs-parser/sources.csv

# User
USER appuser

ENV PORT=5000
EXPOSE ${PORT}

CMD ["./bin/main"]