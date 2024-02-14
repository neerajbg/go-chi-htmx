CREATE TABLE IF NOT EXISTS public.posts
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description text COLLATE pg_catalog."default"
)