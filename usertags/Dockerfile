FROM golang:1.13.4

WORKDIR /app

COPY ./usertags/go.mod ./usertags/go.sum ./
RUN go mod download

COPY ./usertags .

EXPOSE 8001
