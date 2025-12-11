# Classic SQL Injection Module – Report

## I. Goal and Tasks of the Project

### Goal

Показать, как классическая SQL-инъекция в веб-приложении приводит к обходу аутентификации и утечке данных из базы PostgreSQL, и продемонстрировать практическую эксплуатацию с помощью ручных payload-ов и sqlmap.

### Scope (мой модуль)

- Реализовать уязвимый веб‑сервис (classic SQL injection) на Go.
- Настроить PostgreSQL в Docker с тестовыми пользователями и продуктами.
- Реализовать уязвимые эндпоинты:
    - `/login` — обход аутентификации.
    - `/search` — UNION-based SQL-инъекция.
- Подготовить скрипты эксплуатации (curl + sqlmap).
- Оформить документацию по атаке и защите.

### Моя зона ответственности

- Дизайн и реализация сервиса `sql-injection/app` (Go + HTTP + DB).
- Подготовка схемы БД и начальных данных для `users` и `products`.
- Написание эксплойтов в `sql-injection/exploits`.
- Описание attack surface, сценариев атак и механизмов защиты для classic SQLi.

---

## II. Execution Plan and Methodology

### High-Level Plan

1. Спроектировать минимальное веб‑приложение с формой логина и поиском по товарам.
2. Поднять PostgreSQL в Docker с отдельной БД `vuln`.
3. Специально внедрить уязвимость SQL‑инъекции через конкатенацию строк в SQL.
4. Подтвердить эксплуатацию вручную (curl, браузер).
5. Подтвердить эксплуатацию автоматически (sqlmap).
6. Задокументировать уязвимость, атакующие сценарии и меры защиты.

### Архитектура

- **Сервис**: `sql-injection/app` (Go, `net/http`, `database/sql`, `github.com/lib/pq`).
- **База данных**: PostgreSQL в Docker, схема `public`.
- **Связь**: сервис подключается к БД по env‑переменным:

```text
POSTGRES_USER=admin
POSTGRES_PASSWORD=pass
POSTGRES_DB=vuln
```


В приложении параметры соединения берутся из переменных окружения и собираются в DSN для `database/sql`.

- **Развёртывание**: через `docker-compose` с двумя сервисами:
- `db` — PostgreSQL с инициализацией из `database/init.sql`.
- `sql-injection` — Go‑приложение, собранное из `sql-injection/app/Dockerfile`.

### Инфраструктура (словесная схема)

- Пользователь → HTTP (localhost:8081) → Go‑сервер (`/login`, `/search`).
- Go‑сервер → PostgreSQL (по Docker‑сети, хост `db`, БД `vuln`).
- Эксплойты (shell + sqlmap) → HTTP‑запросы к тем же эндпоинтам.

---

## III. Development of Solution and Tests (PoC)

### 1. Attack Surface and Endpoints

**Входные точки:**

1. `POST /login`
- Параметры: `username`, `password`.
- Назначение: аутентификация пользователя.
- Риск: обход аутентификации, доступ под admin без знания пароля, дальнейшая эскалация.

2. `GET /search?q=...`
- Параметр: `q` (строка поиска по имени товара).
- Назначение: поиск продуктов по имени.
- Риск: UNION‑инъекция, извлечение конфиденциальных данных (логины, пароли, e‑mail).

### 2. Vulnerability Description (CWE и CVSS)

**Основная уязвимость:** Improper neutralization of input in SQL queries (classic SQL Injection).

- **CWE‑89**: Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection).
- **Типы в модуле:**
- Authentication bypass через `OR '1'='1'--` (in‑band).
- UNION-based SQL injection для извлечения данных.
- sqlmap дополнительно находит boolean‑based blind и time‑based варианты на параметре `username`.

**CVSS v3.1 (типичный для SQLi такого уровня):**

- **Base Score**: 9.8 (Critical).
- **Vector**: `AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H`.
- Обоснование:
- Доступ по сети, без аутентификации.
- Низкая сложность атаки (простые payload‑ы).
- Нет необходимости во взаимодействии пользователя.
- Полная конфиденциальность, целостность и доступность под угрозой.

### 3. Database Schema and Data

**Схема БД:** `public` в БД `vuln`.

- Таблица `users`:
- `id SERIAL PRIMARY KEY`
- `username VARCHAR(100) NOT NULL`
- `password VARCHAR(100) NOT NULL`
- `email VARCHAR(100)`
- `is_admin BOOLEAN DEFAULT FALSE`

Примеры записей:
- `admin / supersecret123 / admin@example.com / true`
- `user1 / password123 / user1@example.com / false`
- `bob / bob123 / bob@example.com / false`

- Таблица `products`:
- `id SERIAL PRIMARY KEY`
- `name VARCHAR(200)`
- `price VARCHAR(50)`
- `description TEXT`

Примеры записей:
- `Laptop / $999 / High-performance laptop`
- `Mouse / $25 / Wireless mouse`
- `Keyboard / $75 / Mechanical keyboard`

