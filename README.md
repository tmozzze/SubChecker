# SubChecker

SubChecker - REST-сервис для агрегации данных об онлайн-подписках пользователей.

## Структура проекта

```
.
├── cmd/
│   └── app/
│       └── main.go        # Точка входа в приложение
├── database/
│   └── migrations         # SQL миграции
├── docs/                  # Swagger докс
├── docker-compose.yml     # Конфигурация Docker контейнеров
├── internal/
│   ├── http/              # HTTP обработчики
│   ├── config/            # Конфигурация приложения
│   ├── db/                # Работа с базой данных
│   ├── models/            # Модели данных
│   ├── repository/        # Репозиторий для работы с БД
│   ├── service/           # Бизнес-логика
│   └── utils/             # Вспомогательные функции
├── Makefile               # Команды для управления проектом
├── Dockerfile
├── go.mod                 # Модуль
├── go.sum
├── test.txt               # Пример запроса на создание подписки
├── .env.example           # Пример .env
└── .gitignore
```

### Требования для запуска


- `golang-migrate CLI` для управления миграциями
- `Docker` и `Docker-compose`



## Запуск проекта

### Предварительные требования

- Docker и Docker Compose должны быть установлены на вашей системе.
- Убедитесь, что файл `.env` существует и сделан по образцу `.env.example`

### Шаги запуска

1. Запустите сборку в терминале bash:
   ```sh
   make start-service
   ```

2. Дождитесь запуска сервиса

3. Откройте Swagger-документацию по адресу:
   ```
   http://localhost:8080/swagger/index.html
   ```

Теперь ваш сервис готов к использованию!

# Команды для управления проектом

### Запускает сервер Go.
```
make start-service
```

### Создаёт новую миграцию SQL.
```
make create-migration NAME=<имя_миграции>
```
- Файлы миграции создаются в папке database/migrations

### Применяет все миграции к базе данных.
```
make migrate-up-all
```
- Использует golang-migrate CLI
- Применяет все .sql миграции из ./database/migrations

### Применяет только одну миграцию (последнюю).
```
make migrate-up
```
### Применяет только одну миграцию (последнюю).
```
make migrate-down
```
### Показывает текущую версию миграций в базе.
```
make migrate-status
```
### Откатывает все миграции (полностью очищает базу).
```
make migrate-reset
```
### Останавливает контейнеры Docker.
```
make down             # остановка контейнеров
make down-and-clean   # остановка и удаление томов
```
### Генерирует Swagger документацию для API.
```
make swagger-gen
```
- Сканирует cmd/app/main.go
- Создаёт документацию в папке docs

# Основные компоненты
- Gin (HTTP сервер)
- PostgreSQL (СУБД)
- pgx (PostgreSQL драйвер)
- Logrus (Логирование)
- Swagger (Документация API)