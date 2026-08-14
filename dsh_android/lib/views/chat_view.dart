import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../model/session.dart';
import '../providers/session_provider.dart';
import '../widgets/message_bubble.dart';
import '../widgets/tool_call_card.dart';
import 'timeline_panel.dart';

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

  @override
  void initState() {
    super.initState();
    _loadInitialEvents();
  }

  Future<void> _loadInitialEvents() async {
    final snap = ref.read(sessionNotifierProvider(widget.sessionId));
    if (snap.events.isEmpty) {
      await ref.read(sessionNotifierProvider(widget.sessionId).notifier).loadEvents();
    }
  }

  Future<void> _sendPrompt() async {
    final text = _controller.text.trim();
    if (text.isEmpty || _isRunning) return;
    _controller.clear();

    // Append user message locally
    final bridge = ref.read(bridgeProvider);
    final notifier = ref.read(sessionNotifierProvider(widget.sessionId).notifier);

    // Start agent turn
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
          notifier.loadEvents();
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
      final ev = jsonDecode(eventJson) as Map<String, dynamic>;
      // Store event in session
      // For now just reload
      notifier.loadEvents();
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
    _controller.dispose();
    _scrollCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(sessionNotifierProvider(widget.sessionId));
    final showTimeline = ref.watch(configProvider).showTimeline;

    return Scaffold(
      appBar: AppBar(
        title: Text(state.title.isNotEmpty ? state.title : 'Chat'),
        actions: [
          if (_isRunning)
            IconButton(
              icon: const Icon(Icons.stop),
              onPressed: _abort,
              tooltip: 'Abort',
            ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: state.isLoading && state.events.isEmpty
                ? const Center(child: CircularProgressIndicator())
                : state.error != null
                    ? Center(child: Text(state.error!))
                    : ListView(
                        controller: _scrollCtrl,
                        children: [
                          if (showTimeline)
                            const TimelinePanel(sessionId: '', isLive: true),
                          ...state.events.map((e) => _buildEvent(e)),
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
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        child: Row(
          children: [
            Expanded(
              child: TextField(
                controller: controller,
                decoration: const InputDecoration(
                  hintText: 'Type a message...',
                  border: OutlineInputBorder(),
                ),
                maxLines: 3,
                minLines: 1,
                textInputAction: TextInputAction.send,
                onSubmitted: (_) => onSend(),
              ),
            ),
            const SizedBox(width: 8),
            isRunning
                ? IconButton(
                    icon: const Icon(Icons.stop, color: Colors.red),
                    onPressed: onAbort,
                  )
                : IconButton(
                    icon: const Icon(Icons.send),
                    onPressed: onSend,
                  ),
          ],
        ),
      ),
    );
  }
}
