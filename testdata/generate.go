// Requires https://github.com/jteeuwen/go-bindata
//
// Install with go get -u github.com/jteeuwen/go-bindata/...

//go:generate go-bindata -nocompress -o=generated.go -pkg=testdata ./testdata.txt.gz
//go:generate gofmt -w .

package testdata
