FROM golang:latest
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o main .
LABEL author=leeyinghern@gmail.com
EXPOSE 8080
CMD ["/app/main"]