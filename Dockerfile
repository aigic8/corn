FROM golang:alpine

RUN mkdir /app

# copying files and dirs
ADD ./lib /app/lib
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./main.go /app/main.go

# building
WORKDIR /app
RUN go build -o /app/corn /app/main.go

CMD [ "/app/corn" ]
