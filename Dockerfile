FROM golang:stretch AS builder
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/NeoJRotary/GCB-bridge
COPY Gopkg.* ./
RUN dep ensure -vendor-only
COPY . .
RUN go build -o bridge .

FROM google/cloud-sdk:slim
WORKDIR /GCB-bridge
COPY --from=builder /go/src/github.com/NeoJRotary/GCB-bridge/bridge /GCB-bridge/bridge

CMD ["/GCB-bridge/bridge"]
