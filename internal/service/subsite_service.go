package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
	"github.com/Zkeai/DDPay/internal/repo"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// SubsiteService 分站服务接口
type SubsiteService interface {
	// 分站基本操作
	CreateSubsite(ctx context.Context, ownerID int64, name, subdomain string, description string, theme string, commissionRate float64) (*model.Subsite, error)
	CreateSubsiteObject(ctx context.Context, subsite *model.Subsite) (*model.Subsite, error)
	GetSubsiteByID(ctx context.Context, id int64) (*model.Subsite, error)
	GetSubsiteByOwnerID(ctx context.Context, ownerID int64) (*model.Subsite, error)
	GetSubsiteByDomain(ctx context.Context, domain string) (*model.Subsite, error)
	GetSubsiteBySubdomain(ctx context.Context, subdomain string) (*model.Subsite, error)
	UpdateSubsite(ctx context.Context, subsite *model.Subsite) error
	DeleteSubsite(ctx context.Context, id int64) error
	ListSubsites(ctx context.Context, page, pageSize int, status int) ([]*model.Subsite, int, error)
	GetSubsiteInfo(ctx context.Context, subsiteID int64) (*model.SubsiteInfo, error)
	
	// 分站JSON配置操作
	SaveSubsiteJsonConfig(ctx context.Context, subsiteID int64, config map[string]interface{}) error
	GetSubsiteJsonConfig(ctx context.Context, subsiteID int64) (map[string]interface{}, error)
	
	// 分站商品操作
	CreateSubsiteProduct(ctx context.Context, product *model.SubsiteProduct) (int64, error)
	GetSubsiteProduct(ctx context.Context, id int64) (*model.SubsiteProduct, error)
	UpdateSubsiteProduct(ctx context.Context, product *model.SubsiteProduct) error
	ListSubsiteProducts(ctx context.Context, subsiteID int64, page, pageSize int, status int) ([]*model.SubsiteProduct, int, error)
	
	// 分站订单操作
	CreateSubsiteOrder(ctx context.Context, subsiteID, userID, productID int64, quantity int) (*model.SubsiteOrder, error)
	GetSubsiteOrder(ctx context.Context, id int64) (*model.SubsiteOrder, error)
	GetSubsiteOrderByOrderNo(ctx context.Context, orderNo string) (*model.SubsiteOrder, error)
	UpdateSubsiteOrderStatus(ctx context.Context, orderID int64, status int) error
	ListSubsiteOrders(ctx context.Context, subsiteID int64, page, pageSize int, status int) ([]*model.SubsiteOrder, int, error)
	
	// 分站余额操作
	GetSubsiteBalance(ctx context.Context, ownerID int64) (*model.SubsiteBalance, error)
	AddSubsiteBalance(ctx context.Context, ownerID, orderID int64, amount float64, remark string) error
	SubtractSubsiteBalance(ctx context.Context, ownerID int64, amount float64, withdrawalID int64, remark string) error
	ListSubsiteBalanceLogs(ctx context.Context, ownerID int64, page, pageSize int) ([]*model.SubsiteBalanceLog, int, error)
	
	// 分站提现操作
	CreateSubsiteWithdrawal(ctx context.Context, withdrawal *model.SubsiteWithdrawal) (int64, error)
	GetSubsiteWithdrawal(ctx context.Context, id int64) (*model.SubsiteWithdrawal, error)
	ProcessSubsiteWithdrawal(ctx context.Context, id int64, status int, adminRemark string) error
	ListSubsiteWithdrawals(ctx context.Context, ownerID int64, page, pageSize int, status int) ([]*model.SubsiteWithdrawal, int, error)
}

// subsiteService 分站服务实现
type subsiteService struct {
	subsiteRepo  repo.SubsiteRepo
	userRepo     repo.UserRepo
}

// NewSubsiteService 创建分站服务
func NewSubsiteService(subsiteRepo repo.SubsiteRepo, userRepo repo.UserRepo) SubsiteService {
	return &subsiteService{
		subsiteRepo: subsiteRepo,
		userRepo:    userRepo,
	}
}

