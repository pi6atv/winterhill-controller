# template
Replace:

* `drx` -> package name
* `9001` -> port number

Create web:
```bash
vue create web
cd web
vue add vuetify
```

## Installing
This package expects the file `/etc/nginx/snippets/cert.conf` to exist,
with the certificate(s).

See the releases in the repo for the debian packages.

## Building
### Build for x86
```bash
./build-with-docker.sh
```

### Build for armhf
```bash
GOARCH=arm GOARM=7 GOOS=linux CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc ARCH=armhf ./build-with-docker.sh 
```
