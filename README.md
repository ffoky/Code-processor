# Пет-проект на Go, написан по принципам чистой архитектуры, реализован микросервисный подход и работу с брокером сообщений, контейнерами и базами данных. 

---
## Содержание

- [Технологии](#технологии) 
- [Архитектура](#архитектура)  
- [Основные компоненты](#основные-компоненты)  
- [Построение и запуск](#построение-и-запуск)  
- [Docker Compose](#docker-compose)  
- [Процессоры сообщений](#процессоры-сообщений)  
- [Инфраструктура](#инфраструктура)  
- [Мониторинг](#мониторинг)  

# Технологии
- **Go** (1.21), **Chi**, **Swaggo**
- **Docker**, **Docker Compose**
- **RabbitMQ**
- **PostgreSQL**, **Redis**
- **Prometheus**, **Grafana**
- **GitHub Actions**, **Makefile**


# Архитектура

## Общая архитектура микросервисов

```mermaid
flowchart LR
  %% Core components
  HTTPServer["HTTP Server"] 
  RabbitMQ["RabbitMQ"] 
  Processor["Code / Image Processor"]

  %% Horizontal connections
  HTTPServer <--> RabbitMQ
  RabbitMQ <--> Processor

  %% Vertical dependencies
  HTTPServer --> PostgreSQL["PostgreSQL (Task/User storage)"]
  Processor --> Storage["Redis (sessions) & Filesystem (results)"]

  %% Styling
  classDef core fill:#e3f2fd,stroke:#1976d2,stroke-width:2px;
  class HTTPServer,RabbitMQ,Processor core;
```

## Clean Architecture (обобщённая схема)

Здесь отражён стандартный поток данных: от источников хранения через репозиторий, сервисы и контроллеры к внешним интерфейсам.

```mermaid
flowchart LR

  %% Data sources

  RDBMS[(RDBMS)]

  NoSQL[(NoSQL)]

  Micro[(Microservices)]

  %% Core components

  Repo[Repository]

  Domain[Domain/Model/Entity]

  Usecase[Usecase/Service]

  Controller[Controller/Delivery]

  %% Interfaces

  gRPC((gRPC))

  REST((REST))

  CLI[/CLI/]

  Web[/Web/]

  %% Relationships

  RDBMS --> Repo

  NoSQL --> Repo

  Micro --> Repo

  Domain --> Repo

  Repo --> Usecase

  Usecase --> Controller

  Domain --> Controller

  Controller --> gRPC

  Controller --> REST

  Controller --> CLI

  Controller --> Web

  %% Annotation

  classDef note fill:#fff8c6,stroke:#aaa,stroke-dasharray: 2 2;

  Note["Business logic happens here"]:::note

  Usecase --> Note
```

## Слои Clean Architecture в моём Go-проекте

```mermaid
flowchart LR

  subgraph Domain [Domain Entities]

    D1[Task]

    D2[User]

  end

  

  subgraph UseCases [Use Cases / Services]

    U1[TaskService]

    U2[UserService]

  end

  

  subgraph Adapters [Interface Adapters]

    A1[HTTP Handlers]

    A2[PostgresRepo]

    A3[RedisRepo]

    A4[RabbitMQ Client]

  end

  

  subgraph Infra [Frameworks & Drivers]

    F1(Chi HTTP)

    F2(PostgreSQL)

    F3(Redis)

    F4(RabbitMQ)

    F5(Docker)

    F6(Prometheus/Grafana)

  end

  

  %% зависимости

  F1 --> A1

  F2 --> A2

  F3 --> A3

  F4 --> A4

  A1 --> U1

  A1 --> U2

  A2 --> U1

  A3 --> U2

  A4 --> U1

  U1 --> D1

  U2 --> D2
```
- **HTTP Server** (Go, Chi) — REST API, сессии, Swagger-документация.
- **RabbitMQ** — брокер сообщений.
- **Processor** — отдельный сервис для выполнения кода (CodeProcessor) или фильтрации изображений (ImageProcessor).
- **PostgreSQL** — хранение пользователей и задач.
- **Redis** — хранение сессий.
- **Prometheus + Grafana** — сбор и визуализация метрик.
## ## Процессоры сообщений

### CodeProcessor

- На Go + Docker SDK
- Для каждой новой задачи создаёт Docker-контейнер с компиляторами (clang, gcc, python)
- Компилирует/запускает код внутри контейнера
- Возвращает `stdout`/`stderr`
# Основные компоненты

### 1. HTTP-сервер (`cmd/app`)
- **Авторизация** – `AuthMiddleware` проверяет JWT-сессии в Redis.
- **API по задачам** – CRUD-эндпоинты `/task`, `/status/{id}`, `/result/{id}`.
- **API по пользователям** – регистрация `/register`, логин `/login`.
- **Документация** – Swagger UI доступен по `/swagger/*`.
### 2. Сервис работы с задачами
- **TaskService** создаёт задачу, отправляет сообщение в RabbitMQ и хранит метаданные.
- **SessionService** управляет сессиями через Redis.
### 3. Репозитории
- **RamStorage** для быстрого прототипирования (в памяти).
- **PostgreSQL** для хранения Users и Tasks (миграции через `migrate`).
- **Redis** для сессий.
- **RabbitMQSender** для публикации сообщений в очередь `tasks`.
### 4. Processor-микросервисы (`processor/`)
- **CodeProcessor**:  
  - Получает JSON `{ "code": "...", "lang": "c" }` из очереди.  
  - В Docker-контейнере собирает и запускает чужой код (clang, gcc, python).  
  - Возвращает `stdout`/`stderr` через HTTP `/commit`.

---

# Инфраструктура

- **Docker & docker-compose**  
  - Сервисы: `http-server`, `rabbitmq`, `processor`, `postgres`, `redis`.  
  - Общая сеть, тома для БД и Redis.

- **Makefile**  
  - `make build` – сборка всех бинарников и Docker-образов.  
  - `make up` – поднять `docker-compose`.  
  - `make test` – запуск unit & integration тестов.

- **CI/CD (GitHub Actions)**  
  - Проверка `go fmt`/`go vet`, сборка, тесты.  
  - Линтинг Dockerfile, автоматический запуск `docker-compose` и smoke-тесты.  
  - Отдельная ветка `ci` с зелёными галочками.

---

# Метрики и мониторинг

- **Prometheus**  
  - Собирает метрики об обработке запросов и времени выполнения Processor’ов.  
- **Grafana**  
  - Дашборд с графиками latency, throughput и используемых фильтров/компиляторов.

---


  

