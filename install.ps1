# Windows 安装脚本
# 使用方法: powershell -ExecutionPolicy Bypass -File install.ps1

Write-Host "正在安装 CliGool Windows 客户端..."

$SourcePath = "\\wsl$\Ubuntu\home\zhanju\projects\cligool\bin\cligool-windows-amd64.exe"
$DestPath = "$env:USERPROFILE\cligool-windows-amd64.exe"

# 复制文件
Copy-Item $SourcePath $DestPath -Force

Write-Host "✅ 已安装到: $DestPath"
Write-Host ""
Write-Host "使用方法:"
Write-Host "  .\cligool-windows-amd64.exe -cmd claude"
Write-Host "  .\cligool-windows-amd64.exe -cmd gemini"
Write-Host ""
Write-Host "或者添加到 PATH 后可以在任何目录使用"
