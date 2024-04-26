FROM golang:1.20
RUN ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone
ENV GOPROXY https://goproxy.cn
WORKDIR /app
COPY . /app
RUN go build -o main .
CMD ["./main"]