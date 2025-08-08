##### Build Stage #####
FROM golang:alpine AS builder

RUN mkdir /app

# copying files and dirs
ADD ./lib /app/lib
ADD ./internal /app/internal
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./main.go /app/main.go

# adding required libraries (required for CGO enabled builds)
RUN apk add --no-cache build-base git

# building
WORKDIR /app
# sqlite requires CGO enabled
ENV CGO_ENABLED=1
RUN go build -o /app/corn /app/main.go

##### Final Stage #####
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/corn /app/corn

CMD [ "/app/corn", "run" ]