// CreateSubsite 创建分站
func (s *subsiteService) CreateSubsite(ctx context.Context, ownerID int64, name, subdomain string, description string, theme string, commissionRate float64) (*model.Subsite, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetUserByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	
	// 检查用户等级和已有分站数量
	if err := s.checkUserSubsiteLimit(ctx, user); err != nil {
		return nil, err
	}
	
	// 检查子域名是否已存在
	existingSubdomain, err := s.subsiteRepo.GetSubsiteBySubdomain(ctx, subdomain)
	if err != nil {
		return nil, err
	}
	if existingSubdomain != nil {
		return nil, errors.New("子域名已被占用")
	}
	
	// 创建分站
	subsite := &model.Subsite{
		OwnerID:        ownerID,
		Name:           name,
		Subdomain:      subdomain,
		Description:    description,
		Theme:          theme,
		Status:         model.SubsiteStatusEnabled,
		CommissionRate: commissionRate,
	}
	
	// 如果主题为空，设置默认主题
	if subsite.Theme == "" {
		subsite.Theme = "default"
	}
	
	id, err := s.subsiteRepo.CreateSubsite(ctx, subsite)
	if err != nil {
		return nil, err
	}
	
	subsite.ID = id
	return subsite, nil
}

// CreateSubsiteObject 使用完整的Subsite对象创建分站
func (s *subsiteService) CreateSubsiteObject(ctx context.Context, subsite *model.Subsite) (*model.Subsite, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetUserByID(ctx, subsite.OwnerID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	
	// 检查用户等级和已有分站数量
	if err := s.checkUserSubsiteLimit(ctx, user); err != nil {
		return nil, err
	}
	
	// 检查子域名是否已存在
	existingSubdomain, err := s.subsiteRepo.GetSubsiteBySubdomain(ctx, subsite.Subdomain)
	if err != nil {
		return nil, err
	}
	if existingSubdomain != nil {
		return nil, errors.New("子域名已被占用")
	}
	
	// 如果主题为空，设置默认主题
	if subsite.Theme == "" {
		subsite.Theme = "default"
	}
	
	// 如果状态为0，设置为默认启用状态
	if subsite.Status == 0 {
		subsite.Status = model.SubsiteStatusEnabled
	}
	
	id, err := s.subsiteRepo.CreateSubsite(ctx, subsite)
	if err != nil {
		return nil, err
	}
	
	subsite.ID = id
	return subsite, nil
}

// checkUserSubsiteLimit 检查用户是否超过分站数量限制
func (s *subsiteService) checkUserSubsiteLimit(ctx context.Context, user *model.User) error {
	// 获取用户已有的分站数量
	subsites, _, err := s.subsiteRepo.ListSubsitesByOwnerID(ctx, user.ID)
	if err != nil {
		return err
	}
	
	// 根据用户等级判断分站数量限制
	switch user.Level {
	case 1: // 青铜会员最多1个分站
		if len(subsites) >= 1 {
			return errors.New("青铜会员最多只能创建1个分站")
		}
	case 2: // 白银会员最多3个分站
		if len(subsites) >= 3 {
			return errors.New("白银会员最多只能创建3个分站")
		}
	case 3: // 黄金会员最多10个分站
		if len(subsites) >= 10 {
			return errors.New("黄金会员最多只能创建10个分站")
		}
	case 4: // 钻石会员无限制
		// 无限制
	default: // 默认按青铜会员处理
		if len(subsites) >= 1 {
			return errors.New("您的账户最多只能创建1个分站")
		}
	}
	
	return nil
}

// GetSubsiteByID 根据ID获取分站
func (s *subsiteService) GetSubsiteByID(ctx context.Context, id int64) (*model.Subsite, error) {
	return s.subsiteRepo.GetSubsiteByID(ctx, id)
}

// GetSubsiteByOwnerID 根据所有者ID获取分站
func (s *subsiteService) GetSubsiteByOwnerID(ctx context.Context, ownerID int64) (*model.Subsite, error) {
	return s.subsiteRepo.GetSubsiteByOwnerID(ctx, ownerID)
}

// GetSubsiteByDomain 根据域名获取分站
func (s *subsiteService) GetSubsiteByDomain(ctx context.Context, domain string) (*model.Subsite, error) {
	return s.subsiteRepo.GetSubsiteByDomain(ctx, domain)
}

// GetSubsiteBySubdomain 根据子域名获取分站
func (s *subsiteService) GetSubsiteBySubdomain(ctx context.Context, subdomain string) (*model.Subsite, error) {
	return s.subsiteRepo.GetSubsiteBySubdomain(ctx, subdomain)
}

// UpdateSubsite 更新分站
func (s *subsiteService) UpdateSubsite(ctx context.Context, subsite *model.Subsite) error {
	// 检查子域名是否已被占用
	if subsite.Subdomain != "" {
		existingSubsite, err := s.subsiteRepo.GetSubsiteBySubdomain(ctx, subsite.Subdomain)
		if err != nil {
			return err
		}
		if existingSubsite != nil && existingSubsite.ID != subsite.ID {
			return errors.New("子域名已被占用")
		}
	}
	
	// 检查域名是否已被占用
	if subsite.Domain != "" {
		existingSubsite, err := s.subsiteRepo.GetSubsiteByDomain(ctx, subsite.Domain)
		if err != nil {
			return err
		}
		if existingSubsite != nil && existingSubsite.ID != subsite.ID {
			return errors.New("域名已被占用")
		}
	}
	
	return s.subsiteRepo.UpdateSubsite(ctx, subsite)
}

// ListSubsites 获取分站列表
func (s *subsiteService) ListSubsites(ctx context.Context, page, pageSize int, status int) ([]*model.Subsite, int, error) {
	// 获取分站基本信息
	subsites, total, err := s.subsiteRepo.ListSubsites(ctx, page, pageSize, status)
	if err != nil {
		return nil, 0, err
	}
	
	// 为每个分站获取所有者信息
	for _, subsite := range subsites {
		owner, err := s.userRepo.GetUserByID(ctx, subsite.OwnerID)
		if err == nil && owner != nil {
			// 确保用户存在
		}
	}
	
	return subsites, total, nil
}

// GetSubsiteInfo 获取分站详细信息
func (s *subsiteService) GetSubsiteInfo(ctx context.Context, subsiteID int64) (*model.SubsiteInfo, error) {
	// 获取分站信息
	subsite, err := s.subsiteRepo.GetSubsiteByID(ctx, subsiteID)
	if err != nil {
		return nil, err
	}
	if subsite == nil {
		return nil, errors.New("分站不存在")
	}
	
	// 获取所有者信息
	owner, err := s.userRepo.GetUserByID(ctx, subsite.OwnerID)
	if err != nil {
		return nil, err
	}
	if owner == nil {
		return nil, errors.New("分站所有者不存在")
	}
	
	// 获取商品数量
	products, _, err := s.subsiteRepo.ListSubsiteProducts(ctx, subsiteID, 1, 1, -1)
	if err != nil {
		return nil, err
	}
	
	// 获取订单数量
	orders, _, err := s.subsiteRepo.ListSubsiteOrders(ctx, subsiteID, 1, 1, -1)
	if err != nil {
		return nil, err
	}
	
	// 获取余额
	balance, err := s.subsiteRepo.GetSubsiteBalance(ctx, subsite.OwnerID)
	if err != nil {
		return nil, err
	}
	
	return &model.SubsiteInfo{
		Subsite:      subsite,
		Owner:        owner.ToProfile(),
		ProductCount: len(products),
		OrderCount:   len(orders),
		Balance:      balance.Amount,
	}, nil
}

// SaveSubsiteJsonConfig 保存分站JSON配置
func (s *subsiteService) SaveSubsiteJsonConfig(ctx context.Context, subsiteID int64, config map[string]interface{}) error {
	// 将map转换为JSON字符串
	jsonData, err := json.Marshal(config)
	if err != nil {
		return err
	}
	
	// 创建或更新配置
	subsiteConfig := &model.SubsiteConfig{
		SubsiteID: subsiteID,
		Config:    string(jsonData),
	}
	
	return s.subsiteRepo.SaveSubsiteConfig(ctx, subsiteConfig)
}

// GetSubsiteJsonConfig 获取分站JSON配置
func (s *subsiteService) GetSubsiteJsonConfig(ctx context.Context, subsiteID int64) (map[string]interface{}, error) {
	// 获取配置
	config, err := s.subsiteRepo.GetSubsiteConfig(ctx, subsiteID)
	if err != nil {
		return nil, err
	}
	
	// 如果配置不存在，返回空map
	if config == nil {
		return map[string]interface{}{}, nil
	}
	
	// 解析JSON
	var configMap map[string]interface{}
	err = json.Unmarshal([]byte(config.Config), &configMap)
	if err != nil {
		return nil, err
	}
	
	return configMap, nil
}

// CreateSubsiteProduct 创建分站商品
func (s *subsiteService) CreateSubsiteProduct(ctx context.Context, product *model.SubsiteProduct) (int64, error) {
	// 检查分站是否存在
	subsite, err := s.subsiteRepo.GetSubsiteByID(ctx, product.SubsiteID)
	if err != nil {
		return 0, err
	}
	if subsite == nil {
		return 0, errors.New("分站不存在")
	}
	
	// 如果是限时商品，检查时间
	if product.IsTimeLimited == 1 {
		if product.StartTime.IsZero() || product.EndTime.IsZero() {
			return 0, errors.New("限时商品必须设置开始和结束时间")
		}
		if product.EndTime.Before(product.StartTime) {
			return 0, errors.New("结束时间不能早于开始时间")
		}
	}
	
	return s.subsiteRepo.CreateSubsiteProduct(ctx, product)
}

// GetSubsiteProduct 获取分站商品
func (s *subsiteService) GetSubsiteProduct(ctx context.Context, id int64) (*model.SubsiteProduct, error) {
	return s.subsiteRepo.GetSubsiteProductByID(ctx, id)
}

// UpdateSubsiteProduct 更新分站商品
func (s *subsiteService) UpdateSubsiteProduct(ctx context.Context, product *model.SubsiteProduct) error {
	// 检查商品是否存在
	existingProduct, err := s.subsiteRepo.GetSubsiteProductByID(ctx, product.ID)
	if err != nil {
		return err
	}
	if existingProduct == nil {
		return errors.New("商品不存在")
	}
	
	// 如果是限时商品，检查时间
	if product.IsTimeLimited == 1 {
		if product.StartTime.IsZero() || product.EndTime.IsZero() {
			return errors.New("限时商品必须设置开始和结束时间")
		}
		if product.EndTime.Before(product.StartTime) {
			return errors.New("结束时间不能早于开始时间")
		}
	}
	
	return s.subsiteRepo.UpdateSubsiteProduct(ctx, product)
}

// ListSubsiteProducts 获取分站商品列表
func (s *subsiteService) ListSubsiteProducts(ctx context.Context, subsiteID int64, page, pageSize int, status int) ([]*model.SubsiteProduct, int, error) {
	return s.subsiteRepo.ListSubsiteProducts(ctx, subsiteID, page, pageSize, status)
}

// CreateSubsiteOrder 创建分站订单
func (s *subsiteService) CreateSubsiteOrder(ctx context.Context, subsiteID, userID, productID int64, quantity int) (*model.SubsiteOrder, error) {
	// 检查分站是否存在
	subsite, err := s.subsiteRepo.GetSubsiteByID(ctx, subsiteID)
	if err != nil {
		return nil, err
	}
	if subsite == nil {
		return nil, errors.New("分站不存在")
	}
	
	// 检查商品是否存在
	product, err := s.subsiteRepo.GetSubsiteProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("商品不存在")
	}
	
	// 检查商品是否属于该分站
	if product.SubsiteID != subsiteID {
		return nil, errors.New("商品不属于该分站")
	}
	
	// 检查商品状态
	if product.Status != model.SubsiteProductStatusOnline {
		return nil, errors.New("商品已下架")
	}
	
	// 检查库存
	if product.Stock > 0 && quantity > product.Stock {
		return nil, errors.New("库存不足")
	}
	
	// 检查限时商品
	now := time.Now()
	if product.IsTimeLimited == 1 {
		if now.Before(product.StartTime) {
			return nil, errors.New("商品尚未开始销售")
		}
		if now.After(product.EndTime) {
			return nil, errors.New("商品已停止销售")
		}
	}
	
	// 生成订单号
	orderNo, err := gonanoid.New(16)
	if err != nil {
		return nil, err
	}
	
	// 计算订单金额和佣金
	amount := product.Price * float64(quantity)
	commission := amount * (subsite.CommissionRate / 100)
	
	// 创建订单
	order := &model.SubsiteOrder{
		OrderNo:    orderNo,
		SubsiteID:  subsiteID,
		UserID:     userID,
		ProductID:  productID,
		Quantity:   quantity,
		Amount:     amount,
		Commission: commission,
		Status:     model.SubsiteOrderStatusPending,
	}
	
	id, err := s.subsiteRepo.CreateSubsiteOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	
	order.ID = id
	return order, nil
}

