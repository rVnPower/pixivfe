# About tracing in PixivFE

Every request to pixiv websites should go through core/requests.go.

Every request to pixiv websites is traced.
Every server request should be traced, but currently not.

## Todo

- [x] Trace asset route
- [x] Trace every server route

## How to use tracing

```
wget https://repo1.maven.org/maven2/io/zipkin/zipkin-server/3.4.1/zipkin-server-3.4.1-exec.jar
java -jar zipkin-server-3.4.1-exec.jar
# start pixivfe
```

That's it!

To see the spans, open http://localhost:9411/ and click "RUN QUERY".
