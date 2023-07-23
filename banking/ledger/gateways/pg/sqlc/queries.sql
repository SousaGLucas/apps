-- name: CreateAccount :exec
insert into accounts (id, created_at)
values (
       @id,
       @created_at
);

-- name: CreateEvent :exec
insert into account_events (id, account_id, type, amount, created_at)
values (
        @id,
        @account_id,
        @type,
        @amount,
        @created_at
       );

-- name: GetAccount :one
select *
from accounts
where id = @id;

-- name: ListAccountEvents :many
select *
from account_events
where @account_id = @account_id
--         and id < @last_fetched_id
order by id desc
limit @page_size;
