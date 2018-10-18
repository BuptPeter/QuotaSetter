FROM golang:1.10

WORKDIR /go/src
RUN mkdir -p jd.com/quota-setter
#挂载ceph
COPY .  jd.com/quota-setter/

RUN go install jd.com/quota-setter

WORKDIR jd.com/quota-setter

CMD go run *.go

