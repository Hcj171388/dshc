// Data models for session events and metadata.

import 'dart:convert';

class SessionMeta {
  final String id;
  final String title;
  final int createdAt;
  final int updatedAt;

  const SessionMeta({
    required this.id,
    required this.title,
    required this.createdAt,
    required this.updatedAt,
  });

  factory SessionMeta.fromJson(Map<String, dynamic> json) {
    return SessionMeta(
      id: json['id'] as String,
      title: (json['title'] as String?) ?? '',
      createdAt: (json['created_at'] as num?)?.toInt() ?? 0,
      updatedAt: (json['updated_at'] as num?)?.toInt() ?? 0,
    );
  }

  String get conciseTitle => title.isNotEmpty ? title : formatTime(createdAt);
  static String formatTime(int epochMs) {
    final dt = DateTime.fromMillisecondsSinceEpoch(epochMs);
    final now = DateTime.now();
    if (dt.year == now.year && dt.month == now.month && dt.day == now.day) {
      return '${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
    }
    return '${dt.month}/${dt.day} ${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
  }
}

class EventItem {
  final int id;
  final String type;
  final Map<String, dynamic> payload;
  final int timestamp;

  const EventItem({
    required this.id,
    required this.type,
    required this.payload,
    required this.timestamp,
  });

  factory EventItem.fromJson(Map<String, dynamic> json) {
    dynamic payloadRaw = json['payload'];
    Map<String, dynamic> payload;
    if (payloadRaw is String) {
      try {
        payload = jsonDecode(payloadRaw) as Map<String, dynamic>;
      } catch (_) {
        payload = {'raw': payloadRaw};
      }
    } else if (payloadRaw is Map) {
      payload = payloadRaw as Map<String, dynamic>;
    } else {
      payload = {};
    }
    return EventItem(
      id: (json['id'] as num?)?.toInt() ?? 0,
      type: json['type'] as String,
      payload: payload,
      timestamp: (json['timestamp'] as num?)?.toInt() ?? 0,
    );
  }

  bool get isUserMessage => type == 'user_message';
  bool get isAssistantMessage => type == 'assistant_message';
  bool get isToolCallStart => type == 'tool_call_start';
  bool get isToolCallResult => type == 'tool_call_result';
  bool get isTurnStart => type == 'turn_start';
  bool get isTurnComplete => type == 'turn_complete';
  bool get isError => type == 'error';

  String get userContent => payload['content'] as String? ?? '';
  String get assistantContent => payload['content'] as String? ?? '';
  String get toolName => payload['tool_name'] as String? ?? '';
  Map<String, dynamic> get toolArgs =>
      (payload['args'] as Map?)?.cast<String, dynamic>() ?? {};
  bool get toolSuccess => payload['success'] as bool? ?? true;
  String get toolOutput => payload['output'] as String? ?? '';
  String get toolError => payload['error'] as String? ?? '';
  String get prompt => payload['prompt'] as String? ?? '';
  String get response => payload['response'] as String? ?? '';
  int get toolCallCount => (payload['tool_calls'] as num?)?.toInt() ?? 0;
  int get durationMs => (payload['duration_ms'] as num?)?.toInt() ?? 0;
  String get errorMessage => payload['message'] as String? ?? '';
}
