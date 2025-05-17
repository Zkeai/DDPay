package geysergrpc

// SubscribeConfig 用于用户订阅配置
type SubscribeConfig struct {
	Endpoint   string
	Token      string
	Insecure   bool
	Signature  *string
	Resub      uint
	Slots      bool
	Blocks     bool
	BlocksMeta bool
	Accounts   bool

	AccountsFilter              []string
	AccountOwnersFilter         []string
	Transactions                bool
	VoteTransactions            *bool
	FailedTransactions          *bool
	TransactionsAccountsInclude []string
	TransactionsAccountsExclude []string
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "string representation of flag"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// GeyserService 定义订阅服务接口，便于以后扩展
type GeyserService interface {
	RunSubscription(cfg SubscribeConfig) error
}
