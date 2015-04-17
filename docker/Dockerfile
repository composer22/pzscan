FROM google/golang:1.4

MAINTAINER Bryan-Kirk Reinhardt

RUN CGO_ENABLED=0 go get -a --installsuffix cgo -ldflags '-s' github.com/composer22/pzscan
COPY Dockerfile.final /gopath/bin/Dockerfile

CMD docker build -t composer22/pzscan /gopath/bin
