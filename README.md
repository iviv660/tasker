# 📝 Task Manager API

REST API для управления задачами (**todo-list**) с поддержкой регистрации пользователей, аутентификации через JWT и CRUD-операций по задачам.

----

## 🚀 Возможности

- 📌 **Регистрация и логин** пользователей (с хранением пароля в виде bcrypt-хэша).
- 🔑 **JWT-аутентификация**.
- ✅ **CRUD по задачам**:
    - создание задачи
    - получение списка задач
    - получение задачи по ID
    - обновление задачи
    - отметка выполненной
    - удаление задачи
- 📂 Привязка задач к пользователю (`owner_id`).
- 📖 Swagger UI для документации.

---

## ⚙️ Технологии

- **Go (Gin)** — веб-фреймворк
- **PostgreSQL** — база данных
- **migrate** — миграции
- **Docker + docker-compose** — запуск локально
- **swaggo** — генерация Swagger-документации

---

## 📂 Структура проекта

```bash
cmd/app/main.go       # точка входа
internal/config       # конфигурация (.env)
internal/database     # подключение к Postgres
internal/entity       # сущности (User, Task)
internal/repository   # работа с БД
internal/usecase      # бизнес-логика
internal/handler      # HTTP-эндпоинты (Gin)
internal/security     # хэширование пароля, JWT
internal/docs         # swagger-документация (сгенерированная)
migrations/           # SQL-миграции
docker-compose.yml
Dockerfile
README.md
```
---
## 🛠️ Установка и запуск

#### 1.Клонируй проект:
```bash
git clone https://github.com/<yourname>/tasker.git
cd tasker
```
#### 2.Создай .env (пример ниже):
```bash
POSTGRES_USER=postgres
POSTGRES_PASSWORD=070823
POSTGRES_DB=tasker

POSTGRES_URL=postgres://postgres:070823@postgres:5432/tasker?sslmode=disable

BASE_URL=http://localhost:3000
SECRET_KEY=Miromanov070823
```
#### 3.Запусти в Docker:
```bash
docker compose up --build
```
#### 4.API будет доступно по адресу:
```bash
http://localhost:3000
```
#### 5.Swagger UI:
```bash
http://localhost:3000/swagger/index.html
```


---


## 📌 Примеры запросов
```markdown
| Метод  | Endpoint              | Пример запроса                                                                                                          | Ответ (пример)   |
|--------|-----------------------|-------------------------------------------------------------------------------------------------------------------------|------------------|
| POST   | `/auth/register`      | `curl -X POST http://localhost:3000/auth/register -H "Content-Type: application/json" -d '{"email":"x","password":"y"}'`| `{"user_id":1}`  |
| POST   | `/auth/login`         | `curl -X POST http://localhost:3000/auth/login -H "Content-Type: application/json" -d '{"email":"x","password":"y"}'`   | `{"token":"..."}`|
| GET    | `/tasks`              | `curl -X GET http://localhost:3000/tasks -H "Authorization: Bearer <JWT>"`                                              | `{"tasks":[...]}`|
| POST   | `/tasks`              | `curl -X POST http://localhost:3000/tasks -H "Authorization: Bearer <JWT>" -d '{"title":"Test"}'`                       | `{...}`          |
| PUT    | `/tasks/{id}`         | `curl -X PUT http://localhost:3000/tasks/1 -H "Authorization: Bearer <JWT>" -d '{"title":"Update"}'`                    | `{...}`          | 
| PATCH  | `/tasks/{id}/complete`| `curl -X PATCH http://localhost:3000/tasks/1/complete -H "Authorization: Bearer <JWT>"`                                 | `{...}`          |
| DELETE | `/tasks/{id}`         | `curl -X DELETE http://localhost:3000/tasks/1 -H "Authorization: Bearer <JWT>"`                                         | `204 No Content` |
```

