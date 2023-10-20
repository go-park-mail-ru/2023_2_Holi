# 2023_2_Holi
Backend репозиторий команды Holi

## Объяснение нормальных форм(нф):

- 1я нф - т.к. каждый столбец содержит в себе атомарное значение
- 2я нф -  т.к. нет составных ключей
- 3я нф - т.к. каждый столбец таблицы зависит только от PK (нет зависимостей между неключевыми атрибутами)
- НФБК - т.к. нет составных потенциальных ключей

## Функциональные зависимости

Relation VIDEO:

{id} -> name, description, duration, preview_path, media_path, release_path, rating, age_restriction, seasons_count

Relation EPISODE:

{id} -> name, description, duration, preview_path, number

Relation USER:

{id} -> name, email, password, date_joined, image_path

Relation CAST:

{id} -> name

Relation TAG:

{id} -> name

Relation GENRE:

{id} -> name

Relation VIDEO_ESTIMATION:

{user_id, video_id} -> rate


## ER

```mermaid
---
title: Netflix
---
erDiagram
    VIDEO {
        _ id PK
        _ name "not null"
        _ description
        _ duration "not null"
        _ preview_path "not null"
        _ media_path "not null"
        _ release_year
        _ rating
        _ age_restriction
        _ seasons_count "not null default 0"
        _ created_at "default current_timestamp not null"
        _ updated_at "default current_timestamp not null"
    }
    
    VIDEO ||--|{ EPISODE: contains
    EPISODE {
        _ id PK
        _ name "not null"
        _ description
        _ duration "not null"
        _ preview_path "not null"
        _ media_path "not null"
        _ number "not null"
        _ season_number "not null"
        _ created_at "default current_timestamp not null"
        _ updated_at "default current_timestamp not null"
        _ video_id FK
    }

    CAST {
        _ id PK
        _ name "not null unqiue"
    }

    VIDEO-CAST ||--|{ VIDEO: video
    VIDEO-CAST ||--|{ CAST: cast
    VIDEO-CAST {
        _ video_id FK
        _ cast_id FK
    }
    
    GENRE {
        _ id PK
        _ name "not null unque"
    }

    VIDEO-GENRE ||--|{ VIDEO: video 
    VIDEO-GENRE ||--|{ GENRE: genre
    VIDEO-GENRE {
        _ video_id FK
        _ genre_id FK
    }

    USER {
        _ id PK
        _ name
        _ email "not null unique"
        _ password "not null"
        _ date_joined "default now"
        _ image_path
        _ created_at "default current_timestamp not null"
        _ updated_at "default current_timestamp not null"
    }
    
    VIDEO_ESTIMATION ||--o{ VIDEO: video 
    VIDEO_ESTIMATION ||--o{ USER: user
    VIDEO_ESTIMATION {
        _ user_id FK
        _ video_id FK
        _ rate "not null"
        _ created_at "default current_timestamp not null"
        _ updated_at "default current_timestamp not null"
    }

    TAG {
        _ id PK
        _ name "not null unique"
    }

    VIDEO-TAG ||--|{ VIDEO: video
    VIDEO-TAG ||--|{ TAG: tag
    VIDEO-TAG {
        _ video_id FK
        _ tag_id FK
    }
    

```
