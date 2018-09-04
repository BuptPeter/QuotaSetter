FROM ai-image.jd.com/go/golang:1.10

WORKDIR /go/src
#RUN mkdir -p jd.com/gthrottle
#挂载ceph
#RUN mount -t ceph 11.7.148.100,11.7.148.101,11.7.148.102:6789:/ /mnt/cephfs-ht-test2
COPY . ./jd.com/quota-setter
RUN go install jd.com/quota-setter
CMD ["go","run","jd.com/quota-setter/*.go"]
