# DelayedNotifier - Сервис отложенных уведомлений
Сервис принимает запросы на создание уведомлений, складывает их в очередь (RabbitMQ), а затем, в нужное время — отправляет. Если что-то пошло не так, сервис попробует повторить отправку позже.

[Старт](https://github.com/andreyxaxa/Delayed-Notifier?tab=readme-ov-file#%D0%B7%D0%B0%D0%BF%D1%83%D1%81%D0%BA)

## Обзор

- UI - http://localhost:8080/v1/web
- Документация API - Swagger - http://localhost:8080/swagger
- Конфиг - [config/config.go](https://github.com/andreyxaxa/Delayed-Notifier/blob/main/config/config.go). Читается из `.env` файла.
- Логгер - [pkg/logger/logger.go](https://github.com/andreyxaxa/Delayed-Notifier/blob/main/pkg/logger/logger.go). Интерфейс позволяет подменить логгер.
- Кеширование статусов оповещений (Redis) - [internal/repo/cache/notification_cache.go](https://github.com/andreyxaxa/Delayed-Notifier/blob/main/internal/repo/cache/notification_redis.go).
- Удобная и гибкая конфигурация HTTP сервера - [pkg/httpserver/options.go](https://github.com/andreyxaxa/Delayed-Notifier/blob/main/pkg/httpserver/options.go).
  Позволяет конфигурировать сервер в конструкторе таким образом:
  ```go
  httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
  ```
  Аналогичный подход с таким конфигурированием в пакетах rabbitmq, redis, smtp sender, ...
- В слое контроллеров применяется версионирование - [internal/controller/restapi/v1](https://github.com/andreyxaxa/Delayed-Notifier/tree/main/internal/controller/restapi/v1).
  Для версии v2 нужно будет просто добавить папку `restapi/v2` с таким же содержимым, в файле [internal/controller/restapi/router.go](https://github.com/andreyxaxa/Delayed-Notifier/blob/main/internal/controller/restapi/router.go) добавить строку:
```go
{
		v1.NewNotificationRoutes(apiV1Group, n, l)
}

{
		v2.NewNotificationRoutes(apiV1Group, n, l)
}
```
- Graceful shutdown - [internal/app/app.go](https://github.com/andreyxaxa/Delayed-Notifier/blob/main/internal/app/app.go).

## Запуск

Сервис отправляет оповещения на почту через net/smtp, на телеграм через go-telegram-bot-api.

Соответственно, требуются :
- Gmail аккаунт, с него сервис будет слать оповещения. + app-pasword этого аккаунта. [Как создать пароль приложения](https://support.google.com/accounts/answer/185833?hl=ru)
- Токен телеграм бота. В телеграм найдите @BotFather - официальный бот для создания других телеграм-ботов. Создайте нового бота, получите токен.
  Также нужно знать, что телеграм-боты имеют право отправлять сообщения только тем юзерам, которые написали им первыми. API телеграма обязывает слать сообщения по chat ID, а не по @username. Найдите в телеграм @userinfobot - у него можно узнать chat ID любого юзера.

1. Клонируйте репозиторий
2. В корне создайте `.env` файл, скопируйте туда содержимое [env.example](https://github.com/andreyxaxa/Delayed-Notifier/blob/main/.env.example), подставив в `SMTPMAIL_USERNAME` ваш gmail, в `SMTPMAIL_PASSWORD` ваш app-password, в `TELEGRAM_BOT_TOKEN` токен вашего телеграм-бота.
3. Выполните, дождитесь запуска сервиса
   ```
   make compose-up
   ```
4. Перейдите на http://localhost:8080/v1/web и пользуйтесь сервисом.
<img width="931" height="981" alt="image" src="https://github.com/user-attachments/assets/2a2d2a4e-385e-491f-9513-64cfb483a5ba" />


## API

### POST http://localhost:8080/v1/notify
request:
```json
{
    "send_at": "2026-01-22T13:55:30+03:00",
    "payload": {
        "channel": "email",
        "email": {
            "to": "user@example.com",
            "subject": "Alex Birthday",
            "text": "Dont forget about Alexs birthday soon"
        }
    }
}
```
response:
```json
{
    "uid": "9a88d642-6c65-4f0f-b8f0-b920182cceb3",
    "status": "pending"
}
```

### GET http://localhost:8080/v1/notify/{id}
request:
```
GET http://localhost:8080/v1/notify/9a88d642-6c65-4f0f-b8f0-b920182cceb3
```
response:
```json
{
    "status": "pending"
}
```

### DELETE http://localhost:8080/v1/notify/{id}
request:
```
DELETE http://localhost:8080/v1/notify/9a88d642-6c65-4f0f-b8f0-b920182cceb3
```
response:
```json
{
    "status": "cancelled"
}
```

## Прочие `make` команды
Зависимости:
```
make deps
```
docker compose down:
```
make compose-down
```
