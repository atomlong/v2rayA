#!/bin/bash

# 1. 暂停服务
echo "Stopping v2raya.service..."
sudo systemctl stop v2raya.service

# 强制清理可能残留的 v2raya 进程 (防止端口占用)
echo "Killing lingering v2raya processes..."
sudo killall v2raya 2>/dev/null

# 等待进程完全退出
TIMEOUT=10
while sudo pgrep v2raya >/dev/null; do
    if [ $TIMEOUT -le 0 ]; then
        echo "Force killing v2raya..."
        sudo killall -9 v2raya 2>/dev/null
        break
    fi
    echo "Waiting for v2raya to exit..."
    sleep 1
    ((TIMEOUT--))
done

# 2. 编译
echo "Compiling v2raya..."
cd service
if ! go build -o v2raya .; then
    echo "Compilation failed!"
    sudo systemctl start v2raya.service
    exit 1
fi
cd ..

# 3. 准备环境 & 运行
# 加载 EnvironmentFile (如果存在)
if [ -f /etc/default/v2raya ]; then
    echo "Loading /etc/default/v2raya..."
    set -a
    source /etc/default/v2raya
    set +a
fi

# 设置 service 文件中定义的环境变量
export V2RAYA_LOG_FILE=/var/log/v2raya/v2raya.log
# 强制指定资源目录，复用已安装的 Xray 资源文件 (位于 /root/.local/share/xray)
export XRAY_LOCATION_ASSET=/root/.local/share/xray

# 清理旧日志，避免误报
if sudo test -f "$V2RAYA_LOG_FILE"; then
    echo "Clearing old log file..."
    sudo truncate -s 0 "$V2RAYA_LOG_FILE"
fi

echo "Running compiled v2raya..."
# 使用 -E 保留加载的配置，但强制重置 HOME 为 /root，并清空 XDG_RUNTIME_DIR
# 重定向 stdout/stderr 到文件以便调试
sudo -E HOME=/root XDG_RUNTIME_DIR= V2RAYA_LOG_FILE="$V2RAYA_LOG_FILE" XRAY_LOCATION_ASSET="$XRAY_LOCATION_ASSET" ./service/v2raya --log-disable-timestamp > v2raya_stdout.log 2>&1 &
PID=$!

echo "v2raya started with PID $PID."

# 倒计时逻辑
DURATION=1200
DIED=0
while [ $DURATION -gt 0 ]; do
    # 1. 检查主进程是否存活
    if ! sudo kill -0 $PID 2>/dev/null; then
        echo -e "\nERROR: v2raya process died unexpectedly!"
        DIED=1
        break
    fi

    # 2. 检查日志中核心进程是否崩溃 (v2ray-core: exit status 或 unexpected exiting)
    if sudo grep -q -E "v2ray-core: exit status|unexpected exiting" "$V2RAYA_LOG_FILE"; then
        echo -e "\nERROR: Xray core process crashed! (Detected in logs)"
        # 即使主进程还在，核心挂了也算失败
        DIED=1
        break
    fi
    
    # 打印倒计时
    echo -ne "Monitoring... ${DURATION}s remaining \r"
    sleep 1
    ((DURATION--))
done
echo ""

# 4. 检测进程状态
if [ $DIED -eq 1 ]; then
    echo "--- v2raya stdout/stderr tail ---"
    tail -n 20 v2raya_stdout.log
    echo "---------------------------------"

    # 5. 异常提取日志和配置
    echo "Extracting logs and config to current directory..."
    if sudo test -f "$V2RAYA_LOG_FILE"; then
        sudo cp "$V2RAYA_LOG_FILE" ./v2raya_crash.log
        sudo chown $(id -u):$(id -g) ./v2raya_crash.log
        echo "Saved v2raya_crash.log"
        
        echo "--- v2raya log file tail ---"
        tail -n 20 ./v2raya_crash.log
        echo "----------------------------"
    else
        echo "Log file not found at $V2RAYA_LOG_FILE"
    fi
    
    if [ -f "/etc/v2raya/config.json" ]; then
        sudo cp "/etc/v2raya/config.json" ./config_crash.json
        sudo chown $(id -u):$(id -g) ./config_crash.json
        echo "Saved config_crash.json"
    fi
else
    echo "SUCCESS: v2raya is running stable."
    sudo kill $PID
fi

# 6. 重启服务
echo "Restarting v2raya.service..."
sudo systemctl start v2raya.service
