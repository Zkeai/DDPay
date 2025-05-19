package wallet

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"sync"
)

var (
	globalDeriver *WalletDeriver
	once          sync.Once
	initErr       error
)

// WalletDeriver 是用于根据助记词派生 EVM 地址的结构体
type WalletDeriver struct {
	mnemonic string
	wallet   *hdwallet.Wallet
}

// Config 表示配置文件中的助记词
type Config struct {
	Mnemonic string `yaml:"mnemonic"`
}

// DerivedAccount 包含派生出的地址及私钥
type DerivedAccount struct {
	Index      uint32
	Address    string
	PrivateKey string
}

// NewWalletDeriver 创建 WalletDeriver 实例
func NewWalletDeriver(mnemonic string) (*WalletDeriver, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("无效的助记词")
	}
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	return &WalletDeriver{
		mnemonic: mnemonic,
		wallet:   wallet,
	}, nil
}

// InitGlobalDeriver 初始化全局派生器
func InitGlobalDeriver(mnemonic string) error {
	once.Do(func() {
		var deriver *WalletDeriver
		deriver, initErr = NewWalletDeriver(mnemonic)
		if initErr == nil {
			globalDeriver = deriver
		}
	})
	return initErr
}

// GetGlobalDeriver 获取已初始化的派生器
func GetGlobalDeriver() (*WalletDeriver, error) {
	if globalDeriver == nil {
		return nil, errors.New("WalletDeriver 未初始化")
	}
	return globalDeriver, nil
}

// DeriveAccount 根据索引派生账户
func (wd *WalletDeriver) DeriveAccount(index uint32) (*DerivedAccount, error) {
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", index))
	account, err := wd.wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	privateKey, err := wd.wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	return &DerivedAccount{
		Index:      index,
		Address:    account.Address.Hex(),
		PrivateKey: fmt.Sprintf("0x%x", crypto.FromECDSA(privateKey)),
	}, nil
}

// DeriveRange 派生从 start 到 end（包含）的账户
func (wd *WalletDeriver) DeriveRange(start, end uint32) ([]DerivedAccount, error) {
	if end < start {
		return nil, errors.New("end 必须大于等于 start")
	}
	var accounts []DerivedAccount
	for i := start; i <= end; i++ {
		acc, err := wd.DeriveAccount(i)
		if err != nil {
			return nil, fmt.Errorf("派生 index=%d 失败: %w", i, err)
		}
		accounts = append(accounts, *acc)
	}
	return accounts, nil
}
