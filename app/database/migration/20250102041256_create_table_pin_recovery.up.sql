-- Active: 1735301424776@@aws-0-ap-southeast-1.pooler.supabase.com@6543@postgres@public
CREATE TABLE public.pin_recoveries (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    pin_code VARCHAR(255) NOT NULL,
    wallet_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (wallet_id) REFERENCES public.wallets (id) ON DELETE CASCADE
);
