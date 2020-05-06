# BUILDER
FROM fedora:31 as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

WORKDIR /go/src/github.com/grantseltzer/weaver

COPY . /go/src/github.com/grantseltzer/weaver

RUN dnf install -y go bcc bcc-devel make && \
     make && \
     bash -c "./build-helper.sh /tmp/build-dir"


FROM scratch  

ENV LD_LIBRARY_PATH /lib64
ENV PATH /bin
COPY --from=builder /usr/bin/ldd /bin/ldd
COPY --from=builder /tmp/build-dir/* /lib64
COPY --from=builder /go/src/github.com/grantseltzer/weaver/bin/weaver /bin/weaver
ENTRYPOINT ["/bin/weaver"]
