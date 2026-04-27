# LAN Buzzer 设计文档

**项目名称**: lan-buzzer
**创建日期**: 2026-04-27
**目标**: 局域网抢答系统，支持主持人投屏展示，多选手实时抢答

## 1. 项目概述

### 1.1 使用场景
- 小型团队会议/培训（5-10人）
- 主持人电脑投屏到大屏幕
- 选手使用电脑或手机参与抢答
- 同一 Wi-Fi 局域网环境

### 1.2 核心功能
- 主持人控制抢答流程（开始、重置、下一题）
- 实时显示连接的选手信息
- 抢答结果大屏展示
- 支持键盘（空格/回车）和触摸屏操作
- 防抢跑机制（模态框时间惩罚）
- 选手颜色标识，支持重名

## 2. 技术架构

### 2.1 技术栈
- **后端**: Go 1.25+ + Fiber v3 + WebSocket
- **前端**: Vue 3 + TypeScript + Vite + shadcn-vue + TailwindCSS
- **通信**: WebSocket 实时双向通信
- **部署**: Go 编译成单个可执行文件，内嵌前端静态资源

### 2.2 项目结构
```
lan-buzzer/
├── backend/              # Go + Fiber v3
│   ├── main.go          # 入口，内嵌静态文件
│   ├── websocket/       # WebSocket 处理
│   ├── game/            # 游戏状态管理
│   └── embed/           # 前端静态资源（编译时嵌入）
├── frontend/            # Vue3 + TypeScript
│   ├── host/           # 主持人端页面
│   └── player/         # 选手端页面
└── build/              # 编译输出
```

### 2.3 路由设计
- `/` - 自动检测设备类型（电脑=主持人端，手机=选手端）
- `/host` - 强制主持人端
- `/player` - 强制选手端
- `/player?name=xxx` - 选手端（预填名字）
- `/ws` - WebSocket 连接端点

## 3. 核心功能设计

### 3.1 主持人端功能

#### 连接管理
- 显示已连接选手列表（名字、颜色、设备类型、连接状态）
- 右上角显示本机 IP 地址和二维码
- 移除选手功能（列表右侧"×"按钮）

#### 游戏控制
- **开始抢答**: 状态从 `waiting` → `ready`，记录开始时间
- **下一题**: 重置 winner，状态改回 `waiting`
- **取消开始**: 仅在开始 2 秒内可用，防止误操作

#### 信息展示
- **抢答结果**: 超大显示（64px）抢到的选手名字 + 颜色标识
- **选手列表**: 底部一行显示所有已连接选手
- **状态指示**: 顶部显示当前状态
- **连接信息**: 右上角 IP + 二维码（扫码直达选手端）

#### 界面布局（极简聚焦方案）
```
┌─────────────────────────────────────┐
│ 192.168.1.100:3000        📱 二维码  │
├─────────────────────────────────────┤
│                                      │
│        🎉 张三 (红色)                │
│           抢到了！                   │
│                                      │
├─────────────────────────────────────┤
│    [开始抢答]  [下一题]              │
├─────────────────────────────────────┤
│  已连接 5 位选手: 张三 李四 王五...  │
└─────────────────────────────────────┘
```

### 3.2 选手端功能

#### 注册流程（第一屏）
- 输入名字输入框
- "进入抢答"按钮
- 客户端验证：非空、无特殊字符

#### 抢答界面（第二屏）
- **状态显示**: 顶部显示当前状态（等待/可抢答/已锁定/抢到了）
- **抢答按钮**: 大按钮（120px 高），红色背景
- **操作方式**:
  - 手机：触摸按钮
  - 电脑：空格键 或 回车键
- **连接状态**: 显示"已连接到主持人"或断线提示

#### 界面布局
```
┌─────────────────────────┐
│  👋 张三 (红色)         │
│  ✓ 已连接               │
├─────────────────────────┤
│  ⏳ 等待主持人开始...    │
├─────────────────────────┤
│                         │
│     [ 抢答！ ]          │
│   (空格键 或 点击)      │
│                         │
└─────────────────────────┘
```

### 3.3 防抢跑机制
- **触发条件**: 在 `waiting` 状态下按键
- **惩罚方式**: 弹出模态框"请等待主持人开始！"
- **恢复条件**: 必须点击"我知道了"按钮
- **最小停留**: 模态框显示至少 1-2 秒（防止快速点击）

## 4. 数据模型

### 4.1 服务器端状态

```go
type GameState string

const (
    StateWaiting GameState = "waiting"  // 等待开始
    StateReady   GameState = "ready"    // 可抢答
    StateLocked  GameState = "locked"   // 已锁定
)

type Player struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Color       string    `json:"color"`
    DeviceType  string    `json:"deviceType"`  // "desktop" | "mobile"
    ConnectedAt time.Time `json:"connectedAt"`
    Conn        *websocket.Conn `json:"-"`
}

type GameServer struct {
    State     GameState              `json:"state"`
    Players   map[string]*Player     `json:"players"`
    WinnerID  *string                `json:"winnerId"`
    StartTime *time.Time             `json:"startTime"`
    Mutex     sync.RWMutex           `json:"-"`
}
```

### 4.2 颜色分配

```go
var PlayerColors = []string{
    "#ef4444", // 红
    "#f97316", // 橙
    "#eab308", // 黄
    "#22c55e", // 绿
    "#06b6d4", // 青
    "#3b82f6", // 蓝
    "#a855f7", // 紫
    "#ec4899", // 粉
}

// 按连接顺序循环分配
```

