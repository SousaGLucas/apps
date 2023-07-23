-- name: CreateUser :exec
insert into users (id, name, email, hashPassword, created_at)
values (
        @id,
        @name,
        @email,
        @hashpassword,
        @created_at
       );

-- name: ListUsers :many
select *
from users;

-- name: GetUsers :one
select *
from users
where @id = @id;
