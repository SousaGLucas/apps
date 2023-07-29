begin;

create table if not exists users (
    id uuid primary key,
    name text not null,
    email text not null unique,
    hashPassword text not null,
    created_at timestamptz not null,
    updated_at timestamptz
);

commit;
