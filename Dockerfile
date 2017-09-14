FROM golang:1.9

WORKDIR /go/src/actuator

COPY . ./

RUN    go-wrapper download \
    && go-wrapper install \
    && mkdir /actuator \
    && cp /go/bin/actuator /actuator/ \
    && cp /go/src/actuator/actuator.yml /actuator/ \
    && rm -rf /go/src

WORKDIR /actuator

USER 10000
EXPOSE 8080
CMD ["/actuator/actuator"]
