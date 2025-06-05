package db

import (
	"context"
	"database/sql"

	"github.com/Zkeai/DDPay/internal/model"
)

// GetMembershipTransactionByOrderID 根据订单号获取交易记录
func (db *DB) GetMembershipTransactionByOrderID(ctx context.Context, orderID string) (*model.MembershipTransaction, error) {
	query := `SELECT id, user_id, level_id, amount, transaction_type, payment_method, 
              status, order_id, created_at, updated_at 
              FROM membership_transactions WHERE order_id = ?`
	
	transaction := &model.MembershipTransaction{}
	
	err := db.db.QueryRow(ctx, query, orderID).Scan(
		&transaction.ID, &transaction.UserID, &transaction.LevelID,
		&transaction.Amount, &transaction.TransactionType, &transaction.PaymentMethod,
		&transaction.Status, &transaction.OrderID, &transaction.CreatedAt, &transaction.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有找到对应的交易记录
		}
		return nil, err
	}
	
	return transaction, nil
} 