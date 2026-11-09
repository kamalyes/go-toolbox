module github.com/kamalyes/go-toolbox/pkg/imgix

go 1.26.0

require (
	github.com/disintegration/imaging v1.6.2
	github.com/fogleman/gg v1.3.0
	github.com/kamalyes/go-toolbox v0.19.4
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
)

// 本仓库内的父模块：发布前需把上面的 require 版本提到含所需符号的 tag
// replace 仅对本模块的本地构建生效，下游主模块会忽略它
replace github.com/kamalyes/go-toolbox => ../../
