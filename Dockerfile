FROM golang:1.12.5-stretch

RUN go get github.com/aws/aws-sdk-go

COPY ./aem-s3-logsync/ /go/src/github.com/benjamin-maynard/aem-s3-logsync/aem-s3-logsync

RUN go install github.com/benjamin-maynard/aem-s3-logsync/aem-s3-logsync

COPY resources/entrypoint.sh /

RUN chmod +x /entrypoint.sh

ENTRYPOINT /entrypoint.sh