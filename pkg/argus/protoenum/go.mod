module github.com/kamalyes/go-toolbox/pkg/argus/protoenum

go 1.25.0

require (
	github.com/kamalyes/go-toolbox v0.19.4
	github.com/stretchr/testify v1.12.1
	google.golang.org/protobuf v1.36.12
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

// replace 仅对本模块的本地构建生效，下游主模块会忽略它
replace github.com/kamalyes/go-toolbox => ../../../
