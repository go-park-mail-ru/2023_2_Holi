# 2023_2_Holi

## Объяснение нормальных форм(нф):

- 1я нф - т.к. каждый столбец содержит в себе атомарное значение
- 2я нф - т.к. нет зависимостей от части составных ключей
- 3я нф - т.к. каждый столбец таблицы зависит от PK (нет транзитивных зависимостей)
- НФБК - т.к. нет множественных составных потенциальных ключей

## Функциональные зависимости

Relation VIDEO:

{id} -> name, description, preview_path, release_year, rating, age_restriction, seasons_count

Relation EPISODE:

{id} -> name, description, duration, preview_path, media_path, number, season_number, video_id

Relation USER:

{id} -> name, email, password, date_joined, image_path, created_at  
{email} -> name, password, date_joined, image_path, created_at

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
        _ name
        _ description
        _ preview_path
        _ preview_video_path "NOT NULL"
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
        "PK (video_id, cast_id)"
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
        "PK (video_id, genre_id)"
    }

    USER {
        _ id PK
        _ name
        _ email "NOT NULL UNIQUE"
        _ password "NOT NULL UNIQUE"
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
        "PK (video_id, user_id)"
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
        "PK (video_id, tag_id)"
    }

    FAVOURITE ||--|{ VIDEO: video
    FAVOURITE ||--|{ USER: user
    FAVOURITE {
        _ created_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
        _ updated_at "DEFAULT CURRENT_TIMESTAMP NOT NULL"
        _ video_id FK
        _ user_id FK
        "PK (video_id, user_id)"
    }

```
