FROM golang:alpine

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

# EXPOSE 5001

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main