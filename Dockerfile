FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

# need to run primsa items
RUN go get github.com/steebchen/prisma-client-go
RUN go run github.com/steebchen/prisma-client-go generate
RUN go get github.com/steebchen/prisma-client-go/engine@v0.47.0

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


