FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go download

COPY . .
RUN go build -o server .

FROM alpine:latest 

WORKDIR /srv

COPY --from=builder /app/server ./server
MKDIR -p /srv/server/public/img

CMD [ "./server" ]

