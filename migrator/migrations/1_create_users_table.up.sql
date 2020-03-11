CREATE TABLE public.users
(
    id uuid NOT NULL,
    username text NOT NULL,
    "password" text NOT NULL,
    email text NOT NULL,
    firstname text NOT NULL,
    lastname text NOT NULL,
    phone text NOT NULL,
    createdat timestamp NOT NULL,
    updatedat timestamp NULL,
    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_un UNIQUE (username)
);
