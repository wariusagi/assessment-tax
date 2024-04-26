# Build stage

FROM golang:1.21-alpine AS build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/go-assessment-tax .

# Run stage

FROM alpine:3.19.1

COPY --from=build-base /app/out/go-assessment-tax /app/go-assessment-tax

CMD ["/app/go-assessment-tax"]
