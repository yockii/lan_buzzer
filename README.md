# LAN Buzzer

局域网抢答系统 - 支持 5-10 人实时抢答的简单工具

## 功能特点

- 🎯 实时抢答，毫秒级响应
- 👥 支持 5-10 人同时参与
- 💻 电脑、手机均可使用
- 🎨 颜色标识，支持重名
- ⚡ 单文件运行，无需安装依赖
- 🔒 局域网通信，安全可靠

## 游戏模式

系统支持两种模式：

### 抢答模式（默认）
传统抢答系统，选手抢答看谁手速快。

### 答题模式
当存在 `questions.txt` 文件时，系统自动进入答题模式：

1. 在 `lan-buzzer.exe` 同目录下创建 `questions.txt` 文件
2. 格式说明：
   ```
   [单选]
   问题|选项A|选项B|选项C|选项D|正确答案

   [判断]
   问题|正确答案

   [问答]
   问题|正确答案
   ```
3. 运行 `lan-buzzer.exe`
4. 主持人点击"开始答题"开始
5. 选手在设备上看到题目并作答
6. 单选/判断题：自动判定
7. 问答题：主持人点击答案判定（待确认 → ✓ → ✗）
8. 第一个答对的选手显示 👑
9. 主持人点击"下一题"继续

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

## License

MIT
