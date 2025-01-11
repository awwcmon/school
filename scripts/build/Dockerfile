FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# set the time zone to Shanghai
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY school /app/school
RUN chmod +x /app/school

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./school", "-c", "configs/school.yml"]
# if you use the Configuration Center, school.yml is changed to the Configuration Center configuration.
#CMD ["./school", "-c", "configs/school.yml", "-enable-cc"]
