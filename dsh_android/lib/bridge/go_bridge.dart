import 'dart:async';
import 'dart:convert';

import 'package:flutter/services.dart';

/// Bridges Flutter UI to the Go native runtime via MethodChannel/EventChannel.
class GoBridge {
  static const MethodChannel _method =
      MethodChannel('com.deepseek.harness/agent');
  static const MethodChannel _config =
      MethodChannel('com.deepseek.harness/config');

  // Per-session event channels
  final Map<String, StreamSubscription<dynamic>> _listeners = {};

  Future<T> _call<T>(String method, [dynamic args]) async {
    try {
      final result = await _method.invokeMethod(method, args);
      return result as T;
    } on PlatformException catch (e) {
      throw HarnessError(e.code, e.message ?? 'unknown error');
    }
  }

  /// Creates a new session and returns its ID.
  Future<String> createSession() async => _call<String>('createSession');

  /// Returns all sessions as a JSON string (list of maps).
  Future<String> listSessions() async => _call<String>('listSessions');

  /// Returns one session's metadata as JSON string.
  Future<String> getSession(String sessionId) async =>
      _call<String>('getSession', sessionId);

  /// Deletes a session by ID.
  Future<void> deleteSession(String sessionId) async =>
      _call<void>('deleteSession', sessionId);

  /// Archives a session by ID.
  Future<void> archiveSession(String sessionId) async =>
      _call<void>('archiveSession', sessionId);

  /// Updates a session title.
  Future<void> updateSessionTitle(String sessionId, String title) async =>
      _call<void>('updateSessionTitle',
          {'id': sessionId, 'title': title});

  /// Returns recent events for a session, starting after lastEventId.
  Future<String> getEvents(String sessionId,
      {int afterId = 0, int limit = 100}) async {
    return _call<String>('getEvents', {
      'sessionId': sessionId,
      'afterId': afterId,
      'limit': limit,
    });
  }

  /// Starts an agent turn with the given prompt. Returns a listener ID.
  Future<String> runAgent(String sessionId, String prompt) async =>
      _call<String>('runAgent', {'sessionId': sessionId, 'prompt': prompt});

  /// Stream events for a previously-started agent turn.
  /// Passes each event JSON string through the callback channel.
  Stream<String> streamEvents(String listenerId) {
    const eventChannel = EventChannel('com.deepseek.harness/events');
    return eventChannel
        .receiveBroadcastStream({'listenerId': listenerId})
        .map((e) => e as String);
  }

  /// Aborts the currently running agent turn.
  Future<void> abortAgent() async => _call<void>('abortAgent');

  /// Loads the current app configuration as a JSON string.
  Future<String> getConfig() async => _call<String>('getConfig');

  /// Saves the full app configuration (JSON string).
  Future<void> saveConfig(String configJson) async =>
      _call<void>('saveConfig', configJson);

  /// Disposes all pending event subscriptions.
  void dispose() {
    for (final sub in _listeners.values) {
      sub.cancel();
    }
    _listeners.clear();
  }
}

class HarnessError implements Exception {
  final String code;
  final String message;
  HarnessError(this.code, this.message);
  @override
  String toString() => 'HarnessError($code): $message';
}
