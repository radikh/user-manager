CREATE SEQUENCE codes_code_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE codes
(
    code_id integer DEFAULT nextval('codes_code_id_seq'
    ::regclass) NOT NULL,
    user_name varchar
    (20) NOT NULL,
    code varchar
    (64) NOT NULL,
    expired_at timestamp NULL,
    active boolean NOT NULL DEFAULT true,
    CONSTRAINT codes_pk PRIMARY KEY
    (code_id)
);

    GRANT SELECT, INSERT, DELETE, UPDATE  ON codes TO um_user;