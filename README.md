# moacchain-adapter

本项目适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf文件，新建MOAC.ini文件，编辑如下内容：

```ini

# node api url
nodeAPI = "http://ip:port"

# gas limit
gasLimit = "10000"

# is testnet
isTestnet = true

```
