FROM notary-binary as wrapper-binary

# build notary-wrapper binary
ARG NOTARY_WRAPPER_BRANCH
ARG NOTARY_WRAPPER_PKG

RUN git clone -b $NOTARY_WRAPPER_BRANCH https://${NOTARY_WRAPPER_PKG}.git /go/src/${NOTARY_WRAPPER_PKG}

WORKDIR /go/src/${NOTARY_WRAPPER_PKG}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-w -s -X ${NOTARY_WRAPPER_PKG}/version.GitCommit=`git rev-parse --short HEAD` -X ${NOTARY_WRAPPER_PKG}/version.NotaryWrapperVersion=`cat NOTARY_WRAPPER_VERSION`" -o /go/bin/notary-wrapper ${NOTARY_WRAPPER_PKG}/cmd




FROM golang:1.14.4-alpine3.12

COPY --from=notary-binary /.notary /.notary
COPY --from=notary-binary /user/group /user/passwd /etc/
COPY --from=notary-binary /go/bin/notary /notary/notary
COPY --from=wrapper-binary /go/bin/notary-wrapper /notary/notary-wrapper
COPY --from=wrapper-binary /etc/ssl /etc/ssl


## start local only (docker run ...)
# use secret for k8s
# COPY notary-wrapper.crt  /etc/certs/notary/notary-wrapper.crt
# COPY notary-wrapper.key  /etc/certs/notary/notary-wrapper.key
# RUN chown -R notary:notary /etc/certs/notary
# RUN chown -R notary:notary /notary
# RUN chown -R notary:notary /.notary
## end local only


ENV NOTARY_PORT "4445"
ENV NOTARY_CERT_PATH "/etc/certs/notary"
ENV NOTARY_ROOT_CA "root-ca.crt"
ENV NOTARY_CLI_PATH "/notary/notary"

USER notary:notary

EXPOSE 4445

WORKDIR /notary

ENTRYPOINT [ "/notary/notary-wrapper" ]
