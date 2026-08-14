# Task List — dsh-android-apk

## Phase 1: 项目脚手架

- [ ] 1.1 创建 Flutter 项目结构 (`flutter create --org com.deepseek --project-name dsh_android`)
- [ ] 1.2 配置 Go 模块 (`go mod init github.com/deepseek/dsh-android/go-runtime`)
- [ ] 1.3 初始化 Android 项目配置 (minSdk 24, targetSdk 34)
- [ ] 1.4 配置 go_mobile 构建脚本 (`.gobuild.sh`)
- [ ] 1.5 添加依赖：riverpod, sqflite, http, shared_preferences

## Phase 2: Go Runtime 核心

- [ ] 2.1 实现 `go/runtime/session/types.go` — 会话事件类型定义
- [ ] 2.2 实现 `go/runtime/config/store.go` — 本地配置读写
- [ ] 2.3 实现 `go/runtime/session/store.go` — SQLite 会话持久化
- [ ] 2.4 实现 `go/runtime/tools/interface.go` — 工具接口定义
- [ ] 2.5 实现 `go/runtime/tools/registry.go` — 工具注册表
- [ ] 2.6 实现 `go/runtime/tools/bash.go` — Bash 命令执行工具
- [ ] 2.7 实现 `go/runtime/tools/fs.go` — 文件读写搜索工具
- [ ] 2.8 实现 `go/runtime/agent/types.go` — Agent 事件类型
- [ ] 2.9 实现 `go/runtime/agent/loop.go` — Agent 循环核心逻辑
- [ ] 2.10 编写 Go 单元测试覆盖 Phase 2 所有模块

## Phase 3: Go Mobile Bind 暴露

- [ ] 3.1 实现 `go/runtime/mobile/agent.go` — AgentLoop 的 gomobile bind 接口
- [ ] 3.2 实现 `go/runtime/mobile/session.go` — SessionStore 的 gomobile bind 接口
- [ ] 3.3 实现 `go/runtime/mobile/config.go` — ConfigStore 的 gomobile bind 接口
- [ ] 3.4 编写 `gomobile bind` 构建脚本生成 AAR
- [ ] 3.5 验证 AAR 可在 Android 项目中正常导入

## Phase 4: Flutter UI 层

- [ ] 4.1 实现 `lib/main.dart` — 应用入口和路由配置
- [ ] 4.2 实现 `lib/providers/session_provider.dart` — 会话状态管理
- [ ] 4.3 实现 `lib/providers/config_provider.dart` — 配置状态管理
- [ ] 4.4 实现 `lib/bridge/go_bridge.dart` — Method Channel 封装
- [ ] 4.5 实现 `lib/views/home_view.dart` — 主页（会话列表入口）
- [ ] 4.6 实现 `lib/views/chat_view.dart` — 对话界面
- [ ] 4.7 实现 `lib/views/timeline_panel.dart` — 执行过程面板
- [ ] 4.8 实现 `lib/views/settings_view.dart` — 设置页面
- [ ] 4.9 实现 `lib/widgets/message_bubble.dart` — 消息气泡组件
- [ ] 4.10 实现 `lib/widgets/tool_call_card.dart` — 工具调用卡片组件

## Phase 5: 集成与联调

- [ ] 5.1 打通 Flutter → Go 方法通道调用
- [ ] 5.2 实现 Agent 事件流从 Go → Flutter 的增量推送
- [ ] 5.3 实现会话创建/切换/删除的端到端流程
- [ ] 5.4 实现工具调用在时间线面板的实时展示
- [ ] 5.5 实现取消按钮对 Agent loop 的中断信号传递
- [ ] 5.6 实现配置修改的持久化和读取

## Phase 6: 测试与优化

- [ ] 6.1 运行 Dart 单元测试 (`flutter test`)
- [ ] 6.2 运行 Flutter 集成测试 (`flutter test integration_test/`)
- [ ] 6.3 真机部署验证 (Android 10+)
- [ ] 6.4 APK 体积优化 (代码混淆、资源压缩)
- [ ] 6.5 启动性能测试 (冷启动 < 3s)
