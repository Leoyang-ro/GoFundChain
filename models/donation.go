package models

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Donation 表示一笔捐赠记录
type Donation struct {
	// ID 捐赠记录的唯一标识符
	ID string `json:"id"`

	// Donor 捐赠者的以太坊钱包地址
	Donor common.Address `json:"donor"`

	// Amount 捐赠金额（以Wei为单位）
	Amount *big.Int `json:"amount"`

	// ProjectID 关联的项目ID
	ProjectID string `json:"project_id"`

	// Timestamp 捐赠时间戳
	// *big.Int 存储Unix时间戳
	Timestamp *big.Int `json:"timestamp"`

	// TxHash 区块链交易哈希
	// 用于在区块链上验证这笔捐赠
	TxHash string `json:"tx_hash"`

	// Status 捐赠状态
	// pending: 待确认, confirmed: 已确认, failed: 失败
	Status string `json:"status"`

	// Message 捐赠留言（可选）
	Message string `json:"message"`

	// CreatedAt 记录创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 记录更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// 捐赠状态常量
const (
	DonationStatusPending   = "pending"   // 待确认
	DonationStatusConfirmed = "confirmed" // 已确认
	DonationStatusFailed    = "failed"    // 失败
)

// NewDonation 创建一个新的捐赠记录
func NewDonation(donor common.Address, amount *big.Int, projectID string, txHash string) *Donation {
	return &Donation{
		Donor:      donor,
		Amount:     amount,
		ProjectID:  projectID,
		Timestamp:  big.NewInt(time.Now().Unix()),
		TxHash:     txHash,
		Status:     DonationStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// Confirm 确认捐赠
func (d *Donation) Confirm() {
	d.Status = DonationStatusConfirmed
	d.UpdatedAt = time.Now()
}

// Fail 标记捐赠失败
func (d *Donation) Fail() {
	d.Status = DonationStatusFailed
	d.UpdatedAt = time.Now()
}

// IsConfirmed 检查捐赠是否已确认
func (d *Donation) IsConfirmed() bool {
	return d.Status == DonationStatusConfirmed
}

// IsPending 检查捐赠是否待确认
func (d *Donation) IsPending() bool {
	return d.Status == DonationStatusPending
}

// IsFailed 检查捐赠是否失败
func (d *Donation) IsFailed() bool {
	return d.Status == DonationStatusFailed
}

// GetDonationTime 获取捐赠时间
func (d *Donation) GetDonationTime() time.Time {
	return time.Unix(d.Timestamp.Int64(), 0)
}

// SetMessage 设置捐赠留言
func (d *Donation) SetMessage(message string) {
	d.Message = message
	d.UpdatedAt = time.Now()
}
