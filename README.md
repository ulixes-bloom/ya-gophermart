# Gophermart

Индивидуальный дипломный проект курса «Go-разработчик»

## Сводное HTTP API
Накопительная система лояльности «Гофермарт» предоставляет следующие HTTP-хендлеры:

* `POST /api/user/register` — регистрация пользователя;
* `POST /api/user/login` — аутентификация пользователя;
* `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта;
* `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
* `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя;
* `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
* `GET /api/user/withdrawals` — получение информации о выводе средств с накопительного счёта пользователем.

### Регистрация пользоателя

`POST /api/user/register`

##### Параметры запроса

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| user | body | User Registration Information | Yes | [User](#User) |

##### Коды ответа

| Code | Description |
| ---- | ----------- |
| 200 | пользователь успешно зарегистрирован и аутентифицирован |
| 400 | неверный формат запроса |
| 409 | логин уже занят |
| 500 | внутренняя ошибка сервера |

### Аутентификация пользователя 

`POST /api/user/login`

##### Параметры запроса

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| user | body | User Registration Information | Yes | [User](#User) |

##### Коды ответа

| Code | Description |
| ---- | ----------- |
| 200 | пользователь успешно аутентифицирован |
| 400 | неверный формат запроса |
| 401 | неверная пара логин/пароль |
| 500 | внутренняя ошибка сервера |

### Регистрация нового заказа

`GET /api/user/orders`

##### Параметры запроса

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | Bearer | No | string |

##### Коды ответа

| Code | Description |
| ---- | ----------- |
| 200 | успешная обработка запроса |
| 204 | нет данных для ответа |
| 401 | пользователь не авторизован |
| 500 | внутренняя ошибка сервера |

### Получение списка загруженных номеров заказов

`POST /api/user/orders`

##### Параметры запроса

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | Bearer | No | string |
| user | body | User Registration Information | Yes | [OrderRequest](#OrderRequest) |

##### Коды ответа

| Code | Description |
| ---- | ----------- |
| 200 | номер заказа уже был загружен этим пользователем |
| 202 | новый номер заказа принят в обработку |
| 400 | неверный формат запроса |
| 401 | пользователь не аутентифицирован |
| 409 | номер заказа уже был загружен другим пользователем |
| 422 | неверный формат номера заказа |
| 500 | внутренняя ошибка сервера |

### Получение текущего баланса пользователя

`GET /api/user/balance`

##### Параметры запроса

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | Bearer | No | string |

##### Коды ответа

| Code | Description |
| ---- | ----------- |
| 200 | успешная обработка запроса |
| 401 | пользователь не авторизован |
| 500 | внутренняя ошибка сервера |

### Запрос на списание средств

`POST /api/user/balance/withdraw`

##### Параметры запроса

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | Bearer | No | string |
| user | body | User Registration Information | Yes | [WithdrawalRequest](#WithdrawalRequest) |

##### Коды ответа

| Code | Description |
| ---- | ----------- |
| 200 | успешная обработка запроса |
| 401 | пользователь не авторизован |
| 402 | на счету недостаточно средств |
| 422 | неверный номер заказа |
| 500 | внутренняя ошибка сервера |

### Получение информации о выводе средств
`GET /api/user/withdrawals`

##### Параметры запроса

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | Bearer | No | string |

##### Коды ответа

| Code | Description |
| ---- | ----------- |
| 200 | успешная обработка запроса |
| 204 | нет ни одного списания |
| 401 | пользователь не авторизован |
| 500 | внутренняя ошибка сервера |

### Модели данных

#### OrderRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| number | string |  | No |

#### User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| login | string |  | No |
| password | string |  | No |

#### WithdrawalRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| order_number | string |  | No |
| sum | number |  | No |

## Эксплуатация

## Запуск
```
go build -o gophermart ./cmd/gophermart/main.go
./gophermart -d "host=localhost user=<LOGIN> password=<PASSWORD> sslmode=disable" -k "<SECRET_KEY>" -o 5s -l Debug -rl 3
```

| Key | Type | Description | Default value |
| --- | ---- | ----------- | ------------- |
| -a | string | address and port to run service | ":8080" |
| -d | string | database connection string | "" |
| -k | string | secret key to generate jwt token | "SECRET_KEY" |
| -l | string | database connection string | "Info" |
| -o | time.duration | order info update interval | 30s |
| -r | string | accrual system address | "localhost:8081" |
| -rl | int | accrual rate limit | 2 |
| -t | time.duration | jwt token lifetime | 8h |

### Тестирование

```
go test -v ./...
```