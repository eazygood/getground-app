# FROM golang:1.16-alpine

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o bin/app cmd/app/main.go

# EXPOSE 8081

# CMD ["./bin/app"]

FROM golang:1.19-alpine

WORKDIR /app

COPY config.yaml /app/config.yaml
COPY bin/app /app

CMD ["./app"]
