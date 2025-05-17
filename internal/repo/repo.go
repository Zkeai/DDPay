package repo

import (
	"context"
	"time"

	"github.com/Zkeai/go_template/internal/conf"
	"github.com/Zkeai/go_template/internal/repo/db"
)

type Repo struct {
	db *db.DB
}

func NewRepo(conf *conf.Conf) *Repo {
	return &Repo{
		db: db.NewDB(conf.DB),
	}
}

func (r *Repo) UserRegister(ctx context.Context, wallet string, chainId int64) (*db.User, error) {

	return r.db.InsertUser(ctx, wallet, chainId)
}

func (r *Repo) UserQuery(ctx context.Context, wallet string) (*db.User, error) {

	return r.db.QueryUser(ctx, wallet)
}

func (r *Repo) UpdateUser(ctx context.Context, wallet string, ip string) (bool, error) {

	return r.db.UpdateUser(ctx, wallet, ip)
}

// AddMembership 添加会员记录
func (r *Repo) AddMembership(ctx context.Context, userID int64, orderID string, memType string, start, end time.Time) error {
	return r.db.AddMembership(ctx, userID, orderID, memType, start, end)
}

// GetAllMemberships 查询所有会员记录（可排序）
func (r *Repo) GetAllMemberships(ctx context.Context) ([]*db.Membership, error) {
	return r.db.GetAllMemberships(ctx)
}

// GetMembershipsByUserID 查询某个用户的所有会员记录
func (r *Repo) GetMembershipsByUserID(ctx context.Context, userID int64) ([]*db.Membership, error) {
	return r.db.GetMembershipsByUserID(ctx, userID)
}

// GetCurrentValidMembership 查询某用户当前有效会员（未过期）
func (r *Repo) GetCurrentValidMembership(ctx context.Context, userID int64) (*db.Membership, error) {
	return r.db.GetCurrentValidMembership(ctx, userID)
}

// GetMembershipByOrderID 根据orderID查询会员记录
func (r *Repo) GetMembershipByOrderID(ctx context.Context, orderID string) (*db.Membership, error) {
	return r.db.GetMembershipByOrderID(ctx, orderID)
}

// GetLatestEndTimeByUserID 获取用户最晚的有效时间
func (r *Repo) GetLatestEndTimeByUserID(ctx context.Context, userID int64) (time.Time, error) {
	return r.db.GetLatestEndTimeByUserID(ctx, userID)
}

// GetUnexpiredMembershipsByUserID 获取所有未过期的会员记录
func (r *Repo) GetUnexpiredMembershipsByUserID(ctx context.Context, userID int64) ([]*db.Membership, error) {
	return r.db.GetUnexpiredMembershipsByUserID(ctx, userID)
}

// InsertOrder 插入订单
func (r *Repo) InsertOrder(ctx context.Context, order *db.Order) error {
	return r.db.InsertOrder(ctx, order)
}

// QueryOrders 分页查询订单（支持链名、状态、时间范围）
func (r *Repo) QueryOrders(ctx context.Context, params db.OrderQueryParams) ([]*db.Order, error) {
	return r.db.QueryOrders(ctx, params)
}

// CountOrders 获取订单总数（用于分页）
func (r *Repo) CountOrders(ctx context.Context, params db.OrderQueryParams) (int64, error) {
	return r.db.CountOrders(ctx, params)
}

// GetOrderByOrderID 根据 OrderID 获取订单
func (r *Repo) GetOrderByOrderID(ctx context.Context, orderID string) (*db.Order, error) {
	return r.db.GetOrderByOrderID(ctx, orderID)
}

// GetOrderByTxHash 根据 TxHash 获取订单
func (r *Repo) GetOrderByTxHash(ctx context.Context, txHash string) (*db.Order, error) {
	return r.db.GetOrderByTxHash(ctx, txHash)
}

// UpdateOrderStatus 更新订单状态
func (r *Repo) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	return r.db.UpdateOrderStatus(ctx, orderID, status)
}

// Exists 判断指定用户的指定通道是否已存在。
// 如果存在返回 true，不存在返回 false。
func (r *Repo) Exists(ctx context.Context, userID int, channelID string) (bool, error) {
	return r.db.Exists(ctx, userID, channelID)
}

// Create 创建一个新的用户通道记录。
// 必填字段为 userID 和 channelID，其他字段可选，根据前端传参填写。
func (r *Repo) Create(ctx context.Context, c *db.UserChannel) error {
	return r.db.Create(ctx, c)
}

// Update 更新指定用户的指定通道记录。
// 只更新 UserChannel 结构体中提供的字段。
func (r *Repo) Update(ctx context.Context, userID int, channelID string, c *db.UserChannel) error {
	return r.db.Update(ctx, userID, channelID, c)
}

// GetByUserID 根据 userID 查询该用户的所有通道记录。
// 已被标记删除（is_delete = 1）的记录应在 db 层中被过滤。
func (r *Repo) GetByUserID(ctx context.Context, userID int) ([]*db.UserChannel, error) {
	return r.db.GetByUserID(ctx, userID)
}

// SetDisabled 设置指定用户通道的启用/禁用状态。
// disabled 为 0 表示启用，1 表示禁用。
func (r *Repo) SetDisabled(ctx context.Context, userID int, channelID string, disabled int, reson string) error {
	return r.db.SetDisabled(ctx, userID, channelID, disabled, reson)
}

// SoftDelete 逻辑删除指定用户的指定通道。
// 实际为将 is_delete 字段设置为 1，不会真正从数据库中移除记录。
func (r *Repo) SoftDelete(ctx context.Context, userID int, channelID string) error {
	return r.db.SoftDelete(ctx, userID, channelID)
}
