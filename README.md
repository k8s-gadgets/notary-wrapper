# notary-wrapper
this is a notary cli rest interface to get digests via http request  
**this code it is not production ready!**  
for this applicaiton to work a running notary instance is necessary  
if you want to use a private notary instance have a look at following repo: https://github.com/k8s-gadgets/k8s-content-trust  


## env
| Parameter                                   | Description                               | Default                                    |
| ------------------------------------------  | ----------------------------------------  | -------------------------------------------|
| `NOTARY_PORT`                               | port for notary-wrapper                   | `4445`                                     |
| `NOTARY_CERT_PATH`                          | path for certificate folders              | `/etc/certs/notary`                        |
| `NOTARY_ROOT_CA`                            | name of root-ca's                         | `root-ca.crt`                              |
| `NOTARY_CLI_PATH`                           | path for notary cli                       | `/notary/notary`                           |


## certs
- two certs necessary (mount via secret):
  - notary-wrapper cert and key for serving https
  - certs were created with the script in the `k8s-content-trust/notary-k8s/helm/notary/generateCerts.sh` repo 
- optional:
  - root-ca.crt to be able to communicate with notary-server
  - root-ca.crt:
    - ```NOTARY_CERT_PATH``` = "/etc/certs/notary" (defult)
    - mount the cert for each notary in different folder:
    - ```NOTARY_ROOT_CA``` = root-ca.cert  (same name for all root-ca's; in different folders)
    - naming convention folders:
      - name of service: ```notary-server-svc.notary.svc``` 
      - name of folder: ```notary-server-svc:4443```


## local usage
create a new entry in /etc/hosts since certs are only valid for notary-wapper-svc but not for localhost  
```127.0.0.1  notary-wrapper-svc```

## list
### with root-ca
```
curl -X POST https://notary-wrapper-svc:4445/list -H "Content-Type: application/json" -d '{"GUN":"docker.io/dgeiger/nginx", "Tag":"1.15", "notaryServer":"notary-server-svc:4443"}' --cacert /etc/certs/notary/notary-server-svc:4443/root-ca.crt
```

### one image (with tag):
```
curl -X POST https://notary-wrapper-svc:4445/list -H "Content-Type: application/json" -d '{"GUN":"docker.io/library/nginx", "Tag":"1.17", "notaryServer":"notary.docker.io"}' --cacert /etc/certs/notary/notary-server-svc:4443/root-ca.crt
```
### all images (empty tag):
```
curl -X POST https://notary-wrapper-svc:4445/list -H "Content-Type: application/json" -d '{"GUN":"docker.io/library/alpine", "Tag":"", "notaryServer":"notary.docker.io"}' --cacert /etc/certs/notary/notary-server-svc:4443/root-ca.crt
```

## lookup
### notary
```
notary lookup docker.io/library/nginx 1.17 -s https://notary-server-svc:4443 
```
### curl
```
curl -X POST https://notary-wrapper-svc:4445/lookup -H "Content-Type: application/json" -d '{"GUN":"docker.io/library/nginx", "Tag":"1.17", "notaryServer":"notary.docker.io"}' --cacert /etc/certs/notary/notary-server-svc:4443/root-ca.crt
```

## verify
### curl
```
curl -X POST https://notary-wrapper-svc:4445/verify -H "Content-Type: application/json" -d '{"GUN":"docker.io/library/nginx", "SHA":"b2d89d0a210398b4d1120b3e3a7672c16a4ba09c2c4a0395f18b9f7999b768f2", "notaryServer":"notary.docker.io"}' --cacert /etc/certs/notary/notary-server-svc:4443/root-ca.crt
```

### test wrong image/tag
```
curl -X POST https://notary-wrapper-svc:4445/verify -H "Content-Type: application/json" -d '{"GUN":"docker.io/library/nginx", "SHA":"fail", "notaryServer":"notary.docker.io"}' --cacert /etc/certs/notary/notary-server-svc:4443/root-ca.crt
```