// GetSubsiteOrder 获取分站订单
func (s *subsiteService) GetSubsiteOrder(ctx context.Context, id int64) (*model.SubsiteOrder, error) {
	return s.subsiteRepo.GetSubsiteOrderByID(ctx, id)
}

// GetSubsiteOrderByOrderNo 根据订单号获取分站订单
func (s *subsiteService) GetSubsiteOrderByOrderNo(ctx context.Context, orderNo string) (*model.SubsiteOrder, error) {
	return s.subsiteRepo.GetSubsiteOrderByOrderNo(ctx, orderNo)
}

// UpdateSubsiteOrderStatus 更新分站订单状态
func (s *subsiteService) UpdateSubsiteOrderStatus(ctx context.Context, orderID int64, status int) error {
	// 获取订单信息
	order, err := s.subsiteRepo.GetSubsiteOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}
	
	// 检查状态变更是否合法
	if order.Status == model.SubsiteOrderStatusCancelled {
		return errors.New("已取消的订单不能更改状态")
	}
	
	// 更新订单状态
	order.Status = status
	now := time.Now()
	
	// 如果订单状态变为已支付，记录支付时间
	if status == model.SubsiteOrderStatusPaid {
		order.PayTime = now
		
		// 获取商品信息
		product, err := s.subsiteRepo.GetSubsiteProductByID(ctx, order.ProductID)
		if err != nil {
			return err
		}
		if product == nil {
			return errors.New("商品不存在")
		}
		
		// 更新库存
		if product.Stock > 0 {
			product.Stock -= order.Quantity
			if product.Stock < 0 {
				product.Stock = 0
			}
			err = s.subsiteRepo.UpdateSubsiteProduct(ctx, product)
			if err != nil {
				return err
			}
		}
		
		// 获取分站信息
		subsite, err := s.subsiteRepo.GetSubsiteByID(ctx, order.SubsiteID)
		if err != nil {
			return err
		}
		if subsite == nil {
			return errors.New("分站不存在")
		}
		
		// 增加分站余额
		err = s.AddSubsiteBalance(ctx, subsite.OwnerID, order.ID, order.Commission, fmt.Sprintf("订单佣金: %s", order.OrderNo))
		if err != nil {
			return err
		}
	}
	
	// 如果订单状态变为已完成，记录完成时间
	if status == model.SubsiteOrderStatusCompleted {
		order.CompleteTime = now
	}
	
	return s.subsiteRepo.UpdateSubsiteOrder(ctx, order)
}

