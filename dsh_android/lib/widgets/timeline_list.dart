import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/session_provider.dart';
import '../model/session.dart';

class TimelineList extends ConsumerStatefulWidget {
  final String sessionId;
  const TimelineList({super.key, required this.sessionId});

  @override
  ConsumerState<TimelineList> createState() => _TimelineListState();
}

class _TimelineListState extends ConsumerState<TimelineList> {
  List<EventItem> _events = [];
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadEvents();
  }

  Future<void> _loadEvents() async {
    if (widget.sessionId.isEmpty) return;
    setState(() => _isLoading = true);
    try {
      final bridge = ref.read(bridgeProvider);
      final raw = await bridge.getEvents(widget.sessionId);
      final list = jsonDecode(raw) as List;
      setState(() {
        _events = list.map((e) => EventItem.fromJson(e as Map<String, dynamic>)).toList();
        _isLoading = false;
      });
    } catch (_) {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading && _events.isEmpty) {
      return const SizedBox(height: 60, child: Center(child: CircularProgressIndicator()));
    }
    if (_events.isEmpty) {
      return Container(
        margin: const EdgeInsets.only(bottom: 8),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Row(
          children: [
            Icon(Icons.timeline, color: Colors.grey[400], size: 20),
            const SizedBox(width: 8),
            Text('No execution events yet', style: TextStyle(color: Colors.grey[600])),
          ],
        ),
      );
    }

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.grey[200]!),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: Row(
              children: [
                Icon(Icons.timeline, color: Colors.blue[700], size: 18),
                const SizedBox(width: 8),
                const Text('Execution Timeline', style: TextStyle(fontWeight: FontWeight.w600)),
                const Spacer(),
                Text('${_events.length} events', style: TextStyle(fontSize: 12, color: Colors.grey[600])),
              ],
            ),
          ),
          ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _events.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (ctx, i) {
              final e = _events[i];
              return _TimelineEventTile(event: e);
            },
          ),
        ],
      ),
    );
  }
}

class _TimelineEventTile extends StatelessWidget {
  final EventItem event;
  const _TimelineEventTile({required this.event});

  IconData _getIcon() {
    if (event.isUserMessage) return Icons.person;
    if (event.isAssistantMessage) return Icons.smart_toy;
    if (event.isToolCallStart || event.isToolCallResult) return Icons.build;
    if (event.isTurnStart) return Icons.play_arrow;
    if (event.isTurnComplete) return Icons.check_circle;
    if (event.isError) return Icons.error;
    return Icons.circle;
  }

  Color _getColor() {
    if (event.isUserMessage) return Colors.blue;
    if (event.isAssistantMessage) return Colors.green;
    if (event.isToolCallStart) return Colors.orange;
    if (event.isToolCallResult) return event.toolSuccess ? Colors.green : Colors.red;
    if (event.isError) return Colors.red;
    return Colors.grey;
  }

  String _getTitle() {
    if (event.isUserMessage) return 'You';
    if (event.isAssistantMessage) return 'Assistant';
    if (event.isToolCallStart) return event.toolName;
    if (event.isToolCallResult) return event.toolName;
    if (event.isTurnStart) return 'Turn Started';
    if (event.isTurnComplete) return 'Turn Complete';
    if (event.isError) return 'Error';
    return event.type;
  }

  String _getSubtitle() {
    if (event.isUserMessage) return event.userContent.substring(0, event.userContent.length.clamp(0, 50));
    if (event.isAssistantMessage) return event.assistantContent.substring(0, event.assistantContent.length.clamp(0, 50));
    if (event.isToolCallResult) {
      if (!event.toolSuccess) return 'Failed: ${event.toolError}';
      return event.toolOutput.substring(0, event.toolOutput.length.clamp(0, 50));
    }
    return '';
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      child: Row(
        children: [
          CircleAvatar(
            radius: 16,
            backgroundColor: _getColor().withOpacity(0.1),
            child: Icon(_getIcon(), size: 16, color: _getColor()),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(_getTitle(), style: const TextStyle(fontWeight: FontWeight.w500)),
                if (_getSubtitle().isNotEmpty)
                  Text(
                    _getSubtitle(),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(fontSize: 12, color: Colors.grey[600]),
                  ),
              ],
            ),
          ),
          Text(
            _formatTime(event.timestamp),
            style: TextStyle(fontSize: 11, color: Colors.grey[500]),
          ),
        ],
      ),
    );
  }

  String _formatTime(int epochMs) {
    if (epochMs == 0) return '';
    final dt = DateTime.fromMillisecondsSinceEpoch(epochMs);
    return '${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
  }
}
