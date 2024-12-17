-- name: InsertMember :exec
INSERT INTO public.tb_amole_member
("name", email, "password", address)
VALUES($1, $2, $3, $4);

-- name: GetMemberById :one
SELECT id, name, email, password, address
FROM tb_amole_member WHERE id = $1;

-- name: GetMemberByEmail :one
SELECT id, name, email, password, address
FROM tb_amole_member WHERE email = $1;

