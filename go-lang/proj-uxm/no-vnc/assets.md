
This is deprecated and currentl;y unused as of go 1.16 is using native go embed FS

# @go-libs/go/src/github.com/jteeuwen/go-bindata/go-bindata/
# go build -o go-bindata main.go version.go AppendSliceValue.go

go-bindata -o assets.go -prefix assets assets/...

