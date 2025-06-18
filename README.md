---

## 🧱 项目名称：`GoFundChain` —— 去中心化众筹平台

---

## 🎯 项目目标

一个允许用户：

* 发起众筹项目（上传描述、金额、时间等）
* 使用加密钱包地址发起和捐款
* 所有交易信息记录在链上（或链下+链上证明）
* 支持通过 Web 前端或 CLI 访问

---

## 🔧 技术架构

| 层级         | 技术                           | 说明             |
| ---------- | ---------------------------- | -------------- |
| 用户界面       | Web（Vue / React） or CLI      | 提供使用入口         |
| 服务层（Go 实现） | Go + Gin/Fiber + REST API    | 提供众筹创建/查询/捐款功能 |
| 数据层        | IPFS + SQLite/PostgreSQL     | 存储项目信息和本地索引    |
| 区块链        | 以太坊 or Cosmos SDK            | 记录关键交易和验真      |
| 钱包交互       | Go-Ethereum（`geth`）或 web3.js | 支持用户钱包认证、交易签名  |

---

## 🧠 核心模块设计（Go 实现）

### ✅ 1. 用户模块

* 注册/登录（使用钱包地址）
* 钱包绑定与签名验证

🔸 Go 中用 `github.com/ethereum/go-ethereum/crypto` 验签：

```go
func VerifySignature(pubKey, message, sig []byte) bool
```

---

### ✅ 2. 众筹项目模块

* 创建众筹（含目标金额、时间、描述）
* 存储项目至 IPFS（或本地数据库）
* 记录项目 hash 到链上

```go
type Project struct {
    ID        string
    Title     string
    Creator   string // 钱包地址
    Goal      float64
    Raised    float64
    Deadline  time.Time
    IPFSHash  string
}
```

---

### ✅ 3. 捐赠模块

* 用户捐赠资金（用 ETH / 链上代币）
* 后端记录交易信息
* 提交交易哈希 + 签名信息上链

```go
type Donation struct {
    TxHash     string
    Donor      string
    Amount     float64
    Timestamp  time.Time
    ProjectID  string
}
```

---

### ✅ 4. 区块链集成模块

* 用 `go-ethereum` 实现链上交互
* 可部署一个简单的智能合约，记录项目哈希与捐款人记录
* 合约示意：

```solidity
mapping(string => address[]) donations;
mapping(string => string) projectIPFSHash;
```

Go 调用合约：

```go
client, _ := ethclient.Dial("https://rpc.xxx")
contract, _ := NewCrowdContract(addr, client)
tx, _ := contract.Donate(auth, projectID)
```

---

### ✅ 5. 防篡改 & 验证模块

* 项目信息加签后存储
* 用户可通过项目 ID + 签名验证链上记录是否一致

---

## 🚀 Go 项目结构建议

```
go-fund-chain/
├── main.go
├── api/
│   ├── handlers.go     # API 路由处理
│   └── middleware.go
├── blockchain/
│   └── eth.go          # 区块链交互
├── ipfs/
│   └── client.go       # IPFS 上传/下载
├── models/
│   ├── project.go
│   └── donation.go
├── db/
│   └── database.go     # SQLite/Postgres操作
├── utils/
│   └── crypto.go       # 签名/验签等
└── README.md
```

---

## 🌍 项目部署建议

| 组件     | 部署方式                            |
| ------ | ------------------------------- |
| 后端 API | Docker + Go 二进制                 |
| 前端     | 静态网页 + IPFS                     |
| 区块链    | Ganache 本地测试 / Polygon zkEVM 主网 |
| IPFS   | 使用 Infura 或自建节点                 |
| 钱包连接   | 使用 MetaMask 或 WalletConnect     |

---

## 🔑 项目扩展方向

* DAO 治理（用户投票决定项目是否提现）
* 多签钱包控制资金
* NFT 作为贡献凭证
* zk 证明用于隐私保护的捐赠

---

## 🎁 可参考开源项目

* [https://github.com/ConsenSys/go-blockchain-fund](https://github.com/ConsenSys/go-blockchain-fund)
* [https://github.com/ipfs/go-ipfs-api](https://github.com/ipfs/go-ipfs-api)
* [https://github.com/ethereum/go-ethereum](https://github.com/ethereum/go-ethereum)

---

## 🎯 如果你想一步步做：

我可以帮你：

* ✅ 生成初始化代码骨架
* ✅ 编写某个模块的实现（如 IPFS 上传 / 智能合约交互）
* ✅ 本地调试指南（Ganache + go-ethereum）

你现在想从哪个模块开始搭建？我可以帮你写第一个完整模块的代码。