// ListSubsiteOrders 获取分站订单列表
func (s *subsiteService) ListSubsiteOrders(ctx context.Context, subsiteID int64, page, pageSize int, status int) ([]*model.SubsiteOrder, int, error) {
	return s.subsiteRepo.ListSubsiteOrders(ctx, subsiteID, page, pageSize, status)
}

// GetSubsiteBalance 获取分站余额
func (s *subsiteService) GetSubsiteBalance(ctx context.Context, ownerID int64) (*model.SubsiteBalance, error) {
	return s.subsiteRepo.GetSubsiteBalance(ctx, ownerID)
}

// AddSubsiteBalance 增加分站余额
func (s *subsiteService) AddSubsiteBalance(ctx context.Context, ownerID, orderID int64, amount float64, remark string) error {
	if amount <= 0 {
		return errors.New("增加金额必须大于0")
	}
	
	// 获取当前余额
	balance, err := s.subsiteRepo.GetSubsiteBalance(ctx, ownerID)
	if err != nil {
		return err
	}
	
	// 记录变动前的余额
	beforeBalance := balance.Amount
	
	// 增加余额
	balance.Amount += amount
	err = s.subsiteRepo.UpdateSubsiteBalance(ctx, balance)
	if err != nil {
		return err
	}
	
	// 记录余额变动日志
	log := &model.SubsiteBalanceLog{
		OwnerID:       ownerID,
		OrderID:       orderID,
		Amount:        amount,
		BeforeBalance: beforeBalance,
		AfterBalance:  balance.Amount,
		Type:          model.SubsiteBalanceTypeCommission,
		Remark:        remark,
	}
	
	return s.subsiteRepo.CreateSubsiteBalanceLog(ctx, log)
}

