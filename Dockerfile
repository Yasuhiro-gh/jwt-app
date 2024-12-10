FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /jwtapp ./cmd/jwt_app/

FROM alpine:3

WORKDIR /app
COPY --from=builder ./jwtapp ./jwtapp

EXPOSE 8080

CMD ["/app/jwtapp"]