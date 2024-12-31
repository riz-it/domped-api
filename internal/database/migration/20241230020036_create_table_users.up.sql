-- Active: 1735609545514@@pg-174dd6a0-mrizalf-040e.g.aivencloud.com@27154@defaultdb
CREATE TABLE public.users (
    id SERIAL NOT NULL,
    full_name VARCHAR(125),
    email VARCHAR(125) NOT NULL,
    password VARCHAR(255) NOT NULL,
    hashed_rt TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    email_verified_at TIMESTAMP(3),
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
