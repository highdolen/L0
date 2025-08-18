document.getElementById("order-form").addEventListener("submit", async (e) => {
    e.preventDefault();

    const orderId = document.getElementById("order-id").value.trim();
    const resultDiv = document.getElementById("result");

    resultDiv.textContent = "Загрузка...";

    try {
        const response = await fetch(`/order/${orderId}`);
        if (!response.ok) {
            resultDiv.textContent = `Ошибка: ${response.status} ${response.statusText}`;
            return;
        }
        const data = await response.json();
        resultDiv.textContent = JSON.stringify(data, null, 2);
    } catch (err) {
        resultDiv.textContent = `Ошибка запроса: ${err.message}`;
    }
});
