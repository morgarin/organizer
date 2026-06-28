# План по построению приложухи 


## Структура

organizer/
├── cmd/
│   ├── api/                # точка входа REST API
│   │   └── main.go
│   └── worker/             # точка входа для фоновых воркеров (парсинг, нотификации)
│       └── main.go
├── internal/
│   ├── auth/               # JWT, хэширование, middleware
│   ├── handlers/           # обработчики HTTP
│   ├── models/             # структуры БД
│   ├── repository/         # слой работы с БД
│   ├── service/            # бизнес-логика (парсинг, уведомления, планировщик)
│   ├── parser/             # парсер расписания (запросы к сайту, goquery)
│   ├── notifier/           # отправка в Telegram, webpush
│   └── config/             # загрузка .env, структура конфига
├── pkg/                    # переиспользуемые утилиты (logger, db, redis, queue)
├── web/
│   ├── static/             # CSS, JS, изображения
│   └── templates/          # html-шаблоны
├── migrations/             # SQL-миграции (golang-migrate)
├── docker-compose.yml
├── Dockerfile              # для api и worker (под вопросом)
├── Makefile                # команды: run, build, migrate, test (под вопросом)
├── go.mod / go.sum         # ну это я бля объяснять уж не буду
├── .gitignore
└── .env.example


## bd
### Удаление имитации bd для докера
docker stop organizer-postgres
docker rm organizer-postgres
### Запуск bd
docker start organizer-postgres

### Структура bd

-- юзеры
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    telegram_id BIGINT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- расписание
CREATE TABLE IF NOT EXISTS schedule (
    id SERIAL PRIMARY KEY,
    day_of_week INT NOT NULL,
    is_even BOOLEAN DEFAULT false,
    tasks_id VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()  
);

-- задачи юзеров
CREATE TABLE IF NOT EXISTS user_tasks (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    teacher VARCHAR(255),
    place VARCHAR(50),
    time_start TIME NOT NULL,
    time_end TIME NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()        
);

-- предметы расписания
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    teacher VARCHAR(255),
    place VARCHAR(50),
    time_start TIME NOT NULL,
    time_end TIME NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- чатГПТ(допилить в бетке)
CREATE TABLE IF NOT EXISTS chat_history (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


## Примерный расчет времени

1. Базовая инфраструктура

Настройка окружения	Docker для postgree, миграции, .env, подключение к БД и первичная обработка юзеров

2. Авторизация и личный кабинет

Базовая авторизация	Регистрация, логин, сессии/куки, bcrypt, middleware, страницы входа/регистрации (шаблоны)
Настройки пользователя	Профиль, привязка Telegram ID, время уведомления, страница настроек

3. Парсер расписания

Расписание (парсер PDF)	Скачивание PDF, разбор текста (с учетом структуры), сохранение в БД, ручной запуск, интеграция с cron

4. Фронтенд (базовый)

Отображение расписания	Таблица на фронте (Go template + HTMX), переключение недель, фильтрация по дню
Задачи (личные)	CRUD для задач, дашборд со списком, добавление/редактирование через HTMX, дедлайны

5. Уведомления

Уведомления (Telegram)	Создание бота, отправка сообщений, фоновая горутина/воркер для проверки событий, логика напоминаний
Написать планировщик (cron) для отправки напоминаний.

6. ChatGPT

Интеграция ChatGPT	Эндпоинт, вызов OpenAI API, сохранение истории, страница чата с контекстом

7. Тестирование и полировка

Логи, обработка ошибок, кеширование
Документация API (Swagger)
Тестирование безы


## Заметки

В auth_handler поправить на этапе фронта необходимость возврата имени юзера (ну тип нужна ли она вообще)
Надо запуск и мягкий стоп сервака засунуть в пакет server по возможности
Разобраться с моками и интерфесами для тестов