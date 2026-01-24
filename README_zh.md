# v2rayA [![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/v2rayA/v2raya)](https://hub.docker.com/r/mzz2017/v2raya) [![Travis (.org)](https://img.shields.io/travis/v2rayA/v2rayA?label=travis-ci%20build)](https://travis-ci.org/v2rayA/v2rayA)

[**English**](https://github.com/v2rayA/v2rayA/blob/main/README.md)&nbsp;&nbsp;&nbsp;[**简体中文**](https://github.com/v2rayA/v2rayA/blob/main/README_zh.md)

v2rayA 是一个支持全局透明代理的 V2Ray 客户端，同时兼容 SS、SSR、Trojan(trojan-go)、Tuic 与 [Juicity](https://github.com/juicity)协议。 [[SSR支持清单]](https://github.com/v2rayA/dist/shadowsocksR/blob/master/README.md#ss-encrypting-algorithm)

v2rayA 致力于提供最简单的操作，满足绝大部分需求。

得益于 Web 客户端的优势，你不仅可以将其用于本地计算机，还可以轻松地将它部署在路由器或 NAS 上。

项目地址：https://github.com/v2rayA/v2rayA


## 使用方法

v2rayA 主要提供了下述使用方法：

1. 从 APT 软件源或者 AUR 安装
2. Docker
3. 自建 [scoop bucket](https://github.com/v2rayA/v2raya-scoop) (Windows 用户)
4. 自建 [homebrew tap](https://github.com/v2rayA/homebrew-v2raya)
5. 自建 [OpenWrt 仓库](https://github.com/v2rayA/v2raya-openwrt) 和 OpenWrt 官方软件源（从 OpenWrt 22.03 版本开始提供）
6. 微软 winget: https://winstall.app/apps/v2rayA.v2rayA
7. Ubuntu Snap: https://snapcraft.io/v2raya
8. 从 GitHub releases 下载二进制与安装包

详见 [**v2rayA - Docs**](https://v2raya.org/docs/prologue/introduction/)


## 开发者测试

如果你修改了源码并想使用 systemd 服务进行测试，请按照以下步骤操作：

### 1. 编译新版本

```bash
cd service
go build -o v2raya .
```

### 2. 替换系统中的二进制文件

```bash
# 停止当前运行的服务
sudo systemctl stop v2raya

# 备份原始二进制文件
sudo cp /usr/bin/v2raya /usr/bin/v2raya.backup

# 复制新编译的二进制文件
sudo cp ./v2raya /usr/bin/v2raya

# 确保执行权限
sudo chmod +x /usr/bin/v2raya
```

### 3. 启动服务

```bash
sudo systemctl start v2raya
```

### 4. 验证修改

1. 打开浏览器访问 `http://localhost:2017`（默认端口）
2. 进行相应的功能测试
3. 验证修改是否生效

### 5. 查看日志

```bash
# 实时查看服务日志
sudo journalctl -u v2raya -f

# 或查看日志文件
sudo tail -f /var/log/v2raya/v2raya.log
```

### 6. 恢复原版本（可选）

测试完成后，如需恢复原版本：

```bash
sudo systemctl stop v2raya
sudo cp /usr/bin/v2raya.backup /usr/bin/v2raya
sudo systemctl start v2raya
```


## 界面截图

<img src="https://i.loli.net/2020/04/19/kp2oedPiSzVwgHJ.png" border="0">


## 注意

1. 程序不会将任何用户数据保存在云端，所有用户数据存放在用户本地配置文件中。

2. **不要将本项目用于不合法用途。**

## 感谢

[hq450/fancyss](https://github.com/hq450/fancyss)

[ToutyRater/v2ray-guide](https://github.com/ToutyRater/v2ray-guide/blob/master/routing/sitedata.md)

[nadoo/glider](https://github.com/nadoo/glider)

[Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)

[zfl9/ss-tproxy](https://github.com/zfl9/ss-tproxy/blob/master/ss-tproxy)

## 协议

[![License: AGPL v3-only](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
