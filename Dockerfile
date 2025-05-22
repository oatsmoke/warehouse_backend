FROM golang:1.24-alpine AS builder

WORKDIR /user/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/

FROM alpine AS runner

WORKDIR /user/src/app

COPY --from=builder /user/src/app/api .

EXPOSE 8081
CMD ["./api"]