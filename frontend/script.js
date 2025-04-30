async function calculate() {
    const expr = document.getElementById('expression').value;
    const resultElement = document.getElementById('result');

    if (!expr.trim()) {
        resultElement.textContent = "Введите выражение";
        return;
    }

    try {
        resultElement.textContent = "Отправляем выражение...";

        const calcResponse = await fetch('http://localhost:4040/api/v1/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                expression: expr
            })
        });

        if (!calcResponse.ok) {
            const errorData = await calcResponse.json();
            throw new Error(errorData.error || `Ошибка HTTP: ${calcResponse.status}`);
        }

        const { id } = await calcResponse.json();
        resultElement.textContent = `Выражение принято. ID: ${id}. Ожидаем результат...`;

        const result = await pollExpressionResult(id);
        resultElement.textContent = `Результат: ${result}`;

    } catch (error) {
        resultElement.textContent = `Ошибка: ${error.message}`;
        console.error('Ошибка:', error);
    }
}

async function pollExpressionResult(id) {
    const maxAttempts = 30;
    const delay = 1000;

    for (let attempt = 0; attempt < maxAttempts; attempt++) {
        try {
            const response = await fetch(`http://localhost:4040/api/v1/expressions/${id}`);

            if (!response.ok) {
                throw new Error(`Ошибка HTTP: ${response.status}`);
            }

            const data = await response.json();

            if (data.status === "completed") {
                return data.result;
            } else if (data.status === "error") {
                throw new Error(data.error || "Ошибка при вычислении");
            }

        } catch (error) {
            console.error(`Попытка ${attempt + 1}:`, error);
        }

        await new Promise(resolve => setTimeout(resolve, delay));
    }

    throw new Error("Превышено время ожидания результата");
}
document.getElementById('expression').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') calculate();
});