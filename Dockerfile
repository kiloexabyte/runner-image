FROM alpine:3.19

# Environment variables for Go
ENV GOLANG_VERSION=1.24.3 \
    GOROOT=/usr/local/go \
    GOPATH=/go \
    PATH=/usr/local/go/bin:/go/bin:$PATH

# Install system packages, Go, and Docker CLI
RUN apk add --no-cache \
    bash \
    curl \
    git \
    docker-cli \
    ca-certificates \
    nodejs \
    npm \
    tar \
    zip \
    aws-cli \
 && curl -LO https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz \
 && rm go${GOLANG_VERSION}.linux-amd64.tar.gz \
 && mkdir -p /go/src /go/bin /go/pkg \
 && wget https://releases.hashicorp.com/terraform/1.8.4/terraform_1.8.4_linux_amd64.zip \
 && unzip terraform_1.8.4_linux_amd64.zip -d /usr/local/bin \
 && rm terraform_1.8.4_linux_amd64.zip

# Install the 'op' tool
RUN go install lesiw.io/op@latest

# Install the 'golangci-lint' tool
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

# Default command
CMD go version && node -v
