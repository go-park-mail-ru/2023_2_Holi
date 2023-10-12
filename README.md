# 2023_2_Holi
Backend репозиторий команды Holi

---
title: Netflix
---
erDiagram
    VIDEO {
        _ id PK
        _ name 
        _ description
        _ duration
        _ preview_path
        _ media_path
        _ release_year
        _ rating
        _ age_restriction
        _ seasons_count
    }
    
    VIDEO ||--|{ EPISODE: contains
    EPISODE {
        _ id PK
        _ name 
        _ description
        _ duration
        _ preview_path
        _ number
        _ season_number
        _ video_id FK
    }

    CAST {
        _ id PK
        _ name
    }

    VIDEO-CAST ||--|{ VIDEO: video
    VIDEO-CAST ||--|{ CAST: cast
    VIDEO-CAST {
        _ video_id FK
        _ cast_id FK
    }
    
    GENRE {
        _ id PK
        _ name
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
        _ email
        _ password
        _ date_joined
        _ image_path
    }
    
    VIDEO_ESTIMATION ||--o{ VIDEO: video 
    VIDEO_ESTIMATION ||--o{ USER: user
    VIDEO_ESTIMATION {
        _ user_id FK
        _ video_id FK
        _ rate
    }

