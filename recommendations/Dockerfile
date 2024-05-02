# Использование официального образа Python как базового
FROM python:3.9

# Установка рабочей директории в контейнере
WORKDIR /app

# Копирование файла зависимостей в рабочую директорию
COPY recommendations/requirements.txt .

# Установка зависимостей
RUN pip install --no-cache-dir -r requirements.txt

# Копирование остальных файлов проекта в контейнер
COPY /recommendations .

# Команда для запуска приложения
CMD ["python", "./item_based.py"]