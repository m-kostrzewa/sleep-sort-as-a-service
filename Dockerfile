FROM golang:1.11

WORKDIR /src
COPY . .

RUN go mod download
RUN go build -v -o /app .

EXPOSE 8080
ENTRYPOINT ["/app"]
