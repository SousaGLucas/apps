begin;

create table if not exists accounts (
    id uuid primary key,
    user_id uuid not null,
    ledger_account_id uuid not null,
    created_at timestamptz not null
);

commit;
