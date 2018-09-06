FROM google/cloud-sdk:slim

# install golang
RUN curl -O https://dl.google.com/go/go1.11.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.11.linux-amd64.tar.gz
RUN export PATH="/usr/local/go/bin:$PATH"; go version
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"


# install and run dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/NeoJRotary/GCB-bridge
COPY Gopkg.* ./
RUN dep ensure -vendor-only