FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine

WORKDIR /app

RUN adduser -S -D -H -h /app appuser
RUN chown -R appuser: /app
USER appuser

COPY --from=builder /build/main ./

EXPOSE 8080
CMD ["./main"]
