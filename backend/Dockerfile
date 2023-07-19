FROM golang:1.20-alpine

WORKDIR /src

COPY . .
RUN rm -rf configs/*.yml && go mod download && go build -o main . && mkdir -p /app && mv main /app && rm -rf /src
WORKDIR /app
COPY configs/config-docker.yml configs/config.yml
COPY scripts scripts

EXPOSE 8080

CMD ["./main"]
