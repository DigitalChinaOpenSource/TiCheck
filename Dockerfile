FROM node:14.19-alpine3.15 as web

WORKDIR /webpack
COPY ./web .
#RUN yarn config set registry https://registry.npm.taobao.org
RUN yarn install
RUN yarn build

FROM golang:1.16-buster as service
WORKDIR /app
COPY . .
RUN make build

FROM ubuntu:20.04 as final

RUN DEBIAN_FRONTEND noninteractive

RUN apt update \
    && apt install -y python3 \
        python3-pip \
        mysql-client \
        expect \
    && ln -fs /usr/share/zoneinfo/Etc/UTC /etc/localtime \
    && echo 'Etc/UTC' /etc/timezone \
    && rm -rf /var /lib/apt/list/*
RUN pip3 install requests

WORKDIR /ticheck
COPY --from=web /webpack/dist ./web/dist
COPY --from=service /app/bin/ticheck-server ./service/bin
COPY ./probes ./probes
COPY ./config ./config
COPY ./executor ./excutor
COPY ./logpath.sh .
COPY ./run.sh .

WORKDIR /ticheck/service/bin
ENV GIN_MODE=release

EXPOSE 8081

ENTRYPOINT ["./ticheck-server"]


