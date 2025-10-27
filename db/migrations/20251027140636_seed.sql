-- +goose Up
-- +goose StatementBegin
-- админ
INSERT INTO "user" (username, email, password_hash, role) VALUES 
('admin', 'admin@roadmap.sh', '', 'admin');

-- категории
INSERT INTO category (name, description) VALUES 
('Design', 'Design related roadmaps and learning paths'),
('Frontend', 'Frontend development technologies and frameworks'),
('Backend', 'Backend development and server-side technologies'),
('Data Science', 'Data analysis, machine learning and AI'),
('DevOps', 'Infrastructure, deployment and operations'),
('UI/UX', 'User interface and user experience design');

-- roadmap_info для категории Design
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Design System',
    'Component Libraries, Design Tokens, Documentation, Accessibility, Versioning',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Design';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'UI Design',
    'Visual Design, Layout, Typography, Color Theory, Design Tools',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Design';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Graphic Design',
    'Brand Identity, Illustration, Print Design, Motion Graphics, Adobe Creative Suite',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Design';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Product Design',
    'User Research, Prototyping, Design Thinking, Product Strategy, Design Sprints',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Design';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Motion Design',
    'Animation Principles, After Effects, Lottie, Microinteractions, Visual Storytelling',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Design';

-- roadmap_info для категории Frontend
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'React',
    'Components, Hooks, State Management, React Router, Testing',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Frontend';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Vue',
    'Vue 3, Composition API, Vue Router, Pinia, Vue Testing',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Frontend';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Angular',
    'Components, Services, RxJS, NgRx, Angular CLI',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Frontend';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Next.js',
    'SSR, SSG, API Routes, App Router, Deployment',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Frontend';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'TypeScript',
    'Type System, Generics, Decorators, Configuration, Advanced Types',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Frontend';

-- roadmap_info для категории Backend
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Node.js',
    'Express, REST APIs, Middleware, Authentication, Performance',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Backend';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Python Backend',
    'Django, FastAPI, Flask, Database ORM, Async Programming',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Backend';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Go',
    'Goroutines, Standard Library, Web Frameworks, Concurrency, Performance',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Backend';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Spring Boot',
    'Dependency Injection, Spring MVC, Data JPA, Security, Microservices',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Backend';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'API Design',
    'REST, GraphQL, OpenAPI, Rate Limiting, Versioning',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Backend';

-- roadmap_info для категории Data Science
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Machine Learning',
    'Supervised Learning, Unsupervised Learning, Neural Networks, Model Evaluation, Feature Engineering',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Data Science';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Data Analyst',
    'SQL, Data Visualization, Statistical Analysis, Business Intelligence, Reporting',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Data Science';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Data Engineer',
    'ETL, Data Pipelines, Apache Spark, Data Warehousing, Big Data',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Data Science';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'AI Engineer',
    'Deep Learning, NLP, Computer Vision, Model Deployment, MLOps',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Data Science';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Prompt Engineering',
    'LLM Techniques, RAG, Fine-tuning, Context Management, Evaluation',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'Data Science';

-- roadmap_info для категории DevOps
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Kubernetes',
    'Pods, Services, Deployments, Helm, Cluster Management',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'DevOps';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Docker',
    'Containers, Images, Dockerfile, Docker Compose, Container Orchestration',
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
    'Infrastructure as Code, Providers, Modules, State Management, Best Practices',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'DevOps';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Cyber Security',
    'Network Security, Application Security, Threat Modeling, Incident Response, Compliance',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'DevOps';

-- roadmap_info для категории UI/UX
INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'UX Design',
    'User Research, Personas, Journey Mapping, Usability Testing, Information Architecture',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Interaction Design',
    'Prototyping, User Flows, Microinteractions, Animation, Feedback Systems',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'User Research',
    'Interviews, Surveys, Analytics, A/B Testing, Competitive Analysis',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Information Architecture',
    'Content Strategy, Navigation Design, Taxonomy, Sitemaps, Card Sorting',
    true
FROM "user" u, category c WHERE u.username = 'admin' AND c.name = 'UI/UX';

INSERT INTO roadmap_info (author_id, category_id, name, description, is_public)
SELECT 
    u.id,
    c.id,
    'Accessibility',
    'WCAG, ARIA, Screen Readers, Color Contrast, Keyboard Navigation',
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