// SubtractSubsiteBalance 减少分站余额
func (s *subsiteService) SubtractSubsiteBalance(ctx context.Context, ownerID int64, amount float64, withdrawalID int64, remark string) error {
	if amount <= 0 {
		return errors.New("减少金额必须大于0")
	}
	
	// 获取当前余额
	balance, err := s.subsiteRepo.GetSubsiteBalance(ctx, ownerID)
	if err != nil {
		return err
	}
	
	// 检查余额是否足够
	if balance.Amount < amount {
		return errors.New("余额不足")
	}
	
	// 记录变动前的余额
	beforeBalance := balance.Amount
	
	// 减少余额
	balance.Amount -= amount
	err = s.subsiteRepo.UpdateSubsiteBalance(ctx, balance)
	if err != nil {
		return err
	}
	
	// 记录余额变动日志
	log := &model.SubsiteBalanceLog{
		OwnerID:       ownerID,
		OrderID:       withdrawalID, // 这里用提现ID代替订单ID
		Amount:        -amount,      // 负数表示减少
		BeforeBalance: beforeBalance,
		AfterBalance:  balance.Amount,
		Type:          model.SubsiteBalanceTypeWithdrawal,
		Remark:        remark,
	}
	
	return s.subsiteRepo.CreateSubsiteBalanceLog(ctx, log)
}

