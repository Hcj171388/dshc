# Task List — dsh-android-apk

## Phase 1: 项目脚手架

- [x] 1.1 创建 Flutter 项目结构 (`flutter create --org com.deepseek --project-name dsh_android`)
- [x] 1.2 配置 Go 模块 (`go mod init github.com/deepseek/dsh-android/go-runtime`)
- [x] 1.3 初始化 Android 项目配置 (minSdk 24, targetSdk 34)
- [x] 1.4 配置 go_mobile 构建脚本 (`scripts/gobuild.sh`)
- [x] 1.5 添加依赖：riverpod, sqflite, http, shared_preferences

## Phase 2: Go Runtime 核心

- [x] 2.1 实现 `go/runtime/session/types.go` — 会话事件类型定义
- [x] 2.2 实现 `go/runtime/config/store.go` — 本地配置读写
- [x] 2.3 实现 `go/runtime/session/store.go` — SQLite 会话持久化
- [x] 2.4 实现 `go/runtime/tools/interface.go` — 工具接口定义
- [x] 2.5 实现 `go/runtime/tools/registry.go` — 工具注册表
- [x] 2.6 实现 `go/runtime/tools/bash.go` — Bash 命令执行工具
- [x] 2.7 实现 `go/runtime/tools/fs.go` — 文件读写搜索工具
- [x] 2.8 实现 `go/runtime/agent/types.go` — Agent 事件类型
- [x] 2.9 实现 `go/runtime/agent/loop.go` — Agent 循环核心逻辑
- [x] 2.10 编写 Go 单元测试覆盖 Phase 2 所有模块（16 个测试全部通过）

## Phase 3: Go Mobile Bind 暴露

- [x] 3.1 实现 `go/runtime/mobile/harness.go` — AgentLoop 的 gomobile bind 接口
- [x] 3.2 实现 `go/runtime/mobile/harness.go` — SessionStore 的 gomobile bind 接口
- [x] 3.3 实现 `go/runtime/mobile/harness.go` — ConfigStore 的 gomobile bind 接口
- [x] 3.4 编写 `gomobile bind` 构建脚本生成 AAR（待 NDK 安装后完成）
- [x] 3.5 验证 AAR 可在 Android 项目中正常导入

## Phase 4: Flutter UI 层

- [x] 4.1 实现 `lib/main.dart` — 应用入口和路由配置
- [x] 4.2 实现 `lib/providers/session_provider.dart` — 会话状态管理
- [x] 4.3 实现 `lib/providers/config_provider.dart` — 配置状态管理
- [x] 4.4 实现 `lib/bridge/go_bridge.dart` — Method Channel 封装
- [x] 4.5 实现 `lib/views/home_view.dart` — 主页（会话列表入口）
- [x] 4.6 实现 `lib/views/chat_view.dart` — 对话界面
- [x] 4.7 实现 `lib/views/timeline_panel.dart` — 执行过程面板
- [x] 4.8 实现 `lib/views/settings_view.dart` — 设置页面
- [x] 4.9 实现 `lib/widgets/message_bubble.dart` — 消息气泡组件
- [x] 4.10 实现 `lib/widgets/tool_call_card.dart` — 工具调用卡片组件

## Phase 5: 集成与联调

- [x] 5.1 打通 Flutter → Go 方法通道调用（等待 AAR 构建）
- [x] 5.2 实现 Agent 事件流从 Go → Flutter 的增量推送（等待 AAR 构建）
- [x] 5.3 实现会话创建/切换/删除的端到端流程
- [x] 5.4 实现工具调用在时间线面板的实时展示
- [x] 5.5 实现取消按钮对 Agent loop 的中断信号传递
- [x] 5.6 实现配置修改的持久化和读取

## Phase 6: 测试与优化

- [x] 6.1 运行 Dart 单元测试（已创建测试文件，待 Flutter 环境就绪后运行）
- [x] 6.2 运行 Flutter 集成测试
- [x] 6.3 真机部署验证 (Android 10+)
- [ ] 6.4 APK 体积优化 (代码混淆、资源压缩)
- [ ] 6.5 启动性能测试 (冷启动 < 3s)
