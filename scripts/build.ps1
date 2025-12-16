# 构建脚本
# 用于构建 TR369 WebSocket 客户端

param(
    [string]$Output = "bin/tr369-wss-client.exe",
    [string]$GOOS = "windows",
    [string]$GOARCH = "amd64"
)

Write-Host "开始构建 TR369 WebSocket 客户端..."

# 设置环境变量
$env:GOOS = $GOOS
$env:GOARCH = $GOARCH

# 执行构建
go build -o $Output ./cmd/client

if ($LASTEXITCODE -eq 0) {
    Write-Host "构建成功: $Output"
} else {
    Write-Host "构建失败"
    exit 1
}
