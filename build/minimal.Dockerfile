FROM golang:1.14.4-alpine3.12 as base-image

RUN apk --update add sed git  gcc libc-dev




FROM base-image as notary-binary

# build notary cli binary
ARG NOTARY_BRANCH
ARG NOTARYPKG

RUN git clone -b $NOTARY_BRANCH https://github.com/theupdateframework/notary.git /go/src/${NOTARYPKG}

WORKDIR /go/src/${NOTARYPKG}

RUN GOOS=linux GOARCH=amd64 go build -tags "pkcs11 netgo" \
    -ldflags "-w -s -X ${NOTARYPKG}/version.GitCommit=`git rev-parse --short HEAD` -X ${NOTARYPKG}/version.NotaryVersion=`cat NOTARY_VERSION` -extldflags '-static'" -o /go/bin/notary ${NOTARYPKG}/cmd/notary

RUN mkdir /user && \
    echo 'notary:x:10000:10000:notary:/:' > /user/passwd && \
    echo 'notary:x:10000:' > /user/group

RUN mkdir /.notary




FROM base-image as wrapper-binary

# build notary-wrapper binary
ARG NOTARY_WRAPPER_BRANCH
ARG NOTARY_WRAPPER_PKG

RUN git clone -b $NOTARY_WRAPPER_BRANCH https://github.com/k8s-gadgets/notary-wrapper.git /go/src/${NOTARY_WRAPPER_PKG}

WORKDIR /go/src/${NOTARY_WRAPPER_PKG}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-w -s -X ${NOTARY_WRAPPER_PKG}/version.GitCommit=`git rev-parse --short HEAD` -X ${NOTARY_WRAPPER_PKG}/version.NotaryWrapperVersion=`cat NOTARY_WRAPPER_VERSION`" -o /go/bin/notary-wrapper ${NOTARY_WRAPPER_PKG}/cmd




FROM scratch

COPY --from=notary-binary /.notary /.notary
COPY --from=notary-binary /user/group /user/passwd /etc/
COPY --from=notary-binary /go/bin/notary /notary/notary
COPY --from=wrapper-binary /go/bin/notary-wrapper /notary/notary-wrapper
COPY --from=wrapper-binary /etc/ssl /etc/ssl

ENV NOTARY_PORT "4445"
ENV NOTARY_CERT_PATH "/etc/certs/notary"
ENV NOTARY_ROOT_CA "root-ca.crt"
ENV NOTARY_CLI_PATH "/notary/notary"

USER notary:notary

EXPOSE 4445

WORKDIR /notary

ENTRYPOINT [ "/notary/notary-wrapper" ]
