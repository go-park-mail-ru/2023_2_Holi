from sqlalchemy import create_engine
from sqlalchemy.sql import text
import os
from dotenv import load_dotenv
import redis

load_dotenv()


def save_recommendations(recommendations):
    user = os.getenv('POSTGRES_RECOM_USER')
    password = os.getenv('POSTGRES_RECOM_PASSWORD')
    db = os.getenv('POSTGRES_RECOM_DB')
    host = os.getenv('POSTGRES_RECOM_HOST')
    database_url = f'postgresql://{user}:{password}@{host}:5432/{db}'
    engine = create_engine(database_url)

    insert_query = text("INSERT INTO recommendations (user_id, movie_id) VALUES (:user_id, :movie_id)")

    with engine.connect() as conn:
        for user_id, user_recommendations in recommendations.items():
            for movie_id, _ in user_recommendations:
                conn.execute(insert_query, {'user_id': user_id, 'movie_id': movie_id})

        conn.commit()

    load_to_cache(recommendations)


def load_to_cache(recommendations):
    password = os.getenv('REDIS_PASSWORD')
    host = os.getenv('REDIS_HOST')
    port = os.getenv('REDIS_PORT')
    # r = redis.Redis(host=host, port=port, password=password)
    r = redis.Redis(host=host, port=port)

    # Проходим по всем рекомендациям
    for user_id, user_recommendations in recommendations.items():
        # Создаем ключ для пользователя
        key = f"recommendations:{user_id}"
        # print(user_recommendations)
        # Сначала очищаем существующие рекомендации в Redis для этого пользователя
        r.delete(key)

        for movie_id, _ in user_recommendations:
            r.rpush(key, movie_id)

        # Опционально, можно установить время жизни для ключа, чтобы он автоматически удалялся
        # например, через неделю (604800 секунд)
        r.expire(key, 604800)
