FROM golang:latest

WORKDIR /app
COPY . .

ENV GOPROXY=https://goproxy.cn,direct
RUN go build cmd/fsender/main.go

RUN ls | grep -v -E "main|examples" | xargs rm -rf

CMD [ "./main", "run", "-c", "examples/docker.yaml" ]