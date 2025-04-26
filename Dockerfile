FROM golang:1.23.0-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --no-cache make

RUN go build main.go

FROM alpine AS runner

WORKDIR app

RUN apk add --no-cache curl

COPY --from=build /build/main ./main
COPY --from=build /build/.env .env

CMD ["/app/main" ]
