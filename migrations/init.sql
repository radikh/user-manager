CREATE TABLE public.users
(
    id uuid NOT NULL,
    user_name text NOT NULL,
    "password" text NOT NULL,
    email text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    phone text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NULL,
    CONSTRAINT users_pk PRIMARY KEY (id)
);
ALTER TABLE public.users ADD CONSTRAINT users_un_unigue UNIQUE (user_name);
ALTER TABLE public.users ADD CONSTRAINT users_email_unigue UNIQUE (email);
ALTER TABLE public.users ADD CONSTRAINT users_phone_unigue UNIQUE (phone);
CREATE USER gateway
WITH PASSWORD $UM_PASSWORD;
GRANT SELECT, INSERT, UPDATE  ON users TO gateway;