## 5. WebSocket 通信协议

### 5.1 消息格式

```typescript
// 基础消息结构
interface Message {
  type: string;
  payload: any;
}
```

### 5.2 选手 → 服务器

```typescript
// 加入游戏
{
  type: "join",
  payload: {
    name: string;
    deviceType: "desktop" | "mobile";
  }
}

// 抢答
{
  type: "buzz",
  payload: {
    timestamp: number;  // 客户端时间戳
  }
}
```

### 5.3 服务器 → 全员广播

```typescript
// 状态变更
{
  type: "state_changed",
  payload: {
    state: "waiting" | "ready" | "locked";
    startTime?: number;  // ready 状态时提供
  }
}

// 选手列表更新
{
  type: "player_list",
  payload: {
    players: Array<{
      id: string;
      name: string;
      color: string;
      deviceType: string;
    }>;
  }
}

// 抢答结果
{
  type: "buzz_result",
  payload: {
    winner: {
      id: string;
      name: string;
      color: string;
    } | null;
    isEarly: boolean;  // 是否抢跑
  }
}
```

### 5.4 服务器 → 个人

```typescript
// 错误消息
{
  type: "error",
  payload: {
    message: string;
  }
}

// 抢跑警告（触发模态框）
{
  type: "early_buzz_warning",
  payload: {
    message: string;
  }
}
```

## 6. 抢答逻辑流程

### 6.1 正常流程
1. **等待阶段** (`waiting`)
   - 选手可加入，显示在列表中
   - 抢答按钮禁用，显示"等待主持人开始"
   - 抢跑触发模态框惩罚

2. **开始抢答** (`ready`)
   - 主持人点击"开始抢答"
   - 服务器广播 `state_changed: ready`
   - 记录开始时间
   - 抢答按钮激活，显示"可抢答"

3. **选手抢答**
   - 收到第一个 `buzz` 消息
   - 记录 winner，状态改为 `locked`
   - 广播 `buzz_result`（winner 信息）
   - 后续抢答消息被忽略

4. **已锁定** (`locked`)
   - 显示抢答结果
   - 抢答按钮禁用
   - 等待主持人点击"下一题"

5. **下一题** → 返回步骤 1

### 6.2 网络异常处理

#### 选手断线
- 服务器检测心跳（30秒超时）
- 自动从玩家列表移除
- 通知主持人更新列表
- 不影响已确定的抢答结果

#### 主持人断线
- 所有选手显示"主持人已断开"
- 客户端尝试自动重连（指数退避）
- 重连成功后恢复状态

#### 抢答时断线
- 已抢到的选手断线不影响结果
- 只需在下次抢答前重新连接

## 7. 部署方案

### 7.1 构建流程

```bash
# 1. 构建前端
cd frontend
npm run build

# 2. Go 编译（自动嵌入前端）
cd ../backend
go build -o ../build/lan-buzzer.exe ./...

# 3. 运行
./build/lan-buzzer.exe
# 自动打开 http://localhost:3000
```

### 7.2 使用流程

1. **主持人**：
   - 双击运行 `lan-buzzer.exe`
   - 浏览器自动打开主持人端
   - 投屏到大屏幕

2. **选手**：
   - 手机扫描主持人端二维码
   - 或手动输入 URL（如 `http://192.168.1.100:3000/player`）
   - 输入名字，等待开始

3. **开始抢答**：
   - 主持人确认选手都已连接
   - 点击"开始抢答"
   - 选手按键抢答
   - 显示结果，点击"下一题"继续

## 8. 技术实现要点

### 8.1 后端关键点

- **并发安全**: 使用 `sync.RWMutex` 保护 GameServer 状态
- **WebSocket 管理**: 维护连接映射，处理断线重连
- **心跳检测**: 定期 ping，30秒无响应视为断线
- **静态资源嵌入**: 使用 `go:embed` 打包前端

### 8.2 前端关键点

- **设备检测**: 通过 User-Agent 判断设备类型
- **按键监听**: 全局监听空格/回车，防止重复触发
- **状态同步**: WebSocket 消息驱动 UI 更新
- **模态框**: 抢跑惩罚，最小停留时间
- **响应式设计**: 适配手机和电脑屏幕

### 8.3 性能考虑

- **消息去重**: 客户端防抖，避免重复发送
- **状态广播**: 使用单个 goroutine 广播，避免并发写
- **资源优化**: 前端按需加载，减少首屏时间

## 9. 未来扩展（可选）

- 抢答历史记录（哪一题谁抢到了）
- 积分系统（累计抢答成功次数）
- 声音效果（开始提示、抢答成功音）
- 多轮抢答模式（比如三局两胜）
- 导出结果（CSV/JSON）
- 皮肤主题切换

## 10. 验收标准

- [ ] 主持人能正常启动，浏览器自动打开
- [ ] 选手能通过二维码或 URL 连接
- [ ] 选手列表实时更新
- [ ] 抢答响应时间 < 100ms（局域网）
- [ ] 支持 5-10 人同时抢答
- [ ] 抢跑惩罚模态框正常工作
- [ ] 颜色分配正确，支持重名
- [ ] 断线重连机制正常
- [ ] 编译成单个可执行文件
- [ ] 零依赖部署（无需 Node.js）
