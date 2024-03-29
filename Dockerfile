FROM golang:latest

WORKDIR /app

COPY . .
RUN CGO_ENABLED=1 go build
RUN mkdir -p db

EXPOSE 8080

ENV PORT 8080

VOLUME db db

CMD ["./todo"]
