# Jaeger go example

```shell
# run jaeger
docker run -d -p 5775:5775/udp -p 16686:16686 -p 14250:14250 -p 14268:14268 jaegertracing/all-in-one:1.48.0
```

参考: [go opentracing tutorial](https://github.com/yurishkuro/opentracing-tutorial/tree/master/go)
