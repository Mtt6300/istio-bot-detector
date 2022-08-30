# Istio-bot-detector
Istio-bot-detector is a wasm http filter which developed with golang, designed for Envoy to detect and reject incoming spam request's from bot.

This module works with `user-agent` header to detect bot's and will compare them with regular expressions provided by you. This filter **does not set any initial rules**, it means up to you to configure it.

# Configuring
This is an example of `WasmPlugin` configuration which you can apply filter to envoy in istio service mesh. More about [WasmPlugin](https://istio.io/latest/docs/reference/config/proxy_extensions/wasm-plugin/).

We are talking about `pluginConfig`, part of `WasmPlugin` which designed for plugin configured.

```yaml
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
metadata:
  name: bot-detector
  namespace: default
spec:
  selector:
    matchLabels:
      app: httpbin
      
  url: oci://docker.io/mtt6300/istio-bot-detector:latest
  pluginConfig:
    allow:
      - "Curl.*"
    deny:
      - "Scrapy.*"

    denyAll: false
    cacheSize: 200
```

`allow`: An array with all the [regex2](https://github.com/google/re2) that define bots. Matching bots are **allowed**.

`deny`: An array with all the [regex2](https://github.com/google/re2) that define bots. Matching bots are **denied**.

`denyAll`: This option give you ability to deny all user agents. It's useful when you want have a filter which accept incoming requests with **only** have matched `user-agent` with your allow expressions list. You can also ignore or remove deny section in your config when you use this option. More example.

*note*: default value is **false**.

`cacheSize`: A fixed integer which provide Size of cache which help you to speed up detection process. This is number of user agents you want to keep them on memory. More about [LRU caching system](https://github.com/hashicorp/golang-lru).

*note*: default value is **200**.

*note*: `cacheSize` can'not be lower than 1.

# Compile
My tested environment:

* `go` version `1.18.3`
* `tinygo` version `0.25.0`
* `LLVM` version `14.0.0`

Clone repository:

```bash
git clone https://github.com/Mtt6300/istio-bot-detector/
cd istio-bot-detector
```

Compile:

```bash
go mod tidy
tinygo build -o main.wasm -scheduler=none -target=wasi main.go
```

# Contributing , idea ,issue
Feel free to fill an issue or create a pull request, I'll check it ASAP
