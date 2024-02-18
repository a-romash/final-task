# Финальная задача второго спринта

Проект для Яндекс.Лицея

Отправная точка: http://localhost:8080/
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
"id": 2_p_2_m_1,
"expression": 2+2-1,
"answer": 3,
"status": completed,
"createdAt": 2024-02-18 20:23:54.361591737 +0000 UTC,
"completedAt": 2024-02-18 20:23:59.376610678 +0000 UTC
```


Get `http://localhost:8080/api/getimpodencekey`
```
2_p_2
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
