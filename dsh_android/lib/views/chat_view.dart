import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../model/session.dart';
import '../providers/session_provider.dart';
import '../widgets/message_bubble.dart';
import '../widgets/tool_call_card.dart';
import '../widgets/timeline_list.dart';

class ChatView extends ConsumerStatefulWidget {
  final String sessionId;
  const ChatView({super.key, required this.sessionId});

  @override
  ConsumerState<ChatView> createState() => _ChatViewState();
}

class _ChatViewState extends ConsumerState<ChatView> {
  final TextEditingController _controller = TextEditingController();
  final ScrollController _scrollCtrl = ScrollController();
  StreamSubscription<String>? _eventSub;
  bool _isRunning = false;
  String _lastListenerId = '';
  bool _showTimeline = false;

  @override
  void initState() {
    super.initState();
    _loadInitialEvents();
    _controller.addListener(_scrollToBottom);
  }

  Future<void> _loadInitialEvents() async {
    final snap = ref.read(sessionNotifierProvider(widget.sessionId));
    if (snap.events.isEmpty) {
      await ref.read(sessionNotifierProvider(widget.sessionId).notifier).loadEvents();
    }
  }

  void _scrollToBottom() {
    if (_scrollCtrl.hasClients && !_isRunning) {
      Future.delayed(const Duration(milliseconds: 100), () {
        if (_scrollCtrl.hasClients) {
          _scrollCtrl.animateTo(
            _scrollCtrl.position.maxScrollExtent,
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOut,
          );
        }
      });
    }
  }

  Future<void> _sendPrompt() async {
    final text = _controller.text.trim();
    if (text.isEmpty || _isRunning) return;
    _controller.clear();

    final bridge = ref.read(bridgeProvider);
    final notifier = ref.read(sessionNotifierProvider(widget.sessionId).notifier);

    _isRunning = true;
    setState(() {});
    try {
      _lastListenerId = await bridge.runAgent(widget.sessionId, text);
      _eventSub = bridge.streamEvents(_lastListenerId).listen(
        (eventJson) {
          _handleAgentEvent(eventJson, notifier);
        },
        onDone: () {
          setState(() => _isRunning = false);
          _eventSub?.cancel();
          _eventSub = null;
          notifier.loadEvents().then((_) => _scrollToBottom());
        },
      );
    } catch (e) {
      setState(() => _isRunning = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    }
  }

  void _handleAgentEvent(String eventJson, SessionNotifier notifier) {
    try {
      notifier.loadEvents().then((_) => _scrollToBottom());
    } catch (_) {}
  }

  void _abort() {
    ref.read(bridgeProvider).abortAgent();
    setState(() => _isRunning = false);
    _eventSub?.cancel();
    _eventSub = null;
  }

  @override
  void dispose() {
    _eventSub?.cancel();
    _controller.removeListener(_scrollToBottom);
    _controller.dispose();
    _scrollCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(sessionNotifierProvider(widget.sessionId));
    final config = ref.watch(configProvider);
    _showTimeline = config.showTimeline;

    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        title: Text(state.title.isNotEmpty ? state.title : 'New Conversation'),
        backgroundColor: Colors.white,
        elevation: 1,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.black87),
          onPressed: () => Navigator.pop(context),
        ),
        actions: [
          IconButton(
            icon: Icon(_showTimeline ? Icons.timeline : Icons.timeline_outlined, color: Colors.black87),
            onPressed: () {
              setState(() => _showTimeline = !_showTimeline);
              ref.read(configProvider.notifier).state = config.copyWith(showTimeline: !_showTimeline);
            },
          ),
          if (_isRunning)
            IconButton(
              icon: const Icon(Icons.stop_circle, color: Colors.red),
              onPressed: _abort,
            ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: state.isLoading && state.events.isEmpty
                ? const Center(child: CircularProgressIndicator())
                : state.error != null
                    ? Center(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(Icons.error_outline, size: 48, color: Colors.red),
                            const SizedBox(height: 16),
                            Text(state.error!, style: const TextStyle(color: Colors.red)),
                            const SizedBox(height: 16),
                            FilledButton(onPressed: () => ref.read(sessionNotifierProvider(widget.sessionId).notifier).loadEvents(), child: const Text('Retry')),
                          ],
                        ),
                      )
                    : ListView(
                        controller: _scrollCtrl,
                        padding: const EdgeInsets.symmetric(vertical: 8),
                        children: [
                          if (_showTimeline)
                            const Padding(
                              padding: EdgeInsets.symmetric(horizontal: 12),
                              child: TimelineList(sessionId: ''),
                            ),
                          ...state.events.map((e) => _buildEvent(e)),
                          if (_isRunning)
                            const Padding(
                              padding: EdgeInsets.all(16),
                              child: Row(
                                children: [
                                  SizedBox(
                                    width: 20,
                                    height: 20,
                                    child: CircularProgressIndicator(strokeWidth: 2),
                                  ),
                                  SizedBox(width: 12),
                                  Text('Agent is thinking...', style: TextStyle(color: Colors.grey)),
                                ],
                              ),
                            ),
                        ],
                      ),
          ),
          _BuildInputBar(
            controller: _controller,
            isRunning: _isRunning,
            onSend: _sendPrompt,
            onAbort: _abort,
          ),
        ],
      ),
    );
  }

  Widget _buildEvent(EventItem e) {
    if (e.isUserMessage) {
      return MessageBubble(message: e.userContent, isUser: true);
    } else if (e.isAssistantMessage) {
      return MessageBubble(message: e.assistantContent, isUser: false);
    } else if (e.isToolCallStart) {
      return ToolCallCard(
        toolName: e.toolName,
        args: e.toolArgs,
        status: ToolCallStatus.running,
      );
    } else if (e.isToolCallResult) {
      return ToolCallCard(
        toolName: e.toolName,
        output: e.toolOutput,
        error: e.toolError,
        status: e.toolSuccess ? ToolCallStatus.done : ToolCallStatus.error,
      );
    }
    return const SizedBox.shrink();
  }
}

class _BuildInputBar extends StatelessWidget {
  final TextEditingController controller;
  final bool isRunning;
  final VoidCallback onSend;
  final VoidCallback onAbort;

  const _BuildInputBar({
    required this.controller,
    required this.isRunning,
    required this.onSend,
    required this.onAbort,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 4, offset: const Offset(0, -2)),
        ],
      ),
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
      child: SafeArea(
        child: Row(
          children: [
            Expanded(
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.grey[100],
                  borderRadius: BorderRadius.circular(24),
                ),
                child: TextField(
                  controller: controller,
                  decoration: const InputDecoration(
                    hintText: 'Send a message...',
                    border: InputBorder.none,
                    contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                  ),
                  maxLines: 4,
                  minLines: 1,
                  textInputAction: TextInputAction.send,
                  onSubmitted: (_) => onSend(),
                ),
              ),
            ),
            const SizedBox(width: 8),
            isRunning
                ? CircleAvatar(
                    radius: 22,
                    backgroundColor: Colors.red,
                    child: IconButton(
                      icon: const Icon(Icons.stop, color: Colors.white, size: 20),
                      onPressed: onAbort,
                    ),
                  )
                : CircleAvatar(
                    radius: 22,
                    backgroundColor: Theme.of(context).primaryColor,
                    child: IconButton(
                      icon: const Icon(Icons.send, color: Colors.white, size: 20),
                      onPressed: onSend,
                    ),
                  ),
          ],
        ),
      ),
    );
  }
}
