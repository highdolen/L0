@echo off
echo ==========================================
echo     Демонстрация Kafka Order Producer
echo ==========================================
echo.
echo Сейчас мы отправим несколько тестовых заказов
echo в Kafka и покажем их обработку системой.
echo.
pause

echo.
echo 1. Отправляем 1 заказ...
go run kafka_producer.go -count=1 -delay=0s
echo.

echo Подождем немного для обработки...
timeout /t 3 /nobreak >nul

echo.
echo 2. Отправляем 3 заказа с задержкой 1 секунда...
go run kafka_producer.go -count=3 -delay=1s
echo.

echo Подождем для обработки...
timeout /t 5 /nobreak >nul

echo.
echo 3. Отправляем 5 заказов быстро (без задержки)...
go run kafka_producer.go -count=5 -delay=0s
echo.

echo.
echo ==========================================
echo            ДЕМОНСТРАЦИЯ ЗАВЕРШЕНА
echo ==========================================
echo.
echo Всего было отправлено 9 заказов!
echo.
echo Теперь вы можете:
echo - Открыть веб-интерфейс: http://localhost:8080/
echo - Проверить Kafka UI: http://localhost:8081/
echo - Проверить логи: docker-compose logs app
echo.
pause
