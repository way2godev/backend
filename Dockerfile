ARG GOLANG_VERSION=1.22.2
FROM go:${GOLANG_VERSION} as builder

WORKDIR /app

# User
ARG UID=10001
RUN --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

FROM alpine:latest as production

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# User
COPY --from=builder /etc/passwd /etc/passwd
USER appuser

ENV PORT=5000
EXPOSE ${PORT}

CMD ["./main"]