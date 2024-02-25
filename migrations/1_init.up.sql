CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    email VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    patronymic VARCHAR NOT NULL,
    surname VARCHAR,
    sex VARCHAR,
    age INTEGER,
    birthday DATE
)