CREATE SEQUENCE public.codes_code_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.codes
(
    code_id integer DEFAULT nextval('public.codes_code_id_seq'
    ::regclass) NOT NULL,
    user_name varchar
    (20) NOT NULL,
    code varchar
    (64) NOT NULL,
    expired_at timestamp NULL,
    active boolean NOT NULL DEFAULT true  
);
