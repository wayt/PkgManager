FROM dockerfile/go

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get install -y mercurial

RUN go get github.com/my-sign/gosto

WORKDIR /gopath/src/github.com/my-sign/gosto

RUN go build && \
    go install

ADD gosto.conf /gopath/bin/gosto.conf

WORKDIR /gopath/bin

CMD /gopath/bin/gosto -config=/gopath/bin/gosto.conf

EXPOSE 1024
