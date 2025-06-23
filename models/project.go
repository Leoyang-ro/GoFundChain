package models

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Project 表示一个众筹项目
// 这是Go语言中的结构体定义，类似于其他语言中的类
type Project struct {
	// ID 项目的唯一标识符
	// string 是Go语言中的字符串类型
	ID string `json:"id"`

	// Title 项目标题
	Title string `json:"title"`

	// Description 项目详细描述
	Description string `json:"description"`

	// Creator 项目创建者的以太坊钱包地址
	// common.Address 是go-ethereum库中定义的钱包地址类型
	Creator common.Address `json:"creator"`

	// Goal 众筹目标金额（以Wei为单位）
	// *big.Int 是大整数类型，用于处理大金额，避免精度问题
	Goal *big.Int `json:"goal"`

	// Raised 已筹集金额（以Wei为单位）
	Raised *big.Int `json:"raised"`

	// Deadline 众筹截止时间
	// *big.Int 存储Unix时间戳
	Deadline *big.Int `json:"deadline"`

	// IPFSHash 项目详细信息在IPFS上的哈希值
	// IPFS是去中心化文件存储系统
	IPFSHash string `json:"ipfs_hash"`

	// IsActive 项目是否处于活跃状态
	IsActive bool `json:"is_active"`

	// CreatedAt 项目创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 项目最后更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// NewProject 创建一个新的项目实例
// 这是Go语言中的构造函数模式
func NewProject(title, description string, creator common.Address, goal *big.Int, deadline *big.Int) *Project {
	return &Project{
		Title:       title,
		Description: description,
		Creator:     creator,
		Goal:        goal,
		Raised:      big.NewInt(0), // 初始筹集金额为0
		Deadline:    deadline,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// IsExpired 检查项目是否已过期
// 这是结构体的方法，类似于其他语言中的实例方法
func (p *Project) IsExpired() bool {
	// 将big.Int转换为int64进行比较
	deadline := p.Deadline.Int64()
	currentTime := time.Now().Unix()
	return currentTime > deadline
}

// GetProgress 获取项目进度百分比
func (p *Project) GetProgress() float64 {
	if p.Goal.Cmp(big.NewInt(0)) == 0 {
		return 0.0
	}
	
	// 计算百分比：(已筹集金额 / 目标金额) * 100
	progress := new(big.Float).Quo(
		new(big.Float).SetInt(p.Raised),
		new(big.Float).SetInt(p.Goal),
	)
	
	progress.Mul(progress, big.NewFloat(100))
	result, _ := progress.Float64()
	return result
}

// AddDonation 添加捐赠金额
func (p *Project) AddDonation(amount *big.Int) {
	// 将新捐赠金额加到已筹集金额上
	p.Raised.Add(p.Raised, amount)
	p.UpdatedAt = time.Now()
}

// IsGoalReached 检查是否达到目标金额
func (p *Project) IsGoalReached() bool {
	// 比较已筹集金额是否大于等于目标金额
	return p.Raised.Cmp(p.Goal) >= 0
}
