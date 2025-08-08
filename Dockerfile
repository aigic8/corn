# Build stage
FROM golang:alpine AS builder

RUN mkdir /app

# copying files and dirs
ADD ./lib /app/lib
ADD ./internal /app/internal
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./main.go /app/main.go

# building
WORKDIR /app
RUN CGO_ENABLED=1 go build -o /app/corn /app/main.go

# Final stage
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/corn /app/corn

CMD [ "/app/corn", "run" ]
