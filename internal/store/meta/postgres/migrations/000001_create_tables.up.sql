-- пользователи
CREATE TABLE IF NOT EXISTS users (
   user_id serial PRIMARY KEY,
   login VARCHAR (100) UNIQUE NOT NULL,
   password VARCHAR (100) NOT NULL
);

-- типы возможных секретов
CREATE TYPE secret_type AS ENUM (
   'TEXT',
   'LOGIN_PASSWORD',
   'CREDIT_CARD',
   'FILE'
);

-- мета-данные секретов
CREATE TABLE IF NOT EXISTS meta (
   meta_id SERIAL PRIMARY KEY
   user_id INT,
   meta_unique_key VARCHAR(100) UNIQUE NOT NULL, 
   alias VARCHAR(100),
   type secret_type, 
   extra text,
   CONSTRAINT fk_user
      FOREIGN KEY (user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE
);
