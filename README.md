# DevPrep — Платформа для подготовки к техническим собеседованиям

DevPrep — это веб-приложение для подготовки к техническим интервью. Платформа предоставляет структурированные вопросы по
темам и уровням сложности, позволяет отслеживать прогресс, сохранять закладки и просматривать историю.

---

## Содержание

- [Возможности](#-возможности)
- [Архитектура](#-архитектура)
- [Технологический стек](#-технологический-стек)
- [Быстрый старт](#-быстрый-старт)
- [Структура проекта](#-структура-проекта)
- [API](#-api)
- [Конфигурация](#-конфигурация)
- [CI/CD](#-cicd)

---

## Возможности

- **Структурированная база знаний:** Вопросы организованы по темам (ООП, Java, Базы данных и т.д.) и уровням.
- **Отслеживание прогресса:** Пользователи могут отмечать вопросы как "Изучено", "Нужно повторить" или "Не знаю".
  Прогресс визуализируется по каждой теме.
- **Закладки:** Возможность сохранять важные вопросы для быстрого доступа.
- **История просмотров:** Автоматическое отслеживание недавно открытых вопросов.
- **Аутентификация через Keycloak:** Безопасный вход с использованием OpenID Connect (OAuth2).
- **Кэширование:** Использование Redis для ускорения ответов API.
- **Документация API:** Интерактивная Swagger-документация для бэкенда.

---

## Архитектура

Проект построен по микросервисной архитектуре и состоит из трёх основных сервисов:

```
┌─────────────┐     OAuth2/OIDC     ┌─────────────┐
│   Frontend  │◄───────────────────►│  Keycloak   │
│ Spring MVC  │                     │             │
│ Thymeleaf   │                     └─────────────┘
└──────┬──────┘
       │ HTTP (Bearer JWT)
       ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Backend   │◄───►│ PostgreSQL  │     │    Redis    │
│   Go API    │     │             │     │    (кэш)    │
└─────────────┘     └─────────────┘     └─────────────┘
```

* **Nginx:** Выступает в роли единой точки входа. Он принимает все HTTP-запросы и направляет их на фронтенд (для
  пользовательских страниц) или на бэкенд (для API-вызовов). Также обслуживает Swagger-документацию.
* **Frontend (Spring Boot):** Отвечает за отрисовку HTML-страниц с помощью Thymeleaf. Он общается с бэкендом через REST
  API, автоматически добавляя JWT-токен аутентифицированного пользователя в запросы.
* **Backend (Go):** Основной REST API сервис. Предоставляет эндпоинты для получения вопросов, тем, тегов, а также для
  управления прогрессом пользователя, закладками и историей. Интегрируется с Keycloak для валидации JWT-токенов и с
  PostgreSQL для хранения данных. Использует Redis для кэширования часто запрашиваемых данных (список тем, детали
  вопроса).
* **Keycloak:** Сервер аутентификации. Обрабатывает процесс входа пользователей, выдачу JWT-токенов и управление
  пользовательскими сессиями.
* **PostgreSQL:** Основная реляционная база данных, хранящая информацию о вопросах, темах, тегах, а также данные о
  прогрессе, закладках и истории пользователей.
* **Redis:** Высокопроизводительное хранилище ключ-значение в памяти. Используется бэкендом для кэширования данных (
  темы, вопросы), чтобы снизить нагрузку на основную базу данных и ускорить ответы API.

Все сервисы запускаются через **Docker Compose**.

---

## Технологический стек

| Компонент        | Технология                                                                       |
|------------------|----------------------------------------------------------------------------------|
| Frontend         | Java 21, Spring Boot 3, Thymeleaf, Tailwind CSS, Maven, Lombok                   |
| Backend          | Go 1.25, Chi Router, pgx, go-redis, jwt-go, keyfunc, redis/go-redis, swaggo/swag |
| База данных      | PostgreSQL 16                                                                    |
| Кэш              | Redis 7                                                                          |
| Аутентификация   | Keycloak 24, OAuth2, JWT                                                         |
| Документация API | Swagger / OpenAPI (swaggo)                                                       |
| Инфраструктура   | Docker, Docker Compose, Nginx                                                    |
| CI/CD            | Jenkins                                                                          |
| Миграции         | golang-migrate                                                                   |

---

## Быстрый старт (локальный запуск)

**Предварительные требования**

- Docker и Docker Compose
- Java 21, Maven (для сборки фронтенда)
- Go 1.25 (для сборки бэкенда)

1. Клонируйте репозиторий:

```bash
git clone https://github.com/ArtemFedorov2004/dev-preparation
cd dev-preparation
```

2. Создайте файл окружения `.env`

3. Запустите все сервисы:

```bash
cd infra
docker compose up -d
```

4. Дождитесь готовности Keycloak, затем откройте браузер:

| Сервис      | URL                           |
|-------------|-------------------------------|
| Frontend    | http://localhost:8082         |
| Backend API | http://localhost:8081/api/v1  |
| Swagger UI  | http://localhost:8081/swagger |
| Keycloak    | http://localhost:8443         |

### Запуск Frontend в режиме mock (без бэкенда)

```bash
cd frontend
mvn spring-boot:run -Dspring-boot.run.profiles=mock
```

---

## Структура проекта

```
dev-preparation/
├── backend/                    # Go REST API
│   ├── cmd/
│   │   ├── api/main.go         # Точка входа сервера
│   │   └── migrate/main.go     # CLI для миграций
│   ├── internal/
│   │   ├── config/             # Конфигурация через env-переменные
│   │   ├── handler/            # HTTP-обработчики
│   │   ├── service/            # Бизнес-логика
│   │   ├── repository/         # Слой данных (PostgreSQL + Redis-кэш)
│   │   ├── model/              # Доменные модели
│   │   ├── middleware/         # Auth, CORS, Logging
│   │   └── keycloak/           # JWKS-клиент для верификации JWT
│   ├── migrations/             # SQL-миграции (embedded)
│   └── docs/                   # Сгенерированный Swagger
│
├── frontend/                   # Spring MVC + Thymeleaf
│   ├── src/main/java/com/devprep/
│   │   ├── controller/         # MVC-контроллеры
│   │   ├── client/             # HTTP-клиент к backend API
│   │   ├── config/             # Security, OAuth2
│   │   └── security/           # OAuth2 interceptor (прокидывает JWT)
│   └── src/main/resources/
│       └── templates/          # Thymeleaf-шаблоны
│
└── infra/                      # Инфраструктура
    ├── docker-compose.yml      # Базовая конфигурация
    ├── docker-compose.prod.yml # Продакшн конфигурация
    ├── nginx/nginx.conf        # Конфигурация nginx
    ├── postgres/init.sql       # Инициализация схемы Keycloak
    └── keycloak/               # Экспорт realm-конфигурации
```

---

## API

Полная документация доступна через Swagger UI по адресу `/swagger`.

### Публичные эндпоинты

| Метод | Путь                       | Описание                             |
|-------|----------------------------|--------------------------------------|
| GET   | `/api/v1/topics`           | Список всех тем                      |
| GET   | `/api/v1/topics/{slug}`    | Тема с вопросами (фильтр по `level`) |
| GET   | `/api/v1/questions`        | Список вопросов (фильтр + пагинация) |
| GET   | `/api/v1/questions/{slug}` | Вопрос с полным ответом              |
| GET   | `/api/v1/tags`             | Список тегов                         |

### Защищённые эндпоинты

| Метод | Путь                                | Описание                         |
|-------|-------------------------------------|----------------------------------|
| POST  | `/api/v1/questions/{slug}/progress` | Установить статус изучения       |
| GET   | `/api/v1/questions/{slug}/progress` | Получить статус изучения         |
| POST  | `/api/v1/questions/{slug}/bookmark` | Добавить / убрать закладку       |
| GET   | `/api/v1/questions/{slug}/bookmark` | Статус закладки                  |
| POST  | `/api/v1/questions/{slug}/view`     | Записать просмотр                |
| GET   | `/api/v1/me/progress`               | Весь прогресс пользователя       |
| GET   | `/api/v1/me/progress/by-topic`      | Агрегированный прогресс по темам |
| GET   | `/api/v1/me/bookmarks`              | Список закладок                  |
| GET   | `/api/v1/me/history`                | История просмотров               |

---

## Конфигурация

### Backend (переменные окружения)

| Переменная       | По умолчанию            | Описание          |
|------------------|-------------------------|-------------------|
| `SERVER_PORT`    | `8081`                  | Порт HTTP-сервера |
| `DB_HOST`        | `localhost`             | Хост PostgreSQL   |
| `DB_PORT`        | `5432`                  | Порт PostgreSQL   |
| `DB_NAME`        | `devprep`               | Имя базы данных   |
| `DB_USER`        | `devprep`               | Пользователь БД   |
| `DB_PASSWORD`    | `devprep`               | Пароль БД         |
| `KEYCLOAK_URL`   | `http://localhost:8443` | URL Keycloak      |
| `KEYCLOAK_REALM` | `devprep`               | Realm Keycloak    |
| `REDIS_ADDR`     | `localhost:6379`        | Адрес Redis       |
| `REDIS_PASSWORD` | —                       | Пароль Redis      |

### Frontend (переменные окружения)

| Переменная               | По умолчанию            | Описание                  |
|--------------------------|-------------------------|---------------------------|
| `API_BASE_URL`           | `http://localhost:8081` | URL backend API           |
| `KEYCLOAK_URL`           | `http://localhost:8443` | URL Keycloak              |
| `KEYCLOAK_REALM`         | `devprep`               | Realm Keycloak            |
| `KEYCLOAK_CLIENT_ID`     | `devprep-frontend`      | Client ID OAuth2          |
| `KEYCLOAK_CLIENT_SECRET` | *(из Keycloak)*         | Client Secret OAuth2      |
| `SPRING_PROFILES_ACTIVE` | —                       | Профиль spring приложения |

---

## CI/CD

Проект использует **Jenkins Pipeline** с автоматическим деплоем на VPS.

### Этапы пайплайна

```
Checkout → Build Frontend (Maven) → Build Backend (Go)
    → Build Docker Images → Push to Registry → Deploy to VPS
```

### Требования

В Jenkins должны быть настроены следующие credentials:

| ID                 | Тип             | Описание                       |
|--------------------|-----------------|--------------------------------|
| `vps-host`         | Secret text     | IP-адрес VPS                   |
| `vps-ssh-key`      | SSH private key | SSH-ключ для подключения к VPS |
| `devprep-env-file` | Secret file     | Файл `.env` для продакшна      |

---

## Миграции базы данных

Миграции применяются автоматически при старте бэкенда. Для ручного управления:

```bash
# Применить все миграции
go run ./cmd/migrate up

# Откатить последнюю миграцию
go run ./cmd/migrate down

# Посмотреть текущую версию
go run ./cmd/migrate version

# Принудительно установить версию (исправить dirty state)
go run ./cmd/migrate force <version>
```