CREATE TABLE public.users
(
    id uuid NOT NULL,
    user_name varchar(20) NOT NULL,
    "password" varchar(64) NOT NULL,
    email varchar(70) NOT NULL,
    first_name varchar(20) NOT NULL,
    last_name varchar(50) NOT NULL,
    phone varchar(16) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NULL,
    CONSTRAINT users_pk PRIMARY KEY
    (id),
    CONSTRAINT users_un_unigue UNIQUE
    (user_name),
    CONSTRAINT users_email_unigue UNIQUE
    (email),
    CONSTRAINT users_phone_unigue UNIQUE
    (phone)
);