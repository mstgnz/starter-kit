-- USERS_COUNT
SELECT count(*) FROM users;

-- USERS_PAGINATE
select * from users where fullname ilike $1 or email ilike $1 or phone ilike $1 order by id desc offset $2 limit $3;

-- USER_EXISTS_WITH_ID
SELECT count(*) FROM users WHERE id=$1;

-- USER_EXISTS_WITH_EMAIL
SELECT count(*) FROM users WHERE email=$1;

-- USER_GET_WITH_ID
SELECT id, fullname, email, is_admin, password FROM users WHERE id=$1 AND deleted_at isnull;

-- USER_GET_WITH_EMAIL
SELECT id, fullname, email, is_admin, password FROM users WHERE email=$1 AND deleted_at isnull;

-- USER_INSERT
INSERT INTO users (fullname,email,password,phone) VALUES ($1,$2,$3,$4) RETURNING id,fullname,email,phone;

-- USER_UPDATE_PASS
UPDATE users SET password=$1, updated_at=$2 WHERE id=$3;

-- USER_LAST_LOGIN
UPDATE users SET last_login=$1 WHERE id=$2;

-- USER_DELETE
UPDATE users SET active=$1, deleted_at=$2, updated_at=$3 WHERE id=$4;