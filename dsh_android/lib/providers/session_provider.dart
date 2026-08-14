import 'dart:convert';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:uuid/uuid.dart';

import '../bridge/go_bridge.dart';
import '../model/session.dart';

final bridgeProvider = Provider<GoBridge>((ref) => GoBridge());

final uuidProvider = Provider((ref) => const Uuid());

// ---- Session state ----

class SessionState {
  final String id;
  final String title;
  final int createdAt;
  final int updatedAt;
  final List<EventItem> events;
  final bool isLoading;
  final String? error;

  const SessionState({
    required this.id,
    required this.title,
    required this.createdAt,
    required this.updatedAt,
    this.events = const [],
    this.isLoading = false,
    this.error,
  });

  SessionState copyWith({
    String? title,
    List<EventItem>? events,
    bool? isLoading,
    String? error,
    int? createdAt,
    int? updatedAt,
  }) {
    return SessionState(
      id: id,
      title: title ?? this.title,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      events: events ?? this.events,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
    );
  }
}

class SessionNotifier extends StateNotifier<SessionState> {
  final GoBridge bridge;
  final String sessionId;

  SessionNotifier({required this.bridge, required this.sessionId})
      : super(SessionState(id: sessionId, title: '', createdAt: 0, updatedAt: 0));

  Future<void> loadEvents() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final meta = await bridge.getSession(sessionId);
      final parsed = jsonDecode(meta) as Map<String, dynamic>;
      final title = parsed['title'] as String? ?? '';
      final createdAt = (parsed['created_at'] as num?)?.toInt() ?? 0;
      final updatedAt = (parsed['updated_at'] as num?)?.toInt() ?? 0;

      final eventsRaw = await bridge.getEvents(sessionId);
      final eventList = (jsonDecode(eventsRaw) as List)
          .map((e) => EventItem.fromJson(e as Map<String, dynamic>))
          .toList();

      state = state.copyWith(
        title: title,
        events: eventList,
        isLoading: false,
        createdAt: createdAt,
        updatedAt: updatedAt,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }
}

final sessionNotifierProvider =
    StateNotifierProvider.family<SessionNotifier, SessionState, String>((ref, id) {
  return SessionNotifier(bridge: ref.watch(bridgeProvider), sessionId: id);
});

// ---- Session list ----

class SessionListState {
  final List<SessionMeta> sessions;
  final bool isLoading;
  final String? error;

  const SessionListState({
    this.sessions = const [],
    this.isLoading = false,
    this.error,
  });

  SessionListState copyWith({
    List<SessionMeta>? sessions,
    bool? isLoading,
    String? error,
  }) {
    return SessionListState(
      sessions: sessions ?? this.sessions,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
    );
  }
}

class SessionListNotifier extends StateNotifier<SessionListState> {
  final GoBridge bridge;
  SessionListNotifier({required this.bridge})
      : super(const SessionListState());

  Future<void> load() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final raw = await bridge.listSessions();
      final list = jsonDecode(raw) as List;
      final sessions = list
          .map((e) => SessionMeta.fromJson(e as Map<String, dynamic>))
          .toList();
      state = state.copyWith(sessions: sessions, isLoading: false);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<String> create() async {
    final id = await bridge.createSession();
    await load();
    return id;
  }

  Future<void> delete(String id) async {
    await bridge.deleteSession(id);
    await load();
  }

  Future<void> archive(String id) async {
    await bridge.archiveSession(id);
    await load();
  }

  Future<void> rename(String id, String title) async {
    await bridge.updateSessionTitle(id, title);
    await load();
  }
}

final sessionListProvider =
    StateNotifierProvider<SessionListNotifier, SessionListState>((ref) {
  return SessionListNotifier(bridge: ref.watch(bridgeProvider));
});

// ---- Config ----

class ConfigState {
  final int maxTurns;
  final int toolTimeoutMs;
  final int maxParallel;
  final bool showTimeline;
  final String theme;

  const ConfigState({
    this.maxTurns = 50,
    this.toolTimeoutMs = 30000,
    this.maxParallel = 3,
    this.showTimeline = true,
    this.theme = 'system',
  });

  ConfigState copyWith({
    int? maxTurns,
    int? toolTimeoutMs,
    int? maxParallel,
    bool? showTimeline,
    String? theme,
  }) {
    return ConfigState(
      maxTurns: maxTurns ?? this.maxTurns,
      toolTimeoutMs: toolTimeoutMs ?? this.toolTimeoutMs,
      maxParallel: maxParallel ?? this.maxParallel,
      showTimeline: showTimeline ?? this.showTimeline,
      theme: theme ?? this.theme,
    );
  }
}

final configProvider = StateProvider<ConfigState>((ref) {
  return const ConfigState();
});
