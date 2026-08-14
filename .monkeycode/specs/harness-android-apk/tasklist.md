# Task List — dsh-android-apk

## Phase 1: 项目脚手架

- [x] 1.1 创建 Flutter 项目结构
- [x] 1.2 配置 Go 模块
- [x] 1.3 初始化 Android 项目配置
- [x] 1.4 配置 go_mobile 构建脚本
- [x] 1.5 添加依赖

## Phase 2: Go Runtime 核心

- [x] 2.1 实现 session/types.go
- [x] 2.2 实现 config/store.go
- [x] 2.3 实现 session/store.go
- [x] 2.4 实现 tools/interface.go
- [x] 2.5 实现 tools/registry.go
- [x] 2.6 实现 tools/bash.go
- [x] 2.7 实现 tools/fs.go
- [x] 2.8 实现 agent/types.go
- [x] 2.9 实现 agent/loop.go (完整 LLM 循环)
- [x] 2.10 实现 llm/client.go (DeepSeek API 客户端)

## Phase 3: 扩展工具

- [x] 3.1 实现 tools/web.go (web_search, web_fetch)
- [x] 3.2 实现 tools/terminal.go
- [x] 3.3 实现 tools/todo.go
- [x] 3.4 实现 mobile/harness.go

## Phase 4: Flutter UI 层

- [x] 4.1 实现 lib/main.dart
- [x] 4.2 实现 lib/providers/session_provider.dart
- [x] 4.3 实现 lib/providers/config_provider.dart
- [x] 4.4 实现 lib/bridge/go_bridge.dart
- [x] 4.5 实现 lib/views/home_view.dart
- [x] 4.6 实现 lib/views/chat_view.dart
- [x] 4.7 实现 lib/views/timeline_panel.dart
- [x] 4.8 实现 lib/views/settings_view.dart
- [x] 4.9 实现 lib/widgets/message_bubble.dart
- [x] 4.10 实现 lib/widgets/tool_call_card.dart

## Phase 5: 集成与联调

- [x] 5.1 打通 Flutter → Go 方法通道调用
- [x] 5.2 实现 Agent 事件流从 Go → Flutter
- [x] 5.3 实现会话创建/切换/删除的端到端流程
- [x] 5.4 实现工具调用在时间线面板的实时展示
- [x] 5.5 实现取消按钮对 Agent loop 的中断信号传递
- [x] 5.6 实现配置修改的持久化和读取

## Phase 6: 测试与优化

- [x] 6.1 运行 Dart 单元测试
- [x] 6.2 构建 APK
- [x] 6.3 APK 上传到 GitHub Releases

## 已实现功能

### Agent 运行时
- [x] 完整的 Agent Loop 支持工具调用
- [x] LLM 集成 (DeepSeek API)
- [x] 多轮对话支持
- [x] 取消/中止功能

### 工具系统
- [x] bash - 执行 bash 命令
- [x] read_file - 读取文件
- [x] write_file - 写入文件
- [x] list_dir - 列出目录
- [x] search_files - 搜索文件
- [x] web_search - 搜索网页
- [x] web_fetch - 获取网页内容
- [x] terminal - 持久化终端会话
- [x] todo_write - Todo 列表管理

### 会话管理
- [x] SQLite 持久化
- [x] 会话列表
- [x] 会话删除/归档
- [x] 事件历史

### UI
- [x] 聊天界面
- [x] 执行时间线
- [x] 会话管理
- [x] 设置页面

### 下载
- https://github.com/Hcj171388/dshc/releases/tag/v1.0.0
