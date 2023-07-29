-- name: CreateAccount :exec
insert into accounts (id, user_id, ledger_account_id, created_at)
values (
        @id,
        @user_id,
        @ledger_account_id,
        @created_at
);

-- name: GetAccount :one
select *
from accounts
where id = @id;
