FROM golang:1.9 AS build

RUN go-wrapper download github.com/gin-gonic/gin
RUN go-wrapper download github.com/google/go-github/github
RUN go-wrapper download github.com/spf13/afero

WORKDIR /go/src/actuator
COPY . ./

RUN go-wrapper download
RUN go install


#### RUN IMAGE #######
FROM debian:stretch AS run

WORKDIR /actuator
COPY --from=build /go/bin/actuator .
COPY actuator.yml .

USER 1000
EXPOSE 8080
CMD "/actuator/actuator"
