FROM golang:1.26-bookworm

WORKDIR /app
COPY . .


RUN go mod tidy
RUN go build -o /app/exe ./cmd/main.go

CMD ["/app/exe"]