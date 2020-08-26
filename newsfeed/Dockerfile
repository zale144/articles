FROM golang:1.13.4

WORKDIR /app

COPY ./newsfeed/go.mod ./newsfeed/go.sum ./
RUN go mod download

COPY ./newsfeed ./

EXPOSE 8002
