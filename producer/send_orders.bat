@echo off
echo ========================================
echo        Kafka Order Producer
echo ========================================
echo.

set /p COUNT="Введите количество заказов для отправки (по умолчанию 5): "
if "%COUNT%"=="" set COUNT=5

set /p DELAY="Введите задержку между сообщениями в секундах (по умолчанию 2): "
if "%DELAY%"=="" set DELAY=2s

echo.
echo Отправляю %COUNT% заказов с задержкой %DELAY%...
echo.

cd /d "%~dp0"
go run kafka_producer.go -count=%COUNT% -delay=%DELAY%

echo.
echo Готово! Нажмите любую клавишу для закрытия...
pause >nul
