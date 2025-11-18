-- name: ListProducts :many
SELECT * FROM products;

-- name: FindProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: UpdateProductStock :exec
UPDATE products SET quantity = $1 WHERE id = $2;