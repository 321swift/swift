# windows
# 32-bit
GOOS=windows GOARCH=386 go  build -vo bin/windows/swift-386.exe main.go

# 64-bit
GOOS=windows GOARCH=amd64 go  build -vo bin/windows/swift-amd64.exe main.go


# macos
# 64-bit
GOOS=darwin GOARCH=amd64 go  build -vo bin/macos/swift-amd64-darwin main.go

# 32-bit
# GOOS=darwin GOARCH=386 go  build -vo bin/macos/swift-386-darwin main.go


# linux
# 64-bit
GOOS=linux GOARCH=amd64 go  build -vo bin/linux/swift-amd64-linux main.go

# 32-bit
GOOS=linux GOARCH=386 go  build -vo bin/linux/swift-386-linux main.go