// ListSubsiteBalanceLogs 获取分站余额变动记录
func (s *subsiteService) ListSubsiteBalanceLogs(ctx context.Context, ownerID int64, page, pageSize int) ([]*model.SubsiteBalanceLog, int, error) {
	return s.subsiteRepo.ListSubsiteBalanceLogs(ctx, ownerID, page, pageSize)
}

// CreateSubsiteWithdrawal 创建分站提现申请
func (s *subsiteService) CreateSubsiteWithdrawal(ctx context.Context, withdrawal *model.SubsiteWithdrawal) (int64, error) {
	// 检查提现金额
	if withdrawal.Amount <= 0 {
		return 0, errors.New("提现金额必须大于0")
	}
	
	// 获取用户余额
	balance, err := s.subsiteRepo.GetSubsiteBalance(ctx, withdrawal.OwnerID)
	if err != nil {
		return 0, err
	}
	
	// 检查余额是否足够
	if balance.Amount < withdrawal.Amount {
		return 0, errors.New("余额不足")
	}
	
	// 设置提现状态为待处理
	withdrawal.Status = model.SubsiteWithdrawalStatusPending
	
	// 创建提现申请
	id, err := s.subsiteRepo.CreateSubsiteWithdrawal(ctx, withdrawal)
	if err != nil {
		return 0, err
	}
	
	// 减少用户余额
	err = s.SubtractSubsiteBalance(ctx, withdrawal.OwnerID, withdrawal.Amount, id, "提现申请")
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

// GetSubsiteWithdrawal 获取分站提现申请
func (s *subsiteService) GetSubsiteWithdrawal(ctx context.Context, id int64) (*model.SubsiteWithdrawal, error) {
	return s.subsiteRepo.GetSubsiteWithdrawalByID(ctx, id)
}

// ProcessSubsiteWithdrawal 处理分站提现申请
func (s *subsiteService) ProcessSubsiteWithdrawal(ctx context.Context, id int64, status int, adminRemark string) error {
	// 获取提现申请
	withdrawal, err := s.subsiteRepo.GetSubsiteWithdrawalByID(ctx, id)
	if err != nil {
		return err
	}
	if withdrawal == nil {
		return errors.New("提现申请不存在")
	}
	
	// 检查提现状态
	if withdrawal.Status != model.SubsiteWithdrawalStatusPending {
		return errors.New("只能处理待处理的提现申请")
	}
	
	// 更新提现状态
	withdrawal.Status = status
	withdrawal.AdminRemark = adminRemark
	withdrawal.ProcessedAt = time.Now()
	
	// 如果拒绝提现，退回余额
	if status == model.SubsiteWithdrawalStatusRejected {
		err = s.AddSubsiteBalance(ctx, withdrawal.OwnerID, id, withdrawal.Amount, "提现被拒绝，退回余额")
		if err != nil {
			return err
		}
	}
	
	return s.subsiteRepo.UpdateSubsiteWithdrawal(ctx, withdrawal)
}

// ListSubsiteWithdrawals 获取分站提现申请列表
func (s *subsiteService) ListSubsiteWithdrawals(ctx context.Context, ownerID int64, page, pageSize int, status int) ([]*model.SubsiteWithdrawal, int, error) {
	return s.subsiteRepo.ListSubsiteWithdrawals(ctx, ownerID, page, pageSize, status)
}

// DeleteSubsite 删除分站
func (s *subsiteService) DeleteSubsite(ctx context.Context, id int64) error {
	return s.subsiteRepo.DeleteSubsite(ctx, id)
} 