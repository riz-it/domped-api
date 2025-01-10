-- Active: 1735301424776@@aws-0-ap-southeast-1.pooler.supabase.com@6543@postgres@public
CREATE TABLE public.topup (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL, 
    amount BIGINT NOT NULL,
    status INT, 
    snap_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES public.users (id) ON DELETE CASCADE
);