# Requirements Document — dsh-android-apk

## Introduction

将 DeepSeek Harness（Node.js/TypeScript 插件化 Agent 运行时）重构为 Android APK。前端使用 Flutter，后端核心运行时使用 Go（Go Mobile），所有数据和配置纯本地存储于设备，不依赖远程服务。提供混合交互模式：用户可通过对话界面输入指令，也可观看 Agent 自动执行任务的过程与结果。

## Glossary

- **dsh-android**: 目标 Android 应用，由 Flutter + Go Mobile 构建
- **Agent Runtime**: 核心 Agent 循环，负责解析用户意图、调度工具、管理会话状态
- **Tool Bridge**: Go 侧实现的工具接口，向 Agent Runtime 暴露 bash、文件操作等能力
- **Session Store**: 本地持久化的会话日志，使用 SQLite 存储
- **Config Store**: 本地配置文件，存储 Agent 行为参数和系统设置
- **Method Channel**: Flutter 与 Go 之间的双向通信机制

## Requirements

### Requirement 1: Agent 运行时

**User Story:** AS 用户, I want an agent that can execute tools autonomously, so that I can complete complex multi-step tasks with minimal interaction.

#### Acceptance Criteria

1. WHEN user submits a text prompt, the system SHALL launch an Agent loop that decomposes the task into tool calls
2. WHILE the Agent loop is running, the system SHALL emit incremental state events (tool call start, tool call result, turn complete)
3. IF the Agent loop encounters a tool error, the system SHALL capture the error and let the Agent decide the recovery action
4. WHILE the Agent loop is running, the system SHALL respect the user's cancellation signal and terminate gracefully within 2 seconds
5. IF the user opens the app after a previous session, the system SHALL restore the last active session automatically

---

### Requirement 2: 工具执行能力

**User Story:** AS Agent Runtime, I want access to shell execution and file operations, so that I can manipulate the device environment to accomplish user tasks.

#### Acceptance Criteria

1. WHEN the Agent calls `bash` tool, the system SHALL execute the command in a sandboxed local shell and return stdout/stderr
2. WHEN the Agent calls `read_file` tool, the system SHALL read the file from local storage and return its content
3. WHEN the Agent calls `write_file` tool, the system SHALL write content to the specified local path
4. WHEN the Agent calls `edit_file` tool, the system SHALL apply the requested edit operation and return the diff
5. WHEN the Agent calls `search_files` tool, the system SHALL search the local filesystem and return matching paths
6. IF a tool call exceeds the configured timeout, the system SHALL abort the tool and return a timeout error

---

### Requirement 3: 会话持久化

**User Story:** AS 用户, I want my conversations saved locally, so that I can resume them later without losing context.

#### Acceptance Criteria

1. AFTER each user message, the system SHALL append the event to the local session log
2. AFTER each Agent turn completes, the system SHALL flush buffered events to persistent storage
3. WHEN the user opens a past session, the system SHALL reconstruct the full conversation history from the local log
4. IF the session log is corrupted, the system SHALL attempt repair and fall back to the last valid checkpoint
5. WHEN the user deletes a session, the system SHALL remove all log entries and associated data from local storage

---

### Requirement 4: 本地配置管理

**User Story:** AS 用户, I want to configure the Agent behavior without any remote dependency, so that the app works fully offline.

#### Acceptance Criteria

1. WHEN the app launches for the first time, the system SHALL create a default configuration file in the app's private directory
2. WHEN the user modifies settings in the UI, the system SHALL write changes to the local config file immediately
3. WHEN the Agent initializes, the system SHALL read configuration from the local config file and apply it
4. IF the config file is missing or invalid, the system SHALL use built-in defaults and log a warning

---

### Requirement 5: 对话式 UI

**User Story:** AS 用户, I want a chat interface to send prompts and view responses, so that I can interact with the Agent naturally.

#### Acceptance Criteria

1. WHEN the user types a message and taps send, the system SHALL display the message in the chat timeline
2. WHILE the Agent is processing, the system SHALL show a typing indicator
3. WHEN the Agent produces a response, the system SHALL stream the response text into the chat timeline
4. WHEN the Agent invokes a tool, the system SHALL show a collapsible tool-call card with name and arguments
5. WHEN the tool result is available, the system SHALL expand the card to show the result summary

---

### Requirement 6: Agent 执行过程可视化

**User Story:** AS 用户, I want to watch the Agent execute steps in real time, so that I can understand what the Agent is doing and intervene if needed.

#### Acceptance Criteria

1. WHEN the Agent starts a tool execution, the system SHALL add a progress item to the execution timeline panel
2. WHEN a tool completes, the system SHALL update the corresponding timeline item with success or failure status
3. WHEN the Agent completes a turn, the system SHALL collapse the timeline item and mark it as finished
4. WHEN the user taps a completed timeline item, the system SHALL show the full tool input and output
5. WHILE the Agent is running, the system SHALL allow the user to tap a cancel button to abort the current turn

---

### Requirement 7: 多会话管理

**User Story:** AS 用户, I want to maintain multiple conversation sessions, so that I can work on different tasks independently.

#### Acceptance Criteria

1. WHEN the user creates a new session, the system SHALL generate a unique session ID and open a blank conversation
2. WHEN the user switches sessions, the system SHALL load the selected session's history from local storage
3. WHEN the user renames a session, the system SHALL update the session title in the local store
4. WHEN the user archives a session, the system SHALL move it to the archived list without deleting data
5. WHEN the user creates the 21st session, the system SHALL prompt to archive or delete an existing session

---

### Requirement 8: Flutter-Go 通信

**User Story:** AS the Flutter UI, I want to invoke Go runtime methods asynchronously, so that UI remains responsive during long-running Agent operations.

#### Acceptance Criteria

1. WHEN the Flutter UI sends a prompt, the system SHALL call the Go method channel and receive a session ID
2. WHILE the Agent runs, the system SHALL stream incremental events back to Flutter through the method channel
3. WHEN the Go runtime completes a turn, the system SHALL send a completion event with the final result
4. WHEN the Flutter UI requests cancellation, the system SHALL send an abort signal to the Go runtime
5. IF the method channel communication fails, the system SHALL display an error toast and retry once

---

### Requirement 9: 离线可用

**User Story:** AS 用户, I want the app to work fully offline, so that I can use it anywhere without network connectivity.

#### Acceptance Criteria

1. WHEN the app is installed, all core runtime code SHALL be bundled within the APK
2. WHEN the Agent executes tools, the system SHALL NOT make any network requests unless explicitly configured by the user
3. WHEN the app launches, the system SHALL function correctly without any network connection
4. IF the user enables an optional network-dependent tool (e.g., web search), the system SHALL clearly indicate the network requirement

---

### Requirement 10: 应用启动与初始化

**User Story:** AS 用户, I want the app to launch quickly and be ready to use immediately, so that I can start interacting with the Agent without long waits.

#### Acceptance Criteria

1. WHEN the user taps the app icon, the system SHALL display the main UI within 3 seconds
2. WHEN the app first launches, the system SHALL initialize the Go runtime and create the default config
3. WHEN the app resumes from background, the system SHALL restore the previous session state
4. IF the Go runtime initialization fails, the system SHALL display a fatal error screen with a restart option
