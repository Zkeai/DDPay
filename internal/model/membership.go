package model

import "time"

// MembershipLevel 会员等级模型
type MembershipLevel struct {
	ID               int64     `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`                         // 等级名称
	Level            int       `json:"level" db:"level"`                       // 等级数值
	Icon             string    `json:"icon" db:"icon"`                         // 等级图标URL
	Price            float64   `json:"price" db:"price"`                       // 升级价格
	Description      string    `json:"description" db:"description"`           // 等级描述
	DiscountRate     float64   `json:"discount_rate" db:"discount_rate"`       // 折扣率 (0.1-1.0)
	MaxSubsites      int       `json:"max_subsites" db:"max_subsites"`         // 最大分站数量
	CustomServiceAccess bool   `json:"custom_service_access" db:"custom_service_access"` // 专属客服权限
	VIPGroupAccess   bool      `json:"vip_group_access" db:"vip_group_access"` // VIP群权限
	Priority         int       `json:"priority" db:"priority"`                 // 优先级(处理订单等)
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	Benefits         []*MembershipBenefit `json:"benefits,omitempty" db:"-"`   // 等级特权列表
	Requirements     []*MembershipRequirement `json:"requirements,omitempty" db:"-"` // 升级条件列表
}

// MembershipBenefit 会员权益模型
type MembershipBenefit struct {
	ID               int64     `json:"id" db:"id"`
	LevelID          int64     `json:"level_id" db:"level_id"`       // 关联的等级ID
	Title            string    `json:"title" db:"title"`             // 权益标题
	Description      string    `json:"description" db:"description"` // 权益描述
	Icon             string    `json:"icon" db:"icon"`               // 权益图标
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// MembershipRequirement 会员升级条件模型
type MembershipRequirement struct {
	ID               int64     `json:"id" db:"id"`
	LevelID          int64     `json:"level_id" db:"level_id"`       // 关联的等级ID
	Type             string    `json:"type" db:"type"`               // 条件类型(充值金额/订单数/交易额/邀请人数)
	Value            float64   `json:"value" db:"value"`             // 条件值
	Description      string    `json:"description" db:"description"` // 条件描述
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UserMembership 用户会员记录模型
type UserMembership struct {
	ID               int64      `json:"id" db:"id"`
	UserID           int64      `json:"user_id" db:"user_id"`             // 用户ID
	LevelID          int64      `json:"level_id" db:"level_id"`           // 会员等级ID
	StartDate        time.Time  `json:"start_date" db:"start_date"`       // 开始日期
	EndDate          *time.Time `json:"end_date" db:"end_date"`           // 结束日期
	IsActive         bool       `json:"is_active" db:"is_active"`         // 是否激活
	PurchaseAmount   float64    `json:"purchase_amount" db:"purchase_amount"` // 购买金额
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// MembershipTransaction 会员交易记录模型
type MembershipTransaction struct {
	ID               int64     `json:"id" db:"id"`
	UserID           int64     `json:"user_id" db:"user_id"`               // 用户ID
	LevelID          int64     `json:"level_id" db:"level_id"`             // 会员等级ID
	Amount           float64   `json:"amount" db:"amount"`                 // 交易金额
	TransactionType  string    `json:"transaction_type" db:"transaction_type"` // 交易类型(购买/续费/升级)
	PaymentMethod    string    `json:"payment_method" db:"payment_method"` // 支付方式
	Status           string    `json:"status" db:"status"`                 // 交易状态
	OrderID          string    `json:"order_id" db:"order_id"`             // 订单号
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// GetDefaultMembershipLevels 获取默认会员等级配置
func GetDefaultMembershipLevels() []*MembershipLevel {
	now := time.Now()
	
	return []*MembershipLevel{
		{
			ID:               1,
			Name:             "青铜会员",
			Level:            1,
			Icon:             "/assets/membership/bronze.png",
			Price:            0,
			Description:      "基础会员级别，享受平台基础服务",
			DiscountRate:     1.0,
			MaxSubsites:      1,
			CustomServiceAccess: false,
			VIPGroupAccess:   false,
			Priority:         1,
			CreatedAt:        now,
			UpdatedAt:        now,
			Benefits: []*MembershipBenefit{
				{
					Title:       "创建1个分站",
					Description: "可创建并管理1个分站",
					Icon:        "/assets/membership/site.png",
				},
				{
					Title:       "标准技术支持",
					Description: "工作日9:00-18:00客服支持",
					Icon:        "/assets/membership/support.png",
				},
			},
			Requirements: []*MembershipRequirement{
				{
					Type:        "register",
					Value:       0,
					Description: "注册成为会员即可",
				},
			},
		},
		{
			ID:               2,
			Name:             "白银会员",
			Level:            2,
			Icon:             "/assets/membership/silver.png",
			Price:            99,
			Description:      "进阶会员级别，享受更多权益与折扣",
			DiscountRate:     0.9,
			MaxSubsites:      3,
			CustomServiceAccess: false,
			VIPGroupAccess:   true,
			Priority:         2,
			CreatedAt:        now,
			UpdatedAt:        now,
			Benefits: []*MembershipBenefit{
				{
					Title:       "创建3个分站",
					Description: "可创建并管理最多3个分站",
					Icon:        "/assets/membership/sites.png",
				},
				{
					Title:       "9折优惠",
					Description: "所有交易享受9折优惠",
					Icon:        "/assets/membership/discount.png",
				},
				{
					Title:       "会员专属群",
					Description: "加入白银会员专属交流群",
					Icon:        "/assets/membership/group.png",
				},
			},
			Requirements: []*MembershipRequirement{
				{
					Type:        "payment",
					Value:       99,
					Description: "一次性支付99元升级",
				},
				{
					Type:        "total_order",
					Value:       5,
					Description: "累计订单达到5笔",
				},
			},
		},
		{
			ID:               3,
			Name:             "黄金会员",
			Level:            3,
			Icon:             "/assets/membership/gold.png",
			Price:            299,
			Description:      "高级会员，享受VIP待遇与最大折扣",
			DiscountRate:     0.8,
			MaxSubsites:      10,
			CustomServiceAccess: true,
			VIPGroupAccess:   true,
			Priority:         3,
			CreatedAt:        now,
			UpdatedAt:        now,
			Benefits: []*MembershipBenefit{
				{
					Title:       "创建10个分站",
					Description: "可创建并管理最多10个分站",
					Icon:        "/assets/membership/sites.png",
				},
				{
					Title:       "8折优惠",
					Description: "所有交易享受8折优惠",
					Icon:        "/assets/membership/discount.png",
				},
				{
					Title:       "专属客服",
					Description: "一对一专属客服7*12小时服务",
					Icon:        "/assets/membership/vip-support.png",
				},
				{
					Title:       "VIP会员群",
					Description: "加入金牌会员VIP交流群",
					Icon:        "/assets/membership/vip-group.png",
				},
				{
					Title:       "优先处理",
					Description: "订单优先处理，技术问题优先解决",
					Icon:        "/assets/membership/priority.png",
				},
			},
			Requirements: []*MembershipRequirement{
				{
					Type:        "payment",
					Value:       299,
					Description: "一次性支付299元升级",
				},
				{
					Type:        "total_payment",
					Value:       1000,
					Description: "累计充值金额达到1000元",
				},
				{
					Type:        "invitation",
					Value:       3,
					Description: "成功邀请3名用户注册",
				},
			},
		},
		{
			ID:               4,
			Name:             "钻石会员",
			Level:            4,
			Icon:             "/assets/membership/diamond.png",
			Price:            999,
			Description:      "顶级会员待遇，无限制使用所有功能",
			DiscountRate:     0.7,
			MaxSubsites:      -1, // 无限制
			CustomServiceAccess: true,
			VIPGroupAccess:   true,
			Priority:         4,
			CreatedAt:        now,
			UpdatedAt:        now,
			Benefits: []*MembershipBenefit{
				{
					Title:       "无限分站",
					Description: "可创建无限数量的分站",
					Icon:        "/assets/membership/unlimited.png",
				},
				{
					Title:       "7折优惠",
					Description: "所有交易享受7折最大优惠",
					Icon:        "/assets/membership/discount.png",
				},
				{
					Title:       "24小时专属客服",
					Description: "一对一专属客服全天候服务",
					Icon:        "/assets/membership/vip-support.png",
				},
				{
					Title:       "钻石VIP群",
					Description: "加入钻石会员核心交流群",
					Icon:        "/assets/membership/vip-group.png",
				},
				{
					Title:       "最高优先级",
					Description: "享受系统最高优先级处理所有请求",
					Icon:        "/assets/membership/priority.png",
				},
				{
					Title:       "专属定制",
					Description: "享受专属定制开发服务",
					Icon:        "/assets/membership/custom.png",
				},
			},
			Requirements: []*MembershipRequirement{
				{
					Type:        "payment",
					Value:       999,
					Description: "一次性支付999元升级",
				},
				{
					Type:        "total_payment",
					Value:       5000,
					Description: "累计充值金额达到5000元",
				},
				{
					Type:        "total_transaction",
					Value:       10000,
					Description: "累计交易额达到10000元",
				},
				{
					Type:        "invitation",
					Value:       10,
					Description: "成功邀请10名用户注册",
				},
			},
		},
	}
} 