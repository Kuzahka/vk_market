# Проект REST API на Go с Чистой Архитектурой, PostgreSQL и Docker

Этот проект представляет собой REST API для условного маркетплейса, реализованный на языке Go с использованием принципов **Чистой Архитектуры**.  
В качестве базы данных используется **PostgreSQL**, а для контейнеризации — **Docker**.

---

## Обзор Проекта

Приложение предоставляет следующие основные функции:

- **Авторизация и регистрация пользователей**  
  Создание новых учетных записей и вход в систему с получением токена.

- **Размещение объявлений**  
  Авторизованные пользователи могут создавать объявления.

- **Лента объявлений**  
  Список всех объявлений с пагинацией, сортировкой и фильтрацией.

---

## Архитектура

Проект построен по принципам **Чистой Архитектуры** — разделение ответственности, тестируемость и независимость от внешних технологий.

Основные слои:

- **Domain (Сущности)** — бизнес-объекты `User`, `Ad` и их правила.
- **Use Cases (Сценарии использования)** — бизнес-логика: `AuthUseCase`, `AdUseCase`.
- **Adapters (Адаптеры)** — HTTP-контроллеры и репозитории.
- **Infrastructure (Инфраструктура)** — реализация репозиториев, HTTP-сервер, утилиты.

---

## Структура Проекта

```
.
├── cmd/api/main.go            # Точка входа
├── internal/
│   ├── domain/                # Сущности
│   │   ├── user.go
│   │   └── ad.go
│   ├── usecase/               # Бизнес-логика
│   │   ├── auth.go
│   │   └── ad.go
│   ├── adapter/
│   │   ├── handler/           # HTTP-контроллеры
│   │   │   ├── auth_handler.go
│   │   │   ├── ad_handler.go
│   │   │   └── middleware.go
│   │   └── repository/        # Интерфейсы репозиториев
│   │       ├── user_repository.go
│   │       └── ad_repository.go
│   └── infrastructure/
│       ├── postgres/          # Репозитории PostgreSQL
│       │   ├── user_pg_repository.go
│       │   └── ad_pg_repository.go
│       ├── util/              # Утилиты
│       │   ├── password.go
│       │   └── token.go
├── migrations/                # Миграции БД
│   ├── 001_create_users_table.sql
│   └── 002_create_ads_table.sql
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Установка и Запуск

1. **Клонируйте репозиторий:**

```bash
git clone https://github.com/Kuzahka/vk_market
cd VK2
```

2. **Настройте модуль Go:**

```bash
go mod tidy
```

3. **Создайте `.env` файл:**

```env
DATABASE_URL="postgresql://user:password@db:5432/marketplace_db?sslmode=disable"
PORT="8080"
JWT_SECRET_KEY="секретный_ключ_для_токенов_минимум_32_символа"

POSTGRES_DB="marketplace_db"
POSTGRES_USER="user"
POSTGRES_PASSWORD="password"
```

4. **Запустите проект через Docker Compose:**

```bash
docker-compose down --volumes  # Очистка предыдущих томов
docker-compose up --build      # Сборка и запуск
```

- Приложение: http://localhost:8080  
- PostgreSQL: `localhost:5433` (внутри Docker — `5432`)

---

## API Эндпоинты

Все запросы — по адресу: `http://localhost:8080`

---

### 1. Регистрация Пользователя

**URL:** `/auth/register`  
**Метод:** `POST`  
**Content-Type:** `application/json`

**Тело запроса:**

```json
{
  "login": "myuser",
  "password": "MyStrongPassword123!"
}
```

**Пример cURL:**

```bash
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d '{"login": "myuser", "password": "MyStrongPassword123!"}'
```

---

### 2. Авторизация

**URL:** `/auth/login`  
**Метод:** `POST`  
**Content-Type:** `application/json`

**Тело запроса:**

```json
{
  "login": "myuser",
  "password": "MyStrongPassword123!"
}
```

**Пример cURL:**

```bash
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{"login": "myuser", "password": "MyStrongPassword123!"}'
```

---

### 3. Создание Объявления (Авторизация обязательна)

**URL:** `/ads`  
**Метод:** `POST`  
**Заголовок:** `Authorization: Bearer <ВАШ_ТОКЕН>`  
**Content-Type:** `application/json`

**Тело запроса:**

```json
{
  "title": "Продам старый велосипед",
  "description": "Отличный велосипед, почти новый...",
  "image_url": "https://placehold.co/600x400/000000/FFFFFF?text=Bicycle",
  "price": 150.00
}
```

**Пример cURL:**

```bash
curl -X POST http://localhost:8080/ads -H "Content-Type: application/json" -H "Authorization: Bearer <ВАШ_ТОКЕН>" -d '{
  "title": "Продам старый велосипед",
  "description": "Отличный велосипед...",
  "image_url": "https://placehold.co/600x400",
  "price": 150.00
}'
```

---

### 4. Получение Ленты Объявлений

**URL:** `/ads`  
**Метод:** `GET`

**Опциональные параметры:**

| Параметр     | Описание                        |
|--------------|----------------------------------|
| `page`       | Номер страницы (по умолчанию 1) |
| `limit`      | Кол-во на странице (по умолчанию 10) |
| `sort_by`    | Поле сортировки (`created_at`, `price`) |
| `sort_order` | Порядок сортировки (`asc`, `desc`)     |
| `min_price`  | Минимальная цена               |
| `max_price`  | Максимальная цена              |

**Пример cURL (базовый):**

```bash
curl -X GET http://localhost:8080/ads
```

**Пример cURL (с параметрами):**

```bash
curl -X GET "http://localhost:8080/ads?page=1&limit=5&sort_by=price&sort_order=desc&min_price=100&max_price=500"
```

---

> Для размещения объявлений необходим действующий JWT-токен, полученный при логине.

---

Готово! 
