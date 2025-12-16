# 测试脚本
# 用于运行 TR369 WebSocket 客户端的测试

param(
    [switch]$Coverage,
    [switch]$Verbose,
    [string]$Package = "./..."
)

Write-Host "开始运行测试..."

$args = @("test")

if ($Coverage) {
    $args += "-cover"
    $args += "-coverprofile=coverage.out"
}

if ($Verbose) {
    $args += "-v"
}

$args += $Package

# 执行测试
& go $args

if ($LASTEXITCODE -eq 0) {
    Write-Host "测试通过"
    
    if ($Coverage) {
        Write-Host "生成覆盖率报告..."
        go tool cover -html=coverage.out -o coverage.html
        Write-Host "覆盖率报告已生成: coverage.html"
    }
} else {
    Write-Host "测试失败"
    exit 1
}
