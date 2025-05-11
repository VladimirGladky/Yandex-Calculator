# Распределённый вычислитель арифметических выражений
## О проекте
Многопользовательский режим , теперь вычислитель работает в контексте конкртеного пользователя и все данные хранятся в СУБД .
Вычислитель состоит из двух элементов: оркестратор и агент , которые взаимодействуют друг с другом через gRPC.

### Регистрация

Пользователь отправляет **POST-запрос** для регистрации:  
`POST /api/v1/register`

### 📝 Параметры запроса (JSON)

| Параметр  | Тип    | Описание |
|-----------|--------|----------|
| `login`   | string | Логин пользователя. <br>• Длина: **3–20 символов**<br>• Допустимые символы: **A-Z**, **a-z**, **0-9**<br>• **Пробелы запрещены** |
| `password`| string | Пароль пользователя. <br>• Длина: **8–64 символа**<br>• **Пробелы запрещены** |


#### Пример запроса (cURL):
```bash
curl --location 'http://localhost:4040/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
    "login": "test10",
    "password": "testify!"
}'
```
Успешный ответ:
```
{
 "status": "success"
}
```


### Дополнительные детали:
- Сервер сохраняет **хэшированный пароль** в СУБД.
- Логин проверяется на уникальность перед регистрацией.


## 📥 Вход в систему

Для авторизации пользователь отправляет **POST-запрос** на получение JWT-токена.

**Метод:**  
`POST /api/v1/login`

### 📝 Параметры запроса (JSON)

| Параметр  | Тип    | Описание |
|-----------|--------|----------|
| `login`   | string | Логин пользователя. <br>• Длина: **3–20 символов**<br>• Допустимые символы: **A-Z**, **a-z**, **0-9**<br>• **Пробелы запрещены** |
| `password`| string | Пароль пользователя. <br>• Длина: **8–64 символа**<br>• **Пробелы запрещены** |

### 💡 Пример запроса (cURL)

```bash
curl --location 'http://localhost:4040/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "login": "test10",
  "password": "testify!"
}'
```

Успешный ответ:
```
{
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QxMCIsImV4cCI6MTc0NzA3Mjk2NSwiaWF0IjoxNzQ2OTg2NTY1LCJuYmYiOjE3NDY5ODY1NjV9.6mD-1dsO73imjQvIMgV8LQBi1zgj4Yybt3ltrIsf58U"
}
```

📌 Примечание:
Полученный токен действует 24 часа. По истечении этого времени необходимо повторно пройти авторизацию.
Во всех последующих запросах токен передаётся в заголовке Authorization в формате:
Authorization: Bearer <token>


### Оркестратор
Оркестратор принимает на вход арифметическое выражение и переводит её в набор последовательных задач и обеспечивает порядок их выполнения. Это производится с помощью дерева(ast). У выражения есть своя структура Expression , в которую включены поля для хранения последовательных задач и узлов дерева. Также оркестратор хранит все выражения , которые вы ему отправляете.Также оркестратор создаёт gRPC сервер для взаимодействия с агентом , который будет выступать в роли клиента  
У оркестратора есть 3 endpoint-а:
```
1) curl --location 'localhost:4040/api/v1/calculate' \
--header 'Authorization: Bearer <token>' \
--header 'Content-Type: application/json' \
--data '{
"expression": <строка с выражение>
}' 
```
С помощью данного эндпоинта вы отправляете выражение оркестратору , которое он обрабатывает и сохраняет себе
```
2) curl --location 'localhost:4040/api/v1/expressions'
--header 'Authorization: Bearer <token>'    
```
Этот эндпоинт можно использовать для получения всего списка выражений
```
3) curl --location 'localhost:4040/api/v1/expressions/:id'
--header 'Authorization: Bearer <token>'   
```
Этот можно использовать для получения конкретного выражения по id , который вы добавили ранее

Примеры использования эндпоинтов будут дальше

