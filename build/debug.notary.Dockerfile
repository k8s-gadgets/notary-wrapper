FROM golang:1.14.4-alpine3.12 as base-image

RUN apk --update add sed git  gcc libc-dev




FROM base-image as notary-binary

# build notary cli binary
ARG NOTARY_BRANCH
ARG NOTARYPKG

RUN git clone -b $NOTARY_BRANCH https://${NOTARYPKG}.git /go/src/${NOTARYPKG}

WORKDIR /go/src/${NOTARYPKG}

RUN GOOS=linux GOARCH=amd64 go build -tags "pkcs11 netgo" \
    -ldflags "-w -s -X ${NOTARYPKG}/version.GitCommit=`git rev-parse --short HEAD` -X ${NOTARYPKG}/version.NotaryVersion=`cat NOTARY_VERSION` -extldflags '-static'" -o /go/bin/notary ${NOTARYPKG}/cmd/notary

RUN mkdir /user && \
    echo 'notary:x:10000:10000:notary:/:' > /user/passwd && \
    echo 'notary:x:10000:' > /user/group

RUN mkdir /.notary
