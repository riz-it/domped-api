-- Active: 1735301424776@@aws-0-ap-southeast-1.pooler.supabase.com@6543@postgres@public
CREATE TABLE public.users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    full_name VARCHAR(125),
    phone VARCHAR(16),
    email VARCHAR(75) NOT NULL,
    password VARCHAR(255) NOT NULL,
    hashed_rt TEXT,
    is_active BOOLEAN NOT NULL DEFAULT false,
    email_verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
