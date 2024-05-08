FROM golang:1.15.8-alpine


RUN apk add bash git
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY="https://proxy.golang.org,direct"
RUN go get github.com/gin-contrib/cors github.com/gin-gonic/gin github.com/go-sql-driver/mysql 
COPY . /data/

WORKDIR /data
# GO111MODULE=on
# RUN go run /data/database.go
# CMD bash
RUN bash sEngine.sh setup