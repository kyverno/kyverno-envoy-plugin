# Configuration

## Nested Request

One of the flags you can pass to the HTTP authz server is the `--nestedRequest` boolean parameter, which controls a key behavior of the authz server. And that is where it finds the request to authenticate in the incoming request its receiving. If `false`, the authz server will run the policies against the request it receives directly. But if `true`, it will expect that the request its receiving contains in its body the bytes of the request it should authenticate. i.e `curl -XPOST https://authz-server -d {<THE-FULL-BYTES-OF-THE-REQUEST-TO-AUTHENTICATE>}`.

The benefits this offers is ease of development for devs writing their servers that delegate their authentication to the Authz server. As they can simply create requests to send to the Authz server without any transformation, headers.. etc.

For example, this is how a golang server would structure the request when `nestedRequest` is true versus when its false:

```golang
// read the request obtained from the client
rawBytes, err := httputil.DumpRequest(r, true)
if err != nil {
  // handle error
}

// put the request in the body of the request to send to the authz server
req, err := http.NewRequest(http.MethodPost, "https://<authz-server-endpoint>", io.NopCloser(bytes.NewReader(rawBytes)))
```

Otherwise, the code would have to do some transformations to convey information (the original host, method, query.. etc) about the original request to the Authz server
