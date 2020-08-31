# Goello

> one-line mDNS service registration

## What?

Dead simple one-line cli to register an mDNS service, based on the awesome brutella library.

## TL;DR

Call the binary from anywhere:
```
make
dist/goello-server -type _http._tcp -name "Fluffy!" -port 1234
```

You can optionally pass along an explicit hostname with `-host "hostname-something"` and/or ip with `-ip 1.2.3.4`.

This is especially useful as a side-car docker container.

```
docker run -d \
    --name goello-sidecar-for-fluffy-container \
    --read-only \
    --cap-drop ALL \
    --net host \
    --rm \
    dubodubonduponey/goello -type _http._tcp -name "Fluffy service!" -host "hostname-for-fluffy" -port 1234 -ip "$(docker inspect fluffy-container | jq -rc '.[0].NetworkSettings.IPAddress')"
```

## Where is the Dockerfile?

https://github.com/dubo-dubon-duponey/docker-goello
