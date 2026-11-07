version=0.8.0

GOOS=windows GOARCH=amd64 go build -o xo-windows main.go
go build -o xo-mac-arm64 main.go
GOOS=darwin GOARCH=amd64 go build -o xo-mac-amd64 main.go
GOOS=linux GOARCH=amd64 go build -o xo-linux-amd64 main.go

zip xo-$version-windows-amd64.zip xo-windows
zip xo-$version-mac-arm64.zip xo-mac-arm64
zip xo-$version-mac-amd64.zip xo-mac-amd64
zip xo-$version-linux-amd64.zip xo-linux-amd64


rm xo-windows
rm xo-mac-arm64
rm xo-mac-amd64
rm xo-linux-amd64


echo "done packing!"
