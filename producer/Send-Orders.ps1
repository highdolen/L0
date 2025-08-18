param(
    [int]$Count = 5,
    [string]$Delay = "2s",
    [string]$Broker = "localhost:9093",
    [string]$Topic = "orders"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "        Kafka Order Producer" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Параметры:" -ForegroundColor Yellow
Write-Host "  Брокер: $Broker" -ForegroundColor White
Write-Host "  Топик: $Topic" -ForegroundColor White
Write-Host "  Количество заказов: $Count" -ForegroundColor White
Write-Host "  Задержка: $Delay" -ForegroundColor White
Write-Host ""

# Переходим в папку со скриптом
$ScriptPath = Split-Path -Parent $MyInvocation.MyCommand.Definition
Push-Location $ScriptPath

try {
    Write-Host "Отправляем заказы..." -ForegroundColor Green
    
    go run kafka_producer.go -broker="$Broker" -topic="$Topic" -count="$Count" -delay="$Delay"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "" 
        Write-Host "✅ Успешно отправлено $Count заказов!" -ForegroundColor Green
    } else {
        Write-Host "" 
        Write-Host "❌ Ошибка при отправке заказов" -ForegroundColor Red
    }
}
catch {
    Write-Host "❌ Ошибка: $($_.Exception.Message)" -ForegroundColor Red
}
finally {
    Pop-Location
}

Write-Host ""
Write-Host "Нажмите любую клавишу для закрытия..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
