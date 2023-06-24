# windows
# 32-bit
GOOS=windows GOARCH=386 go  build  -o dist/windows/swift-386.exe main.go
7z a -tzip ./dist/windows/swift-386-windows.zip ./dist/windows/swift-386.exe ./ui

# 64-bit
GOOS=windows GOARCH=amd64 go  build  -o dist/windows/swift-amd64.exe main.go
7z a -tzip ./dist/windows/swift-amd64-windows.zip ./dist/windows/swift-amd64.exe ./ui

# macos
# 64-bit
GOOS=darwin GOARCH=amd64 go  build  -o dist/macos/swift-amd64-darwin main.go
7z a -tzip ./dist/macos/swift-amd64-darwin.zip ./dist/macos/swift-amd64-darwin ./ui

# # linux
# # 64-bit
GOOS=linux GOARCH=amd64 go  build  -o dist/linux/swift-amd64-linux main.go
7z a -tzip ./dist/linux/swift-amd64-linux.zip ./dist/linux/swift-amd64-linux ./ui

# # 32-bit
GOOS=linux GOARCH=386 go  build  -o dist/linux/swift-386-linux main.go
7z a -tzip ./dist/linux/swift-386-linux.zip ./dist/linux/swift-386-linux ./ui

rm ./dist/linux/swift-386-linux
rm ./dist/linux/swift-amd64-linux
rm ./dist/macos/swift-amd64-darwin
rm ./dist/windows/swift-386.exe
rm ./dist/windows/swift-amd64.exe
