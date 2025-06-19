创建一个 IPFS 网络取决于你的目标：你是想搭建一个**私有网络（Private IPFS Network）**，还是加入/使用**公共网络**？

下面我会分别介绍这两种情况：

---

## 🧭 目标 1：创建自己的 IPFS 网络（私有网络）

### ✅ 适合场景：

* 企业内网、实验室网络、私有链存储、开发环境等
* 不想让文件被广播到整个公网 IPFS

---

### 🧰 步骤一：安装 IPFS

IPFS 提供了官方的 Go 语言实现。你可以用以下命令安装：

```bash
wget https://dist.ipfs.tech/kubo/v0.27.0/kubo_v0.27.0_linux-amd64.tar.gz
tar -xvzf kubo_*.tar.gz
cd kubo
sudo bash install.sh
```

检查是否安装成功：

```bash
ipfs --version
```

---

### 🧰 步骤二：初始化节点

```bash
ipfs init
```

它会在用户目录生成 `.ipfs` 文件夹，包含私钥、默认配置等。

---

### 🧰 步骤三：设置私有网络（重点）

1. **生成 swarm key** 文件：

```bash
# 推荐使用 go-ipfs-swarm-key-gen 工具，也可手动生成
echo "/key/swarm/psk/1.0.0/" > ~/.ipfs/swarm.key
echo "/base16/" >> ~/.ipfs/swarm.key
openssl rand -hex 32 >> ~/.ipfs/swarm.key
```

内容大致如下：

```
/key/swarm/psk/1.0.0/
/base16/
90ff9bfa34055f24b1cb3617b12f...
```

2. **复制该 swarm.key 到所有节点的 \~/.ipfs 文件夹下**

📌 所有节点使用相同的 swarm.key 才能互通，其他节点无法加入。

---

### 🧰 步骤四：修改配置（关闭默认 bootstrap 节点）

```bash
ipfs bootstrap rm --all
```

你可以添加自己的节点地址作为 bootstrap：

```bash
ipfs bootstrap add /ip4/192.168.1.2/tcp/4001/p2p/<PeerID>
```

查看本节点 Peer ID：

```bash
ipfs id
```

---

### 🧰 步骤五：启动节点

```bash
ipfs daemon
```

节点会启动并监听本地 API：

* API: `localhost:5001`
* Gateway: `localhost:8080`
* Swarm: `4001`

---

## 🌐 多节点通信测试

1. 两台机器都配置好 swarm.key
2. 两台都运行 `ipfs daemon`
3. 在一台上添加另一台为 peer：

   ```bash
   ipfs swarm connect /ip4/192.168.1.2/tcp/4001/p2p/<PeerID>
   ```
4. 成功后，可以上传文件，一台上传，另一台 `ipfs cat <CID>` 查看

---

## ✅ 验证你是私有网络的方法：

* 用你自己的 CID，去 `https://ipfs.io/ipfs/<CID>`，**无法访问**，说明你并不在公网网络中 ✅
* 你只和同样 swarm.key 的节点互联

---

## 🧭 目标 2：加入公共 IPFS 网络（无需创建）

只需要安装 IPFS：

```bash
ipfs init
ipfs daemon
```

默认会连接 IPFS 官方的 bootstrap 节点列表，你的 CID 可以被全球访问：

```
https://ipfs.io/ipfs/<your-cid>
```

---

## 🧩 对比：私有 vs 公网 IPFS

| 特性    | 公共 IPFS        | 私有 IPFS               |
| ----- | -------------- | --------------------- |
| 可接入性  | 所有 IPFS 节点都能访问 | 仅配置相同 `swarm.key` 的节点 |
| 文件可用性 | 越多人使用越稳定       | 需自己维护节点和副本            |
| 安全性   | 数据被全球共享        | 数据只在你控制的节点流通          |
| 使用场景  | Web3、NFT、公链资源  | 企业内网、隐私存储、开发环境        |

---

## 📌 开发者建议

如果你是 DApp 开发者、众筹平台、NFT 系统：

* ✅ **开发阶段用私有 IPFS 网络**（方便调试，快）
* ✅ **上线时使用公共网络 + Pin 服务（如 Pinata）**，保证文件稳定可访问

---

