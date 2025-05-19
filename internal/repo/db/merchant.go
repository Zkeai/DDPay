package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Zkeai/DDPay/internal/model"
)

func (db *DB) InsertMerchantWallet(ctx context.Context, wallet *model.MerchantWallet) error {
	query := `
		INSERT INTO merchant_wallets (merchant_id, chain, address, derivation_path)
		VALUES (?, ?, ?, ?)
	`
	_, err := db.db.Exec(ctx, query, wallet.MerchantID, wallet.Chain, wallet.Address, wallet.DerivationPath)
	return err
}

func (db *DB) GetWalletByMerchantAndChain(ctx context.Context, merchantID uint32, chain string) (*model.MerchantWallet, error) {
	query := `
		SELECT id, merchant_id, chain, address, derivation_path, created_at
		FROM merchant_wallets
		WHERE merchant_id = ? AND chain = ?
		LIMIT 1
	`

	row := db.db.QueryRow(ctx, query, merchantID, chain)

	var wallet model.MerchantWallet
	err := row.Scan(
		&wallet.ID,
		&wallet.MerchantID,
		&wallet.Chain,
		&wallet.Address,
		&wallet.DerivationPath,
		&wallet.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 查不到返回 nil
		}
		return nil, err
	}
	return &wallet, nil
}
