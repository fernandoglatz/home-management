FROM golang:1.20-alpine

WORKDIR /src

COPY . .
RUN rm -rf config && go mod download && go build -o main . && mkdir -p /app/resources/config && mv main /app && rm -rf /src
WORKDIR /app
COPY resources/config/config-docker.yml resources/config/config.yml
COPY resources/db resources/db

EXPOSE 8080

CMD ["./main"]
