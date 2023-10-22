# 2023_2_Holi

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

{id} -> name, email, password, date_joined, image_path, created_at

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
        _ name "NOT NULL"
        _ description
        _ duration "NOT NULL"
        _ preview_path "NOT NULL"
        _ media_path "NOT NULL"
        _ release_year
        _ rating
        _ age_restriction
        _ seasons_count "NOT NULL DEFAULT 0"
        _ created_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
        _ updated_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
    }
    
    VIDEO ||--|{ EPISODE: contains
    EPISODE {
        _ id PK
        _ name "NOT NULL"
        _ description
        _ duration "NOT NULL"
        _ preview_path "NOT NULL"
        _ media_path "NOT NULL"
        _ number "NOT NULL"
        _ season_number "NOT NULL"
        _ created_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
        _ updated_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
        _ video_id FK
    }

    CAST {
        _ id PK
        _ name "NOT NULL UNIQUE"
    }

    VIDEO_CAST ||--|{ VIDEO: video
    VIDEO_CAST ||--|{ CAST: cast
    VIDEO_CAST {
        _ video_id FK
        _ cast_id FK
        "UNIQUE (video_id, cast_id)"
    }
    
    GENRE {
        _ id PK
        _ name "NOT NULL UNIQUE"
    }

    VIDEO_GENRE ||--|{ VIDEO: video 
    VIDEO_GENRE ||--|{ GENRE: genre
    VIDEO_GENRE {
        _ video_id FK
        _ genre_id FK
        "UNIQUE (video_id, genre_id)"
    }

    USER {
        _ id PK
        _ name
        _ email "NOT NULL UNIQUE"
        _ password "NOT NULL"
        _ image_path
        _ created_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
        _ updated_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
    }
    
    VIDEO_ESTIMATION ||--o{ VIDEO: video 
    VIDEO_ESTIMATION ||--o{ USER: user
    VIDEO_ESTIMATION {
        _ rate "NOT NULL"
        _ created_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
        _ updated_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
        _ user_id FK
        _ video_id FK
        "UNIQUE (video_id, user_id)"
    }

    TAG {
        _ id PK
        _ name "NOT NULL UNIQUE"
    }

    VIDEO_TAG ||--|{ VIDEO: video
    VIDEO_TAG ||--|{ TAG: tag
    VIDEO_TAG {
        _ video_id FK
        _ tag_id FK
        "UNIQUE (video_id, tag_id)"
    }    

```
