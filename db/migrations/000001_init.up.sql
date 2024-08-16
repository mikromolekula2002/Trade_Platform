CREATE TABLE IF NOT EXISTS "user_data" (
    login VARCHAR(20) PRIMARY KEY NOT NULL,
    name VARCHAR(50),
    first_name VARCHAR(50),
    phone_number VARCHAR(20),
    image_profile VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS "user_ads" (
    ads_id SERIAL PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    image_1 VARCHAR(255),
    image_2 VARCHAR(255),
    image_3 VARCHAR(255),
    ads_name VARCHAR(50),
    ads_description VARCHAR(200),
    ads_price FLOAT
);

CREATE TABLE IF NOT EXISTS "user" (
    login VARCHAR(20) UNIQUE PRIMARY KEY NOT NULL,
    hash_password BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS "likes" (
    user_login VARCHAR(20) NOT NULL,
    ads_id INT NOT NULL
);