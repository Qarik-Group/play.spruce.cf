FROM concourse/concourse-ci
# https://github.com/concourse/concourse/blob/master/ci/dockerfiles/concourse-ci/Dockerfile

# Get some basic packages
RUN apt-get update && apt-get install curl wget bzr git -y

# Install Go
RUN \
  mkdir -p /goroot && \
  curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | \
    tar xvzf - -C /goroot --strip-components=1

# Set environment variables.
ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

RUN go get github.com/tools/godep

# Install the cf-cli, which may have been updated by Concourse
ADD cf-cli/cf-linux-amd64.tgz /usr/bin/
