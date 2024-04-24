FROM go:1.22.2 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest as production

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]