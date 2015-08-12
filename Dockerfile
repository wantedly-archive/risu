FROM gliderlabs/alpine:3.2
ENTRYPOINT ["/bin/risu"]

COPY . /go/src/github.com/wantedly/risu
RUN apk-install -t build-deps go git mercurial \
      && cd /go/src/github.com/wantedly/risu \
      && export GOPATH=/go \
      && go build -ldflags "-X main.Version" -o /bin/risu \
      && rm -rf /go \
      && apk del --purge build-deps go git mercurial
