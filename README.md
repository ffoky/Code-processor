# Code-processor - раннер кода на Go

---
## Содержание

- [Технологии](#технологии) 
- [Архитектура](#архитектура)  
- [Основные компоненты](#основные-компоненты)  
- [Процессоры сообщений](#процессоры-сообщений)  
- [Инфраструктура](#инфраструктура)  
- [Мониторинг](#мониторинг)  

## Технологии

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![Chi](https://img.shields.io/badge/Chi-5.0+-000000?logo=go&logoColor=white)
![Swaggo](https://img.shields.io/badge/Swaggo-1.8+-34ABE0?logo=swagger&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-24.0+-2496ED?logo=docker&logoColor=white)
![Docker Compose](https://img.shields.io/badge/Docker_Compose-2.23+-2496ED?logo=docker&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.12+-FF6600?logo=rabbitmq&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7.0+-DC382D?logo=redis&logoColor=white)
![Grafana](https://img.shields.io/badge/Grafana-10.1+-F46800?logo=grafana&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-3.0+-2088FF?logo=github-actions&logoColor=white)
![Makefile](https://img.shields.io/badge/Makefile-GNU+-000000?logo=gnu&logoColor=white)

## Архитектура

### Общая архитектура микросервисов

```mermaid
flowchart LR 
  HTTPServer["HTTP Server"] 
  RabbitMQ["RabbitMQ"] 
  Processor["Code Processor"]
  PostgreSQL["PostgreSQL (Task/User storage)"]
  Storage["Redis (sessions) & Filesystem (results)"]

  HTTPServer <--> RabbitMQ
  RabbitMQ <--> Processor

  HTTPServer --> PostgreSQL
  Processor --> Storage

  %% Стилизация только рамки (чтобы сохранить адаптивность фона!)
  classDef core stroke:#1976d2,stroke-width:2px;
  class HTTPServer,RabbitMQ,Processor,PostgreSQL,Storage core;
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
  Usecase[Usecase/Service]
  Controller[Controller/Delivery]

  %% "Фикс" позиции Domain над Usecase
  Domain/Model/Entity
  Domain/Model/Entity --> Usecase
  Domain/Model/Entity --> Repo
  Domain/Model/Entity --> Controller

  %% Interfaces
  gRPC((gRPC))
  REST((REST))
  CLI[/CLI/]
  Web[/Web/]

  %% Relationships
  RDBMS --> Repo
  NoSQL --> Repo
  Micro --> Repo
  Repo --> Usecase
  Usecase --> Controller
  Controller --> gRPC
  Controller --> REST
  Controller --> CLI
  Controller --> Web

  %% Annotation
  classDef note fill:#f7f7f7,color:#666,stroke:#aaa,stroke-dasharray: 2 2;
  Note["Business logic happens here"]:::note
  Usecase --> Note

  %% Styling
  classDef outlined stroke:#1976d2,stroke-width:2px;
  classDef dummy fill:transparent,stroke:transparent;
  class RDBMS,NoSQL,Micro,Repo,Domain,Usecase,Controller,gRPC,REST,CLI,Web outlined;
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

  %% Только синяя рамка для всех блоков, без fill!
  classDef outlined stroke:#1976d2,stroke-width:2px;
  class D1,D2,U1,U2,A1,A2,A3,A4,F1,F2,F3,F4,F5,F6 outlined;
```
- **HTTP Server** (Go, Chi) — REST API, сессии, Swagger-документация.
- **RabbitMQ** — брокер сообщений.
- **Processor** — отдельный сервис для выполнения кода (CodeProcessor) или фильтрации изображений (ImageProcessor).
- **PostgreSQL** — хранение пользователей и задач.
- **Redis** — хранение сессий.
- **Prometheus + Grafana** — сбор и визуализация метрик.
## Процессоры сообщений

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

# Мониторинг

- **Prometheus**  
  - Собирает метрики об обработке запросов и времени выполнения Processor’ов.  
- **Grafana**  
  - Дашборд с графиками latency, throughput и используемых фильтров/компиляторов.

---


  

