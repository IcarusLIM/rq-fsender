FROM golang:latest

WORKDIR /app
COPY . .

ENV GOPROXY=https://goproxy.cn,direct
RUN go build app/main.go

RUN ls | grep -v -E "main" | xargs rm -rf
RUN mkdir files

CMD [ "./main" ]