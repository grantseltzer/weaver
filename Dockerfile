FROM fedora:31 as builder

LABEL maintainer="grantseltzer@gmail.com"

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN mkdir -p /go/src/github.com/grantseltzer

COPY . /go/src/github.com/grantseltzer/oster

# Install Dependencies
RUN dnf install -y go bcc bcc-devel make

# Mount /
# privileged for now until making a seccomp profile