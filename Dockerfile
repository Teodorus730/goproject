FROM golang:1.26-alpine AS builder

WORKDIR /build


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./waterService cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /build/waterService /bin/waterService

CMD ["/bin/waterService"]
