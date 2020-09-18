FROM golang:latest
RUN mkdir /app
RUN go get github.com/gosimple/oauth2
ADD . /app/
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]
