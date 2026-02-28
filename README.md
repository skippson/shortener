# Shortener

## Описание
Сервис Shortener предназначен для генерации коротких идентификаторов для URL с возможностью хранения данных в памяти или в PostgreSQL

Релизует REST API:

* создания короткой ссылки
* получения оригинального URL по идентификатору

### Контракт
* POST /api/create_shortened 
* * Создание короткой ссылки

    Тело Запроса:
    ```json
    {
        "url":"http://example.com"
    }
    ```

    Тело ответа:

    200
    ```json
    {
        "data": {
            "shortened": "QbdEIWlNDV"
        }
    }  
    ```

    4xx/5xx
    ```json
    {
        "message": "error message",
        "status": 4xx/5xx
    }
    ```

* GET /api/get_original/:shortened
* * Получение оригинального `URL`

    Тело запроса:

    URL параметр `shortened`

    Тело ответа:

    200
    ```json
    {
        "data": {
            "original": "http://example.com"
        }
    }
    ```

    4xx/5xx
    ```json
    {
        "message": "error message",
        "status": 4xx/5xx
    }
    ```

## Локальное развертывание
* Для настройки переменных окружения смотрите `.example.env`
    * * `GENERATOR_*` - конфигурация генерации `shortened` (обязательные)
    * * `SERVICE_IN_MEMORY_MODE` - режим хранения в памяти
    * * `SERVICE_PROTECTION` - включение валидации `URL` и `shortened`
    * * `SERVICE_MAX_GENERATE_ATTEMPTS` - максимальное количество попыток генерации `shortened`

* Запуск
    ```
    docker-compose up --build -d
    ```

* Сервис готов к работе

* Остановка
    ```
    docker-compose down -v
    ```

## Примеры запросов
* Создание короткой ссылки
    ```
    curl -X POST http://localhost:8080/api/create_shortened \
        -H "Content-Type: application/json" \
        -d '{"url":"https://google.com"}'
    ```

* Получение оригинального `URL`
    ```
    curl -X GET http://localhost:8080/api/get_original/rhJscUXqZi
    ```

## Документация
* `config` - Установка конфига
* `internal/adapters/repository` - Контракт репозитория
* * `memory` - Релизация и логика хранения в памяти
* * `postgres` - Взаимодействия с базой данных
* * * `migrations` - Миграции базы данных
* `internal/controllers/http_handlers` - Транспортный слой(реализация запросов)
* * `middleware` - Промежуточная логика
* `internal/domain` - Доменные модели(ошибки)
* `internal/generator` - Генерация `shortened`
* `internal/server` - Реализация сервера
* `internal/usecase` - Бизнес-логика
* `internal/validator` - Валидация `URL` и `shortened`
* `pkg/logger` - Логгер модель