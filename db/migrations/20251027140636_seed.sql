-- +goose Up
-- +goose StatementBegin
-- администратор
INSERT INTO "user" (username, email, password_hash, role) VALUES 
('admin', 'admin@roadmap.sh', '', 'admin');

-- бот
INSERT INTO "user" (id, username, email, password_hash, avatar_url) VALUES 
('11111111-1111-1111-1111-111111111111', 'bot', 'bot@roadmap.sh', '', 'http://127.0.0.1:9000/avatars/bot.jpg');

-- категории
INSERT INTO category (name, description) VALUES 
('Дизайн', 'Интерактивный, визуальный и продуктовый дизайн'),
('Фронтенд', 'Технологии и фреймворки фронтенд-разработки'),
('Бэкенд', 'Бэкенд-разработка и серверные технологии'),
('Наука о данных', 'Анализ данных, машинное обучение и ИИ'),
('DevOps', 'Инфраструктура, развертывание и эксплуатация'),
('UI/UX', 'Дизайн пользовательского интерфейса и пользовательского опыта');

-- roadmap_info для категории Дизайн
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Системный дизайн',
    'Библиотеки компонентов, Токены дизайна, Документация, Доступность, Версионирование',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Дизайн';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Дизайн UI',
    'Визуальный дизайн, Компоновка, Типографика, Теория цвета, Инструменты дизайна',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Дизайн';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Графический дизайн',
    'Фирменный стиль, Иллюстрация, Полиграфический дизайн, Моушн-графика, Adobe Creative Suite',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Дизайн';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Дизайн продукта',
    'Исследование пользователей, Прототипирование, Дизайн-мышление, Стратегия продукта, Дизайн-спринты',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Дизайн';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Анимация',
    'Принципы анимации, After Effects, Lottie, Микроинтеракции, Визуальное повествование',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Дизайн';

-- roadmap_info для категории Фронтенд
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'React',
    'Компоненты, Хуки, Управление состоянием, React Router, Тестирование',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Фронтенд';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Vue',
    'Vue 3, Composition API, Vue Router, Pinia, Тестирование Vue',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Фронтенд';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Angular',
    'Компоненты, Сервисы, RxJS, NgRx, Angular CLI',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Фронтенд';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Next.js',
    'SSR, SSG, API Routes, App Router, Развертывание',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Фронтенд';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'TypeScript',
    'Система типов, Дженерики, Декораторы, Конфигурация, Продвинутые типы',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Фронтенд';

-- roadmap_info для категории Бэкенд
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Проектирование высоконагруженных систем',
    'Модели параллелизма, Управление памятью, Пул соединений, Профилирование, Оптимизация',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Бэкенд';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Python Бэкенд',
    'Django, FastAPI, Flask, ORM базы данных, Асинхронное программирование',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Бэкенд';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Go',
    'Горутины, Стандартная библиотека, Веб-фреймворки, Параллелизм, Производительность',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Бэкенд';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Spring Boot',
    'Внедрение зависимостей, Spring MVC, Data JPA, Безопасность, Микросервисы',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Бэкенд';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Дизайн API',
    'REST, GraphQL, OpenAPI, Ограничение частоты запросов, Версионирование',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Бэкенд';

-- roadmap_info для категории Наука о данных
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Машинное обучение',
    'Обучение с учителем, Обучение без учителя, Нейронные сети, Оценка моделей, Формирование признаков',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Наука о данных';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Анализ данных',
    'SQL, Визуализация данных, Статистический анализ, Бизнес-аналитика, Отчетность',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Наука о данных';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Data Engineering',
    'ETL, Data Pipeline, Apache Spark, Хранилища данных, Большие данные',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Наука о данных';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'AI Engineering',
    'Глубокое обучение, NLP, Компьютерное зрение, Развертывание моделей, MLOps',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Наука о данных';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Prompt Engineering',
    'Техники LLM, RAG, Тонкая настройка, Управление контекстом, Оценка',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Наука о данных';

-- roadmap_info для категории DevOps
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Kubernetes',
    'Поды, Сервисы, Развертывания, Helm, Управление кластерами',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'DevOps';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Docker',
    'Контейнеры, Образы, Dockerfile, Docker Compose, Оркестрация контейнеров',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'DevOps';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'AWS',
    'EC2, S3, Lambda, RDS, CloudFormation, IAM',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'DevOps';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Terraform',
    'Инфраструктура как код, Провайдеры, Модули, Управление состоянием, Лучшие практики',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'DevOps';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Кибербезопасность',
    'Сетевая безопасность, Безопасность приложений, Моделирование угроз, Реагирование на инциденты, Соответствие требованиям',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'DevOps';

-- roadmap_info для категории UI/UX
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'UX дизайн',
    'Исследование пользователей, Карты пути пользователя, Юзабилити-тестирование, Информационная архитектура',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Интерактивный дизайн',
    'Прототипирование, Пользовательские сценарии, Микроинтеракции, Анимация, Системы обратной связи',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Пользовательский опыт',
    'Интервью, Опросы, Аналитика, A/B-тестирование, Сравнительный анализ',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Стратегии продвижения контента',
    'Контент-стратегия, Дизайн навигации, Карты сайта, Карточная сортировка',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Доступность',
    'WCAG, ARIA, Скринридеры, Контрастность цветов, Навигация с клавиатуры',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаляем все созданные данные в обратном порядке
TRUNCATE TABLE roadmap_info CASCADE;
TRUNCATE TABLE category CASCADE;
TRUNCATE TABLE "user" CASCADE;
-- +goose StatementEnd
