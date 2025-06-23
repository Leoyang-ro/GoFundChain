package models

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// TestNewProject 测试创建新项目
func TestNewProject(t *testing.T) {
	// 准备测试数据
	title := "测试项目"
	description := "这是一个测试项目"
	creator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	goal := big.NewInt(1000000000000000000) // 1 ETH in Wei
	deadline := big.NewInt(time.Now().AddDate(0, 1, 0).Unix()) // 一个月后

	// 创建新项目
	project := NewProject(title, description, creator, goal, deadline)

	// 验证项目属性
	if project.Title != title {
		t.Errorf("期望标题为 %s，实际为 %s", title, project.Title)
	}

	if project.Description != description {
		t.Errorf("期望描述为 %s，实际为 %s", description, project.Description)
	}

	if project.Creator != creator {
		t.Errorf("期望创建者为 %s，实际为 %s", creator.Hex(), project.Creator.Hex())
	}

	if project.Goal.Cmp(goal) != 0 {
		t.Errorf("期望目标金额为 %s，实际为 %s", goal.String(), project.Goal.String())
	}

	if project.Raised.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("期望初始筹集金额为 0，实际为 %s", project.Raised.String())
	}

	if !project.IsActive {
		t.Error("新项目应该是活跃状态")
	}
}

// TestProjectMethods 测试项目的方法
func TestProjectMethods(t *testing.T) {
	// 创建测试项目
	creator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	goal := big.NewInt(1000000000000000000) // 1 ETH
	deadline := big.NewInt(time.Now().AddDate(0, 0, 1).Unix()) // 明天
	project := NewProject("测试项目", "描述", creator, goal, deadline)

	// 测试初始进度
	progress := project.GetProgress()
	if progress != 0.0 {
		t.Errorf("期望初始进度为 0.0%%，实际为 %.2f%%", progress)
	}

	// 测试添加捐赠
	donationAmount := big.NewInt(500000000000000000) // 0.5 ETH
	project.AddDonation(donationAmount)

	// 验证筹集金额
	if project.Raised.Cmp(donationAmount) != 0 {
		t.Errorf("期望筹集金额为 %s，实际为 %s", donationAmount.String(), project.Raised.String())
	}

	// 验证进度
	progress = project.GetProgress()
	expectedProgress := 50.0
	if progress != expectedProgress {
		t.Errorf("期望进度为 %.2f%%，实际为 %.2f%%", expectedProgress, progress)
	}

	// 测试未达到目标
	if project.IsGoalReached() {
		t.Error("项目不应该达到目标")
	}

	// 添加更多捐赠达到目标
	project.AddDonation(donationAmount)

	// 测试达到目标
	if !project.IsGoalReached() {
		t.Error("项目应该达到目标")
	}
}

// TestNewDonation 测试创建新捐赠
func TestNewDonation(t *testing.T) {
	// 准备测试数据
	donor := common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12")
	amount := big.NewInt(100000000000000000) // 0.1 ETH
	projectID := "project-123"
	txHash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

	// 创建新捐赠
	donation := NewDonation(donor, amount, projectID, txHash)

	// 验证捐赠属性
	if donation.Donor != donor {
		t.Errorf("期望捐赠者为 %s，实际为 %s", donor.Hex(), donation.Donor.Hex())
	}

	if donation.Amount.Cmp(amount) != 0 {
		t.Errorf("期望捐赠金额为 %s，实际为 %s", amount.String(), donation.Amount.String())
	}

	if donation.ProjectID != projectID {
		t.Errorf("期望项目ID为 %s，实际为 %s", projectID, donation.ProjectID)
	}

	if donation.TxHash != txHash {
		t.Errorf("期望交易哈希为 %s，实际为 %s", txHash, donation.TxHash)
	}

	if donation.Status != DonationStatusPending {
		t.Errorf("期望状态为 %s，实际为 %s", DonationStatusPending, donation.Status)
	}
}

// TestDonationMethods 测试捐赠的方法
func TestDonationMethods(t *testing.T) {
	// 创建测试捐赠
	donor := common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12")
	amount := big.NewInt(100000000000000000)
	donation := NewDonation(donor, amount, "project-123", "tx-hash")

	// 测试初始状态
	if !donation.IsPending() {
		t.Error("新捐赠应该是待确认状态")
	}

	if donation.IsConfirmed() {
		t.Error("新捐赠不应该是已确认状态")
	}

	// 测试确认捐赠
	donation.Confirm()

	if !donation.IsConfirmed() {
		t.Error("捐赠确认后应该是已确认状态")
	}

	if donation.IsPending() {
		t.Error("已确认的捐赠不应该是待确认状态")
	}

	// 测试设置留言
	message := "支持这个项目！"
	donation.SetMessage(message)

	if donation.Message != message {
		t.Errorf("期望留言为 %s，实际为 %s", message, donation.Message)
	}
} 