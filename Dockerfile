FROM golang:1.18.3-alpine
RUN apk update
RUN apk add ffmpeg
RUN mkdir /app
COPY ./ /app/
WORKDIR /app/
RUN go build -o main ./cmd
CMD ["./main"]