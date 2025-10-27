-- name: StoreRefreshToken :exec
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES(
	$1,
	NOW(),
	NOW(),
	$2,
	NOW() + INTERVAL '60 days'
);

-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens WHERE token = $1;

-- name: RefreshTokenExpired :one
SELECT (expires_at < NOW() OR revoked_at IS NOT NULL) FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(),
	revoked_at = NOW()
WHERE token = $1;
