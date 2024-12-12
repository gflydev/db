-- -----------------------------------------------------
-- Table users
-- -----------------------------------------------------
CREATE TYPE user_status AS ENUM ('pending', 'active', 'blocked');

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR (255) NOT NULL UNIQUE,
                       password VARCHAR (255) NOT NULL,
                       fullname VARCHAR (255) NULL,
                       phone VARCHAR(20) NULL,
                       token VARCHAR (100) NULL,
                       status user_status DEFAULT 'pending',
                       avatar VARCHAR (255) NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NULL,
                       verified_at TIMESTAMP NULL,
                       blocked_at TIMESTAMP NULL,
                       deleted_at TIMESTAMP NULL,
                       last_access_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_users ON users (id);
CREATE UNIQUE INDEX email_users ON users (email ASC);

-- --------------------------------------------------------------------------------------
-- ------------------------------------ Initial data ------------------------------------
-- --------------------------------------------------------------------------------------
-- P@seWor9  ===>  $2a$04$9QD944312deeQjnxF.zNauGx7NQ0GtS.xJhLy.zWqWxOE8B/XCN9i
INSERT INTO users (email, password, fullname, phone, token, status, avatar, created_at, updated_at)
VALUES ('admin@gfly.dev', '$2a$04$9QD944312deeQjnxF.zNauGx7NQ0GtS.xJhLy.zWqWxOE8B/XCN9i', 'Admin', '0989831911', null, 'active', 'https://www.gfly.dev/assets/avatar.png', '2024-05-15 13:07:48.888668 +07:00', '2024-05-15 13:07:48.888668 +07:00');
