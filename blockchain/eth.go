package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GoFundChainContract 智能合约交互结构体
type GoFundChainContract struct {
	client   *ethclient.Client
	contract *bind.BoundContract
	address  common.Address
	abi      abi.ABI
}

// Project 项目结构体
type Project struct {
	ID          string
	Title       string
	Description string
	Creator     common.Address
	Goal        *big.Int
	Raised      *big.Int
	Deadline    *big.Int
	IPFSHash    string
	IsActive    bool
}

// Donation 捐赠结构体
type Donation struct {
	Donor      common.Address
	Amount     *big.Int
	ProjectID  string
	Timestamp  *big.Int
	TxHash     string
}

// NewGoFundChainContract 创建新的合约实例
func NewGoFundChainContract(rpcURL, contractAddress string) (*GoFundChainContract, error) {
	// 连接到以太坊客户端
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接以太坊客户端失败: %v", err)
	}

	// 解析合约地址
	address := common.HexToAddress(contractAddress)

	// 解析合约ABI
	contractABI, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, fmt.Errorf("解析合约ABI失败: %v", err)
	}

	// 创建绑定合约
	contract := bind.NewBoundContract(address, contractABI, client, client, client)

	return &GoFundChainContract{
		client:   client,
		contract: contract,
		address:  address,
		abi:      contractABI,
	}, nil
}

// CreateProject 创建众筹项目
func (gfc *GoFundChainContract) CreateProject(privateKey *ecdsa.PrivateKey, title, description, ipfsHash string, goal *big.Int, deadline *big.Int) (string, error) {
	// 创建认证对象
	auth, err := gfc.createAuth(privateKey)
	if err != nil {
		return "", fmt.Errorf("创建认证失败: %v", err)
	}

	// 准备输入参数
	input, err := gfc.abi.Pack("createProject", title, description, ipfsHash, goal, deadline)
	if err != nil {
		return "", fmt.Errorf("打包输入参数失败: %v", err)
	}

	// 发送交易
	tx, err := gfc.client.SendTransaction(context.Background(), &types.Transaction{
		To:       &gfc.address,
		Data:     input,
		Gas:      300000, // 预估gas
		GasPrice: big.NewInt(20000000000), // 20 Gwei
		Value:    big.NewInt(0),
		Nonce:    auth.Nonce.Uint64(),
	})
	if err != nil {
		return "", fmt.Errorf("发送交易失败: %v", err)
	}

	return tx.Hash().Hex(), nil
}

// Donate 捐赠到项目
func (gfc *GoFundChainContract) Donate(privateKey *ecdsa.PrivateKey, projectID string, amount *big.Int) (string, error) {
	// 创建认证对象
	auth, err := gfc.createAuth(privateKey)
	if err != nil {
		return "", fmt.Errorf("创建认证失败: %v", err)
	}

	// 准备输入参数
	input, err := gfc.abi.Pack("donate", projectID)
	if err != nil {
		return "", fmt.Errorf("打包输入参数失败: %v", err)
	}

	// 发送交易
	tx, err := gfc.client.SendTransaction(context.Background(), &types.Transaction{
		To:       &gfc.address,
		Data:     input,
		Gas:      200000, // 预估gas
		GasPrice: big.NewInt(20000000000), // 20 Gwei
		Value:    amount, // 捐赠金额
		Nonce:    auth.Nonce.Uint64(),
	})
	if err != nil {
		return "", fmt.Errorf("发送捐赠交易失败: %v", err)
	}

	return tx.Hash().Hex(), nil
}

// GetProject 获取项目信息
func (gfc *GoFundChainContract) GetProject(projectID string) (*Project, error) {
	// 准备输入参数
	input, err := gfc.abi.Pack("getProject", projectID)
	if err != nil {
		return nil, fmt.Errorf("打包输入参数失败: %v", err)
	}

	// 调用合约方法
	msg := ethereum.CallMsg{
		To:   &gfc.address,
		Data: input,
	}

	result, err := gfc.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("调用合约失败: %v", err)
	}

	// 解析返回结果
	var project Project
	err = gfc.abi.UnpackIntoInterface(&project, "getProject", result)
	if err != nil {
		return nil, fmt.Errorf("解析项目信息失败: %v", err)
	}

	return &project, nil
}

// GetProjectDonations 获取项目捐赠列表
func (gfc *GoFundChainContract) GetProjectDonations(projectID string) ([]Donation, error) {
	// 准备输入参数
	input, err := gfc.abi.Pack("getProjectDonations", projectID)
	if err != nil {
		return nil, fmt.Errorf("打包输入参数失败: %v", err)
	}

	// 调用合约方法
	msg := ethereum.CallMsg{
		To:   &gfc.address,
		Data: input,
	}

	result, err := gfc.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("调用合约失败: %v", err)
	}

	// 解析返回结果
	var donations []Donation
	err = gfc.abi.UnpackIntoInterface(&donations, "getProjectDonations", result)
	if err != nil {
		return nil, fmt.Errorf("解析捐赠列表失败: %v", err)
	}

	return donations, nil
}

// WithdrawFunds 项目创建者提取资金
func (gfc *GoFundChainContract) WithdrawFunds(privateKey *ecdsa.PrivateKey, projectID string) (string, error) {
	// 创建认证对象
	auth, err := gfc.createAuth(privateKey)
	if err != nil {
		return "", fmt.Errorf("创建认证失败: %v", err)
	}

	// 准备输入参数
	input, err := gfc.abi.Pack("withdrawFunds", projectID)
	if err != nil {
		return "", fmt.Errorf("打包输入参数失败: %v", err)
	}

	// 发送交易
	tx, err := gfc.client.SendTransaction(context.Background(), &types.Transaction{
		To:       &gfc.address,
		Data:     input,
		Gas:      150000, // 预估gas
		GasPrice: big.NewInt(20000000000), // 20 Gwei
		Value:    big.NewInt(0),
		Nonce:    auth.Nonce.Uint64(),
	})
	if err != nil {
		return "", fmt.Errorf("发送提取交易失败: %v", err)
	}

	return tx.Hash().Hex(), nil
}

// createAuth 创建交易认证对象
func (gfc *GoFundChainContract) createAuth(privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("无法获取公钥")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := gfc.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("获取nonce失败: %v", err)
	}

	gasPrice, err := gfc.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取gas价格失败: %v", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	return auth, nil
}

// WaitForTransaction 等待交易确认
func (gfc *GoFundChainContract) WaitForTransaction(txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	
	// 等待交易被挖矿
	receipt, err := bind.WaitMined(context.Background(), gfc.client, &types.Transaction{
		Hash: hash,
	})
	if err != nil {
		return nil, fmt.Errorf("等待交易确认失败: %v", err)
	}

	if receipt.Status == 0 {
		return nil, fmt.Errorf("交易执行失败")
	}

	return receipt, nil
}

// GetBalance 获取合约余额
func (gfc *GoFundChainContract) GetBalance() (*big.Int, error) {
	balance, err := gfc.client.BalanceAt(context.Background(), gfc.address, nil)
	if err != nil {
		return nil, fmt.Errorf("获取合约余额失败: %v", err)
	}
	return balance, nil
}

// GetGasPrice 获取当前gas价格
func (gfc *GoFundChainContract) GetGasPrice() (*big.Int, error) {
	gasPrice, err := gfc.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取gas价格失败: %v", err)
	}
	return gasPrice, nil
}

// Close 关闭客户端连接
func (gfc *GoFundChainContract) Close() {
	if gfc.client != nil {
		gfc.client.Close()
	}
} 