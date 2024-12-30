-- Active: 1735301424776@@aws-0-ap-southeast-1.pooler.supabase.com@6543@postgres@public
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
