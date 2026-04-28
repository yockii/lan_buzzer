# LAN Buzzer

局域网抢答系统 - 支持 5-10 人实时抢答的简单工具

## 功能特点

- 🎯 实时抢答，毫秒级响应
- 👥 支持 5-10 人同时参与
- 💻 电脑、手机均可使用
- 🎨 颜色标识，支持重名
- ⚡ 单文件运行，无需安装依赖
- 🔒 局域网通信，安全可靠

## 🆢 题库模式（新功能）

想要题库答题功能？查看 **[feature/question-bank](https://github.com/yockii/lan_buzzer/tree/feature/question-bank)** 分支！

**题库模式额外功能**：
- ✅ 支持题库文件（单选、判断、问答）
- 🤖 自动判断单选题和判断题
- 👆 主持人可手动覆盖任何答案
- 👑 自动显示最先答对的选手
- 📊 3列优化布局（选手两侧，题目中间）
- 🔄 题目去重，避免重复

**使用题库模式**：
```bash
# 切换到题库分支
git checkout feature/question-bank

# 或直接下载预编译版本
# 访问：https://github.com/yockii/lan_buzzer/tree/feature/question-bank
# 查看详细文档
```

## 使用方法

1. 下载 `lan-buzzer.exe`（Windows）或 `lan-buzzer`（Linux/Mac）
2. 双击运行
3. 浏览器自动打开主持人界面
4. 选手通过二维码或 URL 加入
5. 点击"开始抢答"开始游戏

## 系统要求

- Windows 10+ / Linux / macOS
- 同一 Wi-Fi 局域网
- 现代浏览器（Chrome、Firefox、Safari）

## 技术栈

- 后端: Go 1.25+ + Fiber v2 + WebSocket
- 前端: Vue 3 + TypeScript + TailwindCSS

## 开发

```bash
# 构建
./build.sh

# 运行
./build/lan-buzzer.exe
```

## 分支说明

- **main**：主分支，基础抢答功能
- **feature/question-bank**：题库功能分支，包含完整的答题模式

## License

MIT
