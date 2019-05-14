FROM alpine:latest

LABEL author zs5460@gmail.com
LABEL appname qiniuservice 

WORKDIR /app

ADD ./qiniusvc ./qiniusvc
ADD ./public/ ./public/

EXPOSE 80

ENTRYPOINT [ "./qiniusvc" ]
