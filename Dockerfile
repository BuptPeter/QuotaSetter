FROM ai-image.jd.com/go/golang:1.10

WORKDIR /go/src
RUN mkdir -p ./jd.com/quota-setter
#挂载ceph
COPY . ./jd.com/quota-setter
RUN go install ./jd.com/quota-setter
#CMD ["./jd.com/quota-setter/quota-setter"]
CMD ["go","run","./jd.com/quota-setter/*.go""]
