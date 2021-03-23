### Ruby client auth server

Usage:

```
pack build appimage -b gcr.io/paketo-buildpacks/ca-certificates -b gcr.io/paketo-buildpacks/ruby
```

```
docker run --init -it -e SERVICE_BINDING_ROOT=/bindings -v "$PWD"/binding:/bindings/ca-certificates -e PORT=8080 -p 8080:8080 appimage
```

```
curl --cert cert.pem --key key.pem --cacert ca.pem https://localhost:8443
Hello, world!
```
