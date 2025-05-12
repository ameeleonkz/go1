# API Documentation

## Публичные эндпоинты

### Регистрация пользователя
- **URL**: `/register`
- **Method**: `POST`
- **Request Body**:
```json
{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "securepassword123"
}
```
- **Response**: `200 OK`
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "email": "user@example.com",
        "username": "johndoe",
        "created_at": "2024-03-20T10:00:00Z",
        "updated_at": "2024-03-20T10:00:00Z"
    }
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат запроса
  - `409 Conflict` - Email или username уже существует

### Авторизация
- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:
```json
{
    "email": "user@example.com",
    "password": "securepassword123"
}
```
- **Response**: `200 OK`
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "email": "user@example.com",
        "username": "johndoe",
        "created_at": "2024-03-20T10:00:00Z",
        "updated_at": "2024-03-20T10:00:00Z"
    }
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат запроса
  - `401 Unauthorized` - Неверные учетные данные

## Защищенные эндпоинты (требуют JWT токен)

### Профиль пользователя
- **URL**: `/profile`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
{
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Пользователь не найден

### Счета

#### Создание счета
- **URL**: `/accounts`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
```json
{
    "currency": "RUB"
}
```
- **Response**: `201 Created`
```json
{
    "id": 1,
    "user_id": 1,
    "number": "40702810123456789012",
    "balance": 0,
    "currency": "RUB",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат запроса
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `409 Conflict` - Счет с такой валютой уже существует

#### Получение списка счетов
- **URL**: `/accounts`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
[
    {
        "id": 1,
        "user_id": 1,
        "number": "40702810123456789012",
        "balance": 1000.50,
        "currency": "RUB",
        "created_at": "2024-03-20T10:00:00Z",
        "updated_at": "2024-03-20T10:00:00Z"
    }
]
```
- **Errors**:
  - `401 Unauthorized` - Отсутствует или неверный токен

#### Получение счета по ID
- **URL**: `/accounts/{id}`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
{
    "id": 1,
    "user_id": 1,
    "number": "40702810123456789012",
    "balance": 1000.50,
    "currency": "RUB",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат ID
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Счет не найден

#### Перевод средств
- **URL**: `/transfer`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
```json
{
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 100.50
}
```
- **Response**: `200 OK`
```json
{
    "id": 1,
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 100.50,
    "type": "transfer",
    "status": "completed",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат запроса
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Счет не найден
  - `409 Conflict` - Недостаточно средств

#### Получение истории транзакций
- **URL**: `/transactions`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
[
    {
        "id": 1,
        "from_account_id": 1,
        "to_account_id": 2,
        "amount": 100.50,
        "type": "transfer",
        "status": "completed",
        "created_at": "2024-03-20T10:00:00Z",
        "updated_at": "2024-03-20T10:00:00Z"
    }
]
```
- **Errors**:
  - `401 Unauthorized` - Отсутствует или неверный токен

### Карты

#### Создание карты
- **URL**: `/cards`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
```json
{
    "account_id": 1,
    "cardholder_name": "JOHN DOE"
}
```
- **Response**: `201 Created`
```json
{
    "id": 1,
    "account_id": 1,
    "number": "4111111111111111",
    "expiry_date": "12/25",
    "cardholder_name": "JOHN DOE",
    "is_active": true,
    "created_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат запроса
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Счет не найден

#### Получение карты по ID
- **URL**: `/cards/{id}`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
{
    "id": 1,
    "account_id": 1,
    "number": "4111111111111111",
    "expiry_date": "12/25",
    "cardholder_name": "JOHN DOE",
    "is_active": true,
    "created_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат ID
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Карта не найдена

#### Получение списка карт по ID счета
- **URL**: `/cards`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Query Parameters**: `account_id=1`
- **Response**: `200 OK`
```json
[
    {
        "id": 1,
        "account_id": 1,
        "number": "4111111111111111",
        "expiry_date": "12/25",
        "cardholder_name": "JOHN DOE",
        "is_active": true,
        "created_at": "2024-03-20T10:00:00Z"
    }
]
```
- **Errors**:
  - `400 Bad Request` - Неверный формат account_id
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Счет не найден

#### Деактивация карты
- **URL**: `/cards/{id}/deactivate`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
{
    "id": 1,
    "is_active": false,
    "updated_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат ID
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Карта не найдена

### Кредиты

#### Создание кредита
- **URL**: `/credits`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
```json
{
    "account_id": 1,
    "amount": 100000,
    "term_months": 12
}
```
- **Response**: `201 Created`
```json
{
    "id": 1,
    "user_id": 1,
    "account_id": 1,
    "amount": 100000,
    "interest_rate": 12.5,
    "term_months": 12,
    "status": "active",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат запроса
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Счет не найден
  - `409 Conflict` - Недостаточно средств для погашения

#### Получение кредита по ID
- **URL**: `/credits/{id}`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
{
    "id": 1,
    "user_id": 1,
    "account_id": 1,
    "amount": 100000,
    "interest_rate": 12.5,
    "term_months": 12,
    "status": "active",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат ID
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Кредит не найден

#### Получение списка кредитов
- **URL**: `/credits`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
[
    {
        "id": 1,
        "user_id": 1,
        "account_id": 1,
        "amount": 100000,
        "interest_rate": 12.5,
        "term_months": 12,
        "status": "active",
        "created_at": "2024-03-20T10:00:00Z",
        "updated_at": "2024-03-20T10:00:00Z"
    }
]
```
- **Errors**:
  - `401 Unauthorized` - Отсутствует или неверный токен

#### Получение графика платежей
- **URL**: `/credits/{id}/schedule`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK`
```json
[
    {
        "id": 1,
        "credit_id": 1,
        "payment_number": 1,
        "amount": 9166.67,
        "principal": 8333.33,
        "interest": 833.34,
        "due_date": "2024-04-20T00:00:00Z",
        "status": "pending",
        "created_at": "2024-03-20T10:00:00Z",
        "updated_at": "2024-03-20T10:00:00Z"
    }
]
```
- **Errors**:
  - `400 Bad Request` - Неверный формат ID
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Кредит не найден

#### Внесение платежа
- **URL**: `/credits/{id}/payments/{payment_id}`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
```json
{
    "amount": 9166.67
}
```
- **Response**: `200 OK`
```json
{
    "id": 1,
    "credit_id": 1,
    "payment_number": 1,
    "amount": 9166.67,
    "principal": 8333.33,
    "interest": 833.34,
    "due_date": "2024-04-20T00:00:00Z",
    "status": "paid",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат запроса
  - `401 Unauthorized` - Отсутствует или неверный токен
  - `404 Not Found` - Кредит или платеж не найден
  - `409 Conflict` - Недостаточно средств для платежа

### Аналитика

#### Получение аналитики
- **URL**: `/analytics`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Query Parameters**: 
  - `forecast_days` (опционально) - количество дней для прогноза (по умолчанию 30)
- **Response**: `200 OK`
```json
{
    "monthly_stats": {
        "month": "2024-03-01T00:00:00Z",
        "total_income": 100000,
        "total_expenses": 50000,
        "net_income": 50000,
        "transactions_count": 10,
        "top_categories": [
            {
                "category": "income",
                "amount": 100000,
                "count": 5
            },
            {
                "category": "expense",
                "amount": 50000,
                "count": 5
            }
        ]
    },
    "credit_load": {
        "total_credits": 500000,
        "active_credits": 2,
        "monthly_payments": 25000,
        "total_debt": 500000,
        "credit_utilization": 100,
        "payment_to_income": 25
    },
    "balance_forecast": {
        "start_date": "2024-03-15T00:00:00Z",
        "end_date": "2024-04-14T00:00:00Z",
        "initial_balance": 100000,
        "forecast_days": [
            {
                "date": "2024-03-15T00:00:00Z",
                "expected_balance": 100000,
                "planned_income": 0,
                "planned_expenses": 0
            }
        ]
    }
}
```
- **Errors**:
  - `400 Bad Request` - Неверный формат forecast_days
  - `401 Unauthorized` - Отсутствует или неверный токен

## Коды ошибок

- `400 Bad Request` - Неверный формат запроса
- `401 Unauthorized` - Отсутствует или неверный токен авторизации
- `403 Forbidden` - Нет доступа к ресурсу
- `404 Not Found` - Ресурс не найден
- `409 Conflict` - Конфликт данных (например, email уже существует)
- `500 Internal Server Error` - Внутренняя ошибка сервера
