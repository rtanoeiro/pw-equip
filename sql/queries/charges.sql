-- name: AddCharge :exec
INSERT INTO charges (
    user_email, payment_intent_id, amount, status
)
VALUES (?, ?, ?, ?);

-- name: GetChargesByUserEmail :many
SELECT
    payment_intent_id,
    amount,
    user_email
FROM charges
WHERE
    status = "succeeded"
    AND user_email = ?
ORDER BY created_at DESC;
