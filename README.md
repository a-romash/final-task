# Финальная задача второго спринта

### Проект для Яндекс.Лицея

Отправная точка: http://localhost:8080/

## Техническое задание
<details>
  <summary><color> тык</color></summary>
  
  Пользователь хочет считать арифметические выражения. Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна. Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно.

Front-end часть

GUI, который можно представить как 4 страницы

Форма ввода арифметического выражения. Пользователь вводит арифметическое выражение и отправляет POST http-запрос с этим выражением на back-end. Примечание: Запросы должны быть идемпотентными. К запросам добавляется уникальный идентификатор. Если пользователь отправляет запрос с идентификатором, который уже отправлялся и был принят к обработке - ответ 200. Возможные варианты ответа:
200. Выражение успешно принято, распаршено и принято к обработке
400. Выражение невалидно
500. Что-то не так на back-end. В качестве ответа нужно возвращать id принятного к выполнению выражения.
Страница со списком выражений в виде списка с выражениями. Каждая запись на странице содержит статус, выражение, дату его создания и дату заверщения вычисления. Страница получает данные GET http-запрсом с back-end-а
Страница со списком операций в виде пар: имя операции + время его выполнения (доступное для редактирования поле). Как уже оговаривалось в условии задачи, наши операции выполняются "как будто бы очень долго". Страница получает данные GET http-запрсом с back-end-а. Пользователь может настроить время выполения операции и сохранить изменения.
Страница со списком вычислительных можностей. Страница получает данные GET http-запросом с сервера в виде пар: имя вычислительного ресурса + выполняемая на нём операция.

Требования:
Оркестратор может перезапускаться без потери состояния. Все выражения храним в СУБД.
Оркестратор должен отслеживать задачи, которые выполняются слишком долго (вычислитель тоже может уйти со связи) и делать их повторно доступными для вычислений.

Back-end часть

Состоит из 2 элементов:

Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его оркестратором.
Вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. Далее будем называть его агентом.
Оркестратор
Сервер, который имеет следующие endpoint-ы:

Добавление вычисления арифметического выражения.
Получение списка выражений со статусами.
Получение значения выражения по его идентификатору.
Получение списка доступных операций со временем их выполения.
Получение задачи для выполения.
Приём результата обработки данных.

Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения. При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. Количество горутин регулируется переменной среды.
</details>

## Запуск проекта: 
0) Установите [docker engine](https://docs.docker.com/engine/install/) и [docker compose](https://docs.docker.com/compose/install/)
1) `cd <"path/to/project/root/directory">`
2) `docker compose up -d` (`docker compose up -d --scale agent=X` <- если хотите несколько агентов на вычисление)

# Пример:
POST `http://localhost:8080/expression`
```json
{
    "expression": "((2*(3+4))/5)+(6-7)"
}
```

GET `http://localhost:8080/expression/2_p_2_m_1`
```json
"id": "2_p_2_m_1",
"expression": "2+2-1",
"answer": "3",
"status": "completed",
"createdAt": "2024-02-18 20:23:54.361591737 +0000 UTC",
"completedAt": "2024-02-18 20:23:59.376610678 +0000 UTC"
```


POST `http://localhost:8080/api/getimpodencekey`
```json
{
    "expression": "2+2"
}
```

```
"2_p_2"
```
## Вы можете получить данную ошибку: 

`Error response from daemon: Ports are not available: exposing port TCP 0.0.0.0:5673 -> 0.0.0.0:0: listen tcp 0.0.0.0:5673: bind: address already in use`

В этом случае вы должны поставить порт на любой другой свободный [docker-compose](docker-compose.yml)

Меняйте только проброшенные порты

Например:
- Было:
```yaml
  rabbitmq:
    container_name: rabbitmq
    image: rabbitmq:management
    ports:
    - 8081:15672
    - 5673:5672
```
- Стало:
```yaml
  rabbitmq:
    container_name: rabbitmq
    image: rabbitmq:management
    ports:
    - 8090:15672
    - 5812:5672
```

## Фронтенда нет

Фронтенда нет

# Правила для выражений

1) не должно быть никаких лишних символов
2) скобочки должны стоять правильно

# Схемки

![Схема всего проекта](/docs/schema.png)

![Схема всего проекта](/docs/rabbitmq.png)

![Схема всего проекта](/docs/orchestrator.png)