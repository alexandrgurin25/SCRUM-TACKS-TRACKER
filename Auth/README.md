# Auth 

Микросервис аутентификации и авторизации с использованием gRPC и PostgreSQL.

## Функциональность

- Регистрация пользователей
- Аутентификация пользователей (логин)
- Обновление токенов (access и refresh)
- Верификация access токена

## Технологии

- Go 1.24
- gRPC
- PostgreSQL
- JWT (JSON Web Tokens)
- Docker

## Запуск проекта

1. Перейдите в директорию проекта:
   ```shell
   cd Auth
   ```
2. Запустите сервисы
   ```shell
   docker-compose up --build
   ```
   
## Запуск проекта
  Сборка и запуск: ```docker-compose up --build```<br>
  Генерация GraphQL кода: ```make gen``` <br>
  Запуск в development режиме: ```make run```<br>
  
## Структура проекта
```
Auth/
├── api/
│   └── proto/           # Protobuf спецификация
├── cmd/
│   └── app/             # Точка входа
├── config/              # Конфигурация
├── db/
│   └── migrations/      # SQL миграции
├── internal/
│   ├── entity/          # Бизнес-сущности
│   ├── repository/      # Работа с БД
│   ├── service/         # Бизнес-логика
│   └── transport/       # gRPC handlers
└── pkg/
    ├── logger/          # Логирование
    ├── postgres/        # Подключение к PG
    └── proto/           # Сгенерированный код
```
