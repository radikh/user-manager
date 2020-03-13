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
