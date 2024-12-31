CREATE TABLE public.transactions (
    id BIGSERIAL PRIMARY KEY, 
    sof_number VARCHAR(100) NOT NULL, 
    dof_number VARCHAR(100) NOT NULL, 
    amount BIGINT, 
    transaction_type CHAR(1) NOT NULL CHECK (transaction_type IN ('C', 'D')), 
    wallet_id BIGINT NOT NULL, 
    transaction_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (wallet_id) REFERENCES public.wallets (id) ON DELETE CASCADE
);
