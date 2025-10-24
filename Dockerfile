FROM golang:1.25-alpine AS builder

WORKDIR /user/src/app

RUN apk add --no-cache curl
RUN curl -sSf https://atlasgo.sh | sh

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/

FROM alpine AS runner

WORKDIR /user/src/app

RUN apk add --no-cache postgresql-client

COPY --from=builder /root/.atlas/bin/atlas /usr/local/bin/atlas
COPY --from=builder /user/src/app/api .

COPY migrations ./migrations
COPY atlas.hcl .
COPY schema/root_user.sql .

COPY scripts/entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./api"]