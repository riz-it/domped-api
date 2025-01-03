-- Active: 1735301424776@@aws-0-ap-southeast-1.pooler.supabase.com@6543@postgres@public
CREATE TABLE public.topup (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
