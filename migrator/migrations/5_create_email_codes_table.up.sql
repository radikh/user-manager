CREATE TABLE public.email_codes
(
    id uuid NOT NULL,
    user_name varchar(20) NOT NULL,
    email varchar(70) NOT NULL,
    verification_code varchar (50) NOT NULL,
    created_at timestamp NOT NULL,
    /*CONSTRAINT codes_pk PRIMARY KEY (id)*/
);