### Агент
Взаимодействие оркестратора и агента реализовано с помощью gRPC.Агент создаёт клиент для взаимодействия с gRPC сервером оркестратора. Агент запускает несколько горутин , которые задаются переменной среды COMPUTING_POWER, каждая горутина выступает в роли независемого вычислителя.
Каждую секунду агент с помощью gRPC клиента отправляет запросы оркестратору на TaskGet. Если он получает задачу , то вычисляет её и отправляет обратно оркестратору с помощью TaskPost. Если оркестратор отвечает, что в данный момент нет задачи для выполнения , то агент делает паузу на 15 секунд.

## Схема работы

Схема иллюстрирует то как работает агент и оркестратор взаимодействуют между собой, а также процесс регистрации перед началом работы с вычислителем

![schemework.png](pkg/assets/image/schemework.png)

Также вот так разбивается выражение 2+2×2+2 в дерево в оркестраторе:

![expression.png](pkg/assets/image/expression.png)

## Технологии и библиотеки
Вычислитель написан на языке **Go** и использует следующие библиотеки и инструменты:

#### Язык программирования:
- **Go** (версия 1.23.2)

#### Библиотеки:
- **[github.com/google/uuid](https://github.com/google/uuid)**: Для генерации уникальных идентификаторов (UUID).
- **[github.com/gorilla/mux](https://github.com/gorilla/mux)**: Для создания HTTP-роутера и обработки запросов.
- **[github.com/joho/godotenv](https://github.com/joho/godotenv)**: Для загрузки переменных окружения из файла `.env`.
- **[github.com/stretchr/testify](https://github.com/stretchr/testify)**: Для написания unit-тестов.
- **[go.uber.org/zap](https://go.uber.org/zap)**: Для структурированного и производительного логирования.
- **[github.com/go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation)**: Для валидации данных (валидаторы для структур, строк и других типов).
- **[github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)**: Для создания, подписания и проверки JWT-токенов.
- **[github.com/golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate)**: Для управления миграциями базы данных.
- **[github.com/ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv)**: Для простого чтения конфигурации из переменных окружения и файлов.
- **[github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)**: Драйвер для работы с базой данных SQLite.
- **[golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)**: Набор дополнительных криптографических пакетов (например, bcrypt, hash-функции и др.).
- **[google.golang.org/grpc](https://github.com/grpc/grpc-go)**: Для создания gRPC-сервисов (высокопроизводительная RPC-коммуникация).
- **[google.golang.org/protobuf](https://github.com/protocolbuffers/protobuf-go)**: Для работы с Protocol Buffers (сериализация и десериализация данных).



## Структура проекта 

- **`cmd/`**
  - **`agent/`**
  - **`orchestrator/`**
  - **`migrator/`**
-**`config/`**
-**`frontend/`**
- **`gen/`**
- **`internal/`**
  - **`agent/`**
  - **`app/`**
  - **`config/`**
  - **`domain/`**
  - **`orchestrator/`**
    - **`parser/`**
    - **`service/`**
    - **`storage/`**
    - **`transport/`**
- **`migrations/`**
- **`pkg/`**
  - **`assets/`**
  - **`calculation/`**
  - **`jwt/`**
  - **`logger/`**
- **`storage/`**
- **`proto/`**
- **`tests`**

## Quick start

Для начала нужно склонировать репозиторий командой 
```
git clone https://github.com/VladimirGladky/FinalTaskFirstSprint.git
```

После этого вам нужно перейти в папку с проектом 
```
cd FinalTaskFirstSprint
```

Теперь вы можете запустить проект , но для этогт нужно чтобы был установлен Go версии 1.23.2
Ссылка для скачивания: [Go Download](https://go.dev/doc/install)

Перед запуском агента и оркестратора , воспользуйтесь командой

```bash
go mod download

```

#### Установите переменные среды :

Переменная среды `SECRET` используется для шифрования JWT-токенов, поменяйте ее на свою секретную фразу.
Также вы можете создать файл `local.env` в корневой папке проекта и добавить переменные среды в него.

для Linux/macOS:

```bash
export TIME_ADDITION_MS=200
export TIME_SUBTRACTION_MS=200
export TIME_MULTIPLICATIONS_MS=300
export TIME_DIVISIONS_MS=400
export COMPUTING_POWER=4
export SECRET=secret
```

для Windows:

```bash
set TIME_ADDITION_MS=200
set TIME_SUBTRACTION_MS=200
set TIME_MULTIPLICATIONS_MS=300
set TIME_DIVISIONS_MS=400
set COMPUTING_POWER=4
set SECRET=secret
```


Сначала запускается оркестратор , затем запускается агент

```
go run ./cmd/orchestrator/main.go
```

потом открываете другой терминал и запускаете агента

```
go run ./cmd/agent/main.go
```

Для прекращения работы агента или оркестратора можете нажать сочетание клавиш Ctrl+C

## Примеры использования со всеми возможными сценариями

После запуска проекта вы можете отправлять cURL-запросы к сервису:

Так как в терминале Windows не обрабатываются cURL запросы я использовал git bash.

Нужно отметить , что мой веб-сервис использует порт 4040(надеюсь он у вас не занят 😊)

### Запросы по endpoint-у:

```bash
'127.0.0.1:4040/api/v1/calculate' 
```

cURL команда с ответом сервиса 201:

```bash
 curl --location '127.0.0.1:4040/api/v1/calculate' --header 'Authorization: Bearer <token>' --header  'Content-Type: server/json' --data '{
  "expression": "2+2*2"
}'
```
Ответ:

```bash
{"id":"f67287f7-f29f-4196-b3be-7abff3bec739"}
```

cURL команда с ответом сервиса 400:
```bash
curl --location '127.0.0.1:4040/api/v1/calculate' --header 'Authorization: Bearer <token>' --header 'Content-Type: server/json' --data '{
  "expression": "2+2*2
}'
```
Ответ:
```bash
{"error":"Bad request"}
```

cURL команда с ответом сервиса 405:
```bash
curl --request GET \ --url '127.0.0.1:4040/api/v1/calculate' --header 'Authorization: Bearer <token>' --header 'Content-Type: server/json' --data '{
  "expression": "2+2*2"
}'
```
Ответ:
```bash
{"error":"You can use only POST method"}
```

cURL команда с ответом сервиса 422:
```bash
curl --location '127.0.0.1:4040/api/v1/calculate' --header 'Authorization: Bearer <token>' --header 'Content-Type: server/json' --data '{
  "expression": "2+2*2)"
}'
```
Ответ:
```bash
{"error":"Expression is not valid"}
```

### Запросы по endpoint-у:

```bash
'127.0.0.1:4040/api/v1/expressions' 
```

cURL команда с ответом сервиса 200:
```bash
curl --location '127.0.0.1:4040/api/v1/expressions' --header 'Authorization: Bearer <token>'
```
Ответ:
```bash
{
  "expressions":
    [
      {
        "id":"56e8677e-a058-485d-bc2c-342af7130c4c",
        "status":"done",
        "result":6
      },
      {
        "id":"7ddf6a72-d2a3-4bcc-947c-998fd9eac383",
        "status":"done",
        "result":10
      }
    ]
}
```

### Запросы по endpoint-у:

```bash
'127.0.0.1:4040/api/v1/expressions/{id}' 
```

cURL команда с ответом сервиса 200:
```bash
curl --location '127.0.0.1:4040/api/v1/expressions/56e8677e-a058-485d-bc2c-342af7130c4c' --header 'Authorization: Bearer <token>'
```
Ответ:
```bash
  {
    "id":"56e8677e-a058-485d-bc2c-342af7130c4c",
    "status":"done",
    "result":6}
```

cURL команда с ответом сервиса 404:
```bash
curl --location '127.0.0.1:4040/api/v1/expressions/56e8677e-a058-485d-bc2c-342af7130c4'
```

Ответ:
```bash
{"error":"Expression not found"}
```

### Запросы по endpoint-у:

метод GET

```bash
'127.0.0.1:4040/internal/taskGet' 
```

cURL команда с ответом сервиса 200:
```bash
curl --location '127.0.0.1:4040/internal/taskGet'
```

Ответ:
```bash
{
    "id":"5",
    "arg1":4,
    "arg2":2,
    "operation":"*",
    "operation_time":300
}
```

curl команда с ответом сервиса 404:

```bash
curl --location '127.0.0.1:4040/internal/taskGet'
```

Ответ:
```bash
{"error":"No taskGet available"}
```

Мой телеграм https://t.me/smoothhhhhhh





