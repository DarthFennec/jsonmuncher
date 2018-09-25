FROM golang:1.11

RUN go get github.com/json-iterator/go
RUN go get github.com/a8m/djson
RUN go get github.com/Jeffail/gabs
RUN go get github.com/bitly/go-simplejson
RUN go get github.com/buger/jsonparser
RUN go get github.com/pquerna/ffjson
RUN go get github.com/antonholmquist/jason
RUN go get github.com/mreiferson/go-ujson
RUN go get -tags=unsafe github.com/ugorji/go/codec
RUN go get github.com/mailru/easyjson
RUN go get github.com/bcicen/jstream
RUN go get github.com/francoispqt/gojay

ADD . /go/src/github.com/darthfennec/jsonmuncher
WORKDIR /go/src/github.com/darthfennec/jsonmuncher
