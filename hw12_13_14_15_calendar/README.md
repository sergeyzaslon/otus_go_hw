#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:
- [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
- [Домашнее задание №13 «Внешние API от Календаря»](./docs/13_README.md)
- [Домашнее задание №14 «Кроликизация Календаря»](./docs/14_README.md)
- [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)

#### Ветки при выполнении
- `hw12_calendar` (от `master`) -> Merge Request в `master`
- `hw13_calendar` (от `hw12_calendar`) -> Merge Request в `hw12_calendar` (если уже вмержена, то в `master`)
- `hw14_calendar` (от `hw13_calendar`) -> Merge Request в `hw13_calendar` (если уже вмержена, то в `master`)
- `hw15_calendar` (от `hw14_calendar`) -> Merge Request в `hw14_calendar` (если уже вмержена, то в `master`)

**Домашнее задание не принимается, если не принято ДЗ, предшедствующее ему.**

Планировщик
 - репозиторий событий
 - выбрать те, о которых пора напомнить
 - пометить, что уже напомнили?
 - подключиться к Rabbit и создать очереди
 - создать Уведомление
 - отправить Уведомление в очередь
 - очистить старые события (более 1 года)


 ## Переменные окружения
|Название|Пример|Описание|
|--------|------|--------|
|LOG_LEVEL|`debug`|`[panic, fatal, error, warn, info, debug]`|
|LOG_FILE|`stderr`|Абсолютный путь к файлу или `[stderr|stdout]`|
|LOG_FORMAT|`text`|`[text|json]`|
|STORAGE_TYPE|`sql`|`[sql|memory]`|
|STORAGE_DSN|`postgres://calendar:calendar@localhost:5434/calendar`|Можно использовать `memory://` для хранения в памяти|
|HTTP_HOST|`0.0.0.0`||
|HTTP_PORT|`80`||
|GRPC_HOST=|`0.0.0.0`||
|GRPC_PORT|`8080`||
|QUEUE_DSN|`amqp://guest:guest@rabbit:5673/`|Подключение к очереди сообщений|
|RABBIT_QUEUE|`event_notifications`|Название очереди в RabbitMQ|
|RABBIT_EXCHANGE|`calendar`|Название exchange в RabbitMQ|


```
docker-compose up -d
docker-compose up migrations
```
