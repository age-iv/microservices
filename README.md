```markdown
# News Microservices

Микросервисная система агрегатора новостей с комментариями и цензурой.

## Сервисы

| Сервис                | Порт  | Описание                                      |
|-----------------------|-------|-----------------------------------------------|
| **API Gateway**       | 8080  | Единая точка входа, проксирует запросы        |
| **News Service**      | 8081  | Хранение новостей, RSS-парсер, поиск          |
| **Comments Service**  | 8082  | Древовидные комментарии                       |
| **Censorship Service**| 8083  | Проверка текста на недопустимые слова         |
| **PostgreSQL**        | 5432  | Две базы: `newsdb` и `commentsdb`             |

## Требования

- Docker
- Docker Compose

## Запуск

```bash
# 1. Клонируйте репозиторий
git clone <repo-url> && cd news-microservices

# 2. Соберите и запустите все сервисы
docker-compose up --build
```

После запуска API Gateway доступен на `http://localhost:8080`.  
Новости автоматически подтягиваются из RSS‑лент (lenta.ru, interfax.ru) раз в 10 минут.

## Тестирование

Примеры запросов через `curl`.

### Список новостей (с пагинацией и поиском)

```bash
# Без параметров (первая страница)
curl http://localhost:8080/news

# Вторая страница
curl "http://localhost:8080/news?page=2"

# Поиск по заголовку
curl "http://localhost:8080/news?s=спорт"
```

### Детальная новость с комментариями

```bash
# Замените 1 на реальный ID новости
curl http://localhost:8080/news/1
```

### Создание комментария

```bash
# Комментарий, который пройдёт цензуру
curl -X POST http://localhost:8080/comments \
  -H "Content-Type: application/json" \
  -d '{"news_id":1, "text":"Отличная статья!"}'

# Комментарий с запрещённым словом (qwerty, йцукен, zxvbnm)
curl -X POST http://localhost:8080/comments \
  -H "Content-Type: application/json" \
  -d '{"news_id":1, "text":"qwerty - test"}'
```

### Сквозной идентификатор запроса

Передайте `request_id` для отслеживания запроса в логах:

```bash
curl "http://localhost:8080/news?request_id=my-trace-123"
```

Все сервисы логируют каждый запрос с IP, статусом и `request_id`.

## Структура проекта

```
├── api-gateway/             # API Gateway
├── news-service/            # Сервис новостей
├── comments-service/        # Сервис комментариев
├── censorship-service/      # Сервис цензуры
├── db/init/init.sql         # Создание баз данных
├── docker-compose.yml       # Описание всех сервисов
└── README.md
```
```
