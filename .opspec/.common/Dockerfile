FROM golang:1.8

RUN go get \
  github.com/kardianos/govendor \
  github.com/maxbrunsfeld/counterfeiter \
  github.com/onsi/ginkgo/ginkgo

# equinox.io build tool
RUN curl -o equinox.tgz https://bin.equinox.io/c/mBWdkfai63v/release-tool-stable-linux-amd64.tgz && \
tar -vxf equinox.tgz && \
mv ./equinox /usr/local/bin
