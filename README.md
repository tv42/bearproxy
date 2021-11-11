# bearproxy -- Authorization enforcing HTTP reverse proxy

Bearproxy is a very simple HTTP reverse proxy that checks that requests contain a valid secret as a bearer token.

It is meant to work as a minimal sidecar container, protecting a web service that doesn't have any authorization logic of its own.

```
go build eagain.net/go/bearproxy
echo s3kr1t >secret
./bearproxy -secret-token-file=secret -backend-url=http://localhost:8001/ -listen=:8000
```
