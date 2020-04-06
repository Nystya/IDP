FROM golang:1.13

WORKDIR /go/src/app/
COPY . .

#RUN apt-get install -y libssl1.0.0
RUN apt install pkg-config
RUN wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb
RUN dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb
RUN rm seabolt-1.7.4-Linux-ubuntu-18.04.deb

ENV PKG_CONFIG_PATH .

RUN go build

CMD ["./Jobs"]