Скрипт создания и заполнения находится в `database/init.sql` и автоматически исполняется контейнером PostgreSQL при первом запуске.

### 4. Implementation Details

#### 4.1. Уязвимый логин `/login`

- Обработчик читает `username` и `password` из формы (`POST`).
- Формируется SQL‑запрос через `fmt.Sprintf` с прямой подстановкой значений:

```go
query := fmt.Sprintf(
"SELECT id, username, email, is_admin FROM users WHERE username='%s' AND password='%s'",
username, password,
)
```


- Затем этот запрос передаётся в `db.Query`, без параметров и без экранирования.

**Результат:**

- Атакующий может закрыть строковый литерал и добавить своё условие:

```sql
username=admin' OR '1'='1'--
```


- Это приводит к выполнению запроса без проверки реального пароля и входу под пользователем `admin`.

#### 4.2. Уязвимый поиск `/search`

- Обработчик читает `q` из строки запроса (`GET /search?q=...`).
- Формируется запрос вида:


```go
query := fmt.Sprintf(
"SELECT id, name, price, description FROM products WHERE name LIKE '%%%s%%'",
searchQuery,
)
```


- Поскольку `searchQuery` вставляется без экранирования, злоумышленник может закрыть строку и добавить `UNION SELECT` с нужными полями.

**Пример опасного payload:**

```sql
' UNION SELECT id, username, password, email FROM users--
```


Этот payload позволяет вывести строки из таблицы `users` в таблицу результатов поиска, так как количество и типы колонок совпадают.

### 5. Exploitation Steps

#### 5.1. Manual Exploits (curl / браузер)

**Authentication bypass:**

1. Открыть `http://localhost:8081/login`.
2. Ввести:
    - Username: `admin' OR '1'='1'--`
    - Password: любое значение.
3. Сервер выполняет SQL с условием, которое всегда истинно, и возвращает успешный вход под `admin`.

Аналогично атака повторяется через `curl` в скрипте `01_auth_bypass.sh`.

**UNION-based SQL Injection:**

1. Открыть `http://localhost:8081/search`.
2. В поле поиска ввести:

```sql
' UNION SELECT id, username, password, email FROM users--
```


3. Обработчик выполнит объединённый запрос и выведет строки из `users` вместо обычных продуктов.

В скрипте `02_union_injection.sh` используется `curl` с `--data-urlencode`, чтобы корректно закодировать payload в URL.

#### 5.2. Automated Exploits (sqlmap)

На эндпоинте `POST /login` с телом `username=admin&password=test` sqlmap обнаруживает:

- boolean‑based blind injection.
- stacked queries (через `;SELECT PG_SLEEP(5)--`).
- time‑based blind injection через `PG_SLEEP`.

Дальше sqlmap автоматически:

1. Определяет СУБД как PostgreSQL.
2. Находит доступные схемы, в том числе `public`.
3. Считывает список таблиц: `users`, `products`.
4. Выполняет `--dump` таблицы `public.users`, извлекая все логины, пароли, e‑mail и флаг `is_admin`.

---

## IV. Difficulties and New Skills

### Difficulties

- Корректная работа sqlmap с эндпоинтом, который возвращает HTTP 401 при неуспешной аутентификации:
- пришлось использовать флаг `--ignore-code=401`;
- важно было обеспечить различимое поведение сервера для валидных и инъецированных запросов, чтобы sqlmap мог детектировать SQLi.
- Необходимо было учитывать особенности PostgreSQL:
- таблицы находятся в схеме `public`, поэтому при работе sqlmap нужно явно указывать `-D public`, иначе он "не видит" таблицы.

### New Skills

- Практический опыт написания заведомо уязвимого кода на Go с использованием `database/sql` и драйвера PostgreSQL.
- Практическая работа с sqlmap:
- выбор параметров (`--data`, `--ignore-code`, `-D`, `-T`, `--dump`);
- интерпретация логов (типы инъекций, найденные payload‑ы).
- Глубже понял разницу между:
- классической (in‑band) SQL‑инъекцией;
- boolean-based blind;
- time-based blind (через задержки).

---

## V. Conclusion and Judgement

В рамках модуля classic SQL injection удалось построить полностью рабочий proof‑of‑concept: уязвимое веб‑приложение на Go, реальная БД PostgreSQL в Docker, понятные сценарии атак и автоматизированные эксплойты.

Один и тот же программный дефект — конкатенация строк в SQL‑запросах без параметризации — даёт возможность как простого обхода логина и вытаскивания паролей через UNION, так и более продвинутых blind‑атак, которые находят и эксплуатируют sqlmap. Это хорошо демонстрирует, насколько критична правильно реализованная работа с БД и почему классическая SQL‑инъекция (CWE‑89, CVSS 9.8) до сих пор остаётся одной из самых опасных уязвимостей в веб‑приложениях.
