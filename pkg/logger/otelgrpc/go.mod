module github.com/kamalyes/go-toolbox/pkg/logger/otelgrpc

go 1.26.0

require (
	github.com/kamalyes/go-toolbox v0.19.4
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel/trace v1.47.0
	google.golang.org/grpc v1.85.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.opentelemetry.io/otel v1.47.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.48.0 // indirect
)

// 本仓库内的父模块：发布前需把上面的 require 版本提到含所需符号的 tag
// replace 仅对本模块的本地构建生效，下游主模块会忽略它
replace github.com/kamalyes/go-toolbox => ../../../
