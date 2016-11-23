FROM gliderlabs/alpine:3.4
RUN apk update
RUN apk --update add bash curl go git mercurial gcc linux-headers musl-dev libc-dev

RUN curl -Ls https://github.com/progrium/execd/releases/download/v0.1.0/execd_0.1.0_Linux_x86_64.tgz \
    | tar -zxC /bin \
  && curl -Ls https://github.com/progrium/entrykit/releases/download/v0.2.0/entrykit_0.2.0_Linux_x86_64.tgz \
    | tar -zxC /bin \
  && curl -sL https://get.docker.com/builds/Linux/x86_64/docker-1.7.1 > /bin/docker \
  && chmod +x /bin/docker \
  && entrykit --symlink

ADD ./data /tmp/data

ENV GOPATH /go
COPY . /go/src/github.com/gophertrain/envy
WORKDIR /go/src/github.com/gophertrain/envy
RUN go version
RUN go get && CGO_ENABLED=0 go build -a -installsuffix cgo -o /bin/envy \
  && ln -s /bin/envy /bin/enter \
  && ln -s /bin/envy /bin/auth \
  && ln -s /bin/envy /bin/serve

VOLUME /envy
EXPOSE 22 80
ENTRYPOINT ["codep", "/bin/execd -e -k /tmp/data/id_host /bin/auth /bin/enter", "/bin/serve"]
