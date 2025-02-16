## Welcome to Avito-shop

### О проекте:

Проект представляет собой WEB-сервер, реализующий следующий функционал:

- Авторизация пользователя
- Покупка мерча
- Перевод монеток между пользователями
- Просмотр информации о пользователе (баланс, купленный мерч, история транзакций)

Реализовано юнит-тестирование (83% покрытия) и интеграционное тестирование бизнес сценариев.

### Сборка и запуск

```sh
docker compose build
docker compose up -d
```

### Запуск интеграционных тестов

```sh
docker compose up test-runner
```


### Решения в ходе разработки

1. При проектировании базы данных возник выбор в связи между таблицами `merch` и `inventory`. Всё же в аналогичных системах скорее ожидается владение юзера определенным уникальным товаром, а не его "типом". Я решила связать сущности именно с помощью `merch.id`, а не `merch.name`, в целях нормализации таблиц в БД.
2. В целях удобства реализации взаимодействия с БД (а также внедрения юнит-тестов) был разработан слой репозитория: интерфейс, реализующийся по-отдельности для работы с postgres и для мокания в рамках юнит-тестов.
3. Для того, чтобы репозитории могли работать и в рамках транзакций, и в рамках прямых запросов к БД, был введен интерфейс **Executor** (`config/db`).
4. В целях стандартизации при отсутствии тех или иных компонентов-списков в JSON возвращаются пустые списки, а не null-значения.
5. Для использования в юнит-тестах функции для управления транзакциями в `config/db/transaction.go` были объявлены как переменные.