FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest

WORKDIR /srv

COPY --from=builder /app/server ./server
RUN mkdir -p ./public/img

CMD [ "./server" ]