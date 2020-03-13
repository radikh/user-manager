ALTER TABLE public.users ADD CONSTRAINT users_un_unigue UNIQUE (user_name);
ALTER TABLE public.users ADD CONSTRAINT users_email_unigue UNIQUE (email);
ALTER TABLE public.users ADD CONSTRAINT users_phone_unigue UNIQUE (phone);