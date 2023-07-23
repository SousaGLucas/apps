begin;

create table if not exists accounts (
    id uuid primary key,
    created_at timestamptz not null
);

create table if not exists account_events (
    id uuid primary key,
    account_id uuid not null references accounts(id),
    type text not null check (type in ('cash_in', 'cash_out')),
    amount int not null,
    created_at timestamptz not null
);

commit;
