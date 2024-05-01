from sqlalchemy import create_engine, MetaData, Table, select
import pandas as pd
import os
from dotenv import load_dotenv

load_dotenv()

def get_ratings():
    # Получение переменных окружения
    user = os.getenv('POSTGRES_USER')
    password = os.getenv('POSTGRES_PASSWORD')
    db = os.getenv('POSTGRES_DB')
    host = os.getenv('POSTGRES_HOST')

    database_url = f'postgresql://{user}:{password}@{host}:5432/{db}'
    engine = create_engine(database_url)

    metadata = MetaData()
    metadata.reflect(bind=engine)  # Загружаем информацию о таблицах

    # Доступ к таблице video_estimation после загрузки метаданных
    video_estimation = metadata.tables['video_estimation']

    # Выборка необходимых столбцов
    query = select(
        video_estimation.c.user_id,
        video_estimation.c.video_id,
        video_estimation.c.rate
    )

    # Выполнение запроса и загрузка результатов в DataFrame
    with engine.connect() as conn:
        result = conn.execute(query)
        df = pd.DataFrame(result.fetchall(), columns=result.keys())

    return df