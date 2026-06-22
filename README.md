# TODO-приложение: REST API на Go

## 📋 Описание проекта
Проект реализует полноценный backend на Go (с REST API, аутентификацией JWT + Refresh)

### Технологии
- REST API с `gin-gonic/gin`
- PostgreSQL через `sqlx`
- Миграции бд через  `migrate`
- Docker-среда разработки
- Конфигурация приложения через `cleanenv`
- JWT (передача в заголовке) аутентификация и Refresh токены (HTTPOnly)

## ⚙️ Как запустить проект

### 1. Клонировать репозиторий
```bash
git clone https://github.com/wesorat/todo.git

cd todo
```
### 2. Создать .env файл на основе .env.example

### 3. Запустить через Docker
```bash
docker-compose up --build -d
```
### 4. Применить миграции
```bash
migrate -path migration -database 'postgres://postgres:1234@0.0.0.0:5433/todo_db?sslmode=disable' up
```