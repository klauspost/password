// Requires https://github.com/jteeuwen/go-bindata
//
// Install with go get -u github.com/jteeuwen/go-bindata/...

//go:generate go-bindata -nocompress -o=generated.go -pkg=testdata -ignore=generate .

package testdata
