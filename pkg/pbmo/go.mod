module github.com/kamalyes/go-toolbox/pkg/pbmo

go 1.26.0

require (
	github.com/kamalyes/go-toolbox v0.19.4
	github.com/stretchr/testify v1.12.1
	google.golang.org/protobuf v1.36.12
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

// 本仓库内的父模块：发布前需把上面的 require 版本提到含所需符号的 tag
// replace 仅对本模块的本地构建生效，下游主模块会忽略它
replace github.com/kamalyes/go-toolbox => ../../
