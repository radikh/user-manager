CREATE SEQUENCE admins_code_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE admins
(
    admin_id integer DEFAULT nextval('admins_code_id_seq'
    ::regclass) NOT NULL,
    admin varchar
    (20) NOT NULL,
    password varchar
    (80) NOT NULL,
    CONSTRAINT admins_pk PRIMARY KEY
    (admin_id)
);

    GRANT SELECT, INSERT, DELETE, UPDATE  ON admins TO um_user;

    INSERT INTO admins
        (admin, password)
    VALUES
        ('admin', '$argon2id$v=19$m=65536,t=3,p=1$L/cXOPSeeE9f68JKienFug$8tkLKbOpaOnxPniURMxUsg');