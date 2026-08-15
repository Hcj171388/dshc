import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/session_provider.dart';
import '../model/session.dart';
import 'chat_view.dart';

class HomeView extends ConsumerStatefulWidget {
  const HomeView({super.key});

  @override
  ConsumerState<HomeView> createState() => _HomeViewState();
}

class _HomeViewState extends ConsumerState<HomeView> {
  @override
  void initState() {
    super.initState();
    ref.read(sessionListProvider.notifier).load();
  }

  Future<void> _createAndNavigate() async {
    final sid = await ref.read(sessionListProvider.notifier).create();
    if (mounted && sid.isNotEmpty) {
      Navigator.push(
        context,
        MaterialPageRoute(builder: (_) => ChatView(sessionId: sid)),
      );
    }
  }

  Future<void> _renameSession(SessionMeta session) async {
    final result = await showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Rename Session'),
        child: TextField(
          autofocus: true,
          decoration: const InputDecoration(hintText: 'Enter new name'),
          onChanged: (v) => {},
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cancel')),
          FilledButton(onPressed: () => Navigator.pop(ctx), child: const Text('Save')),
        ],
      ),
    );
    if (result != null && result.isNotEmpty) {
      await ref.read(sessionListProvider.notifier).rename(session.id, result);
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(sessionListProvider);
    final sessions = state.sessions;

    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        title: const Text('DeepSeek Harness'),
        backgroundColor: Colors.white,
        elevation: 0,
        leading: null,
        actions: [
          IconButton(
            icon: const Icon(Icons.settings, color: Colors.black87),
            onPressed: () => Navigator.pushNamed(context, '/settings'),
          ),
        ],
      ),
      body: state.isLoading
          ? const Center(child: CircularProgressIndicator())
          : state.error != null
              ? Center(child: Text(state.error!, style: const TextStyle(color: Colors.red)))
              : sessions.isEmpty
                  ? _EmptyState(onCreate: _createAndNavigate)
                  : ListView.builder(
                      padding: const EdgeInsets.only(top: 8),
                      itemCount: sessions.length,
                      itemBuilder: (ctx, i) {
                        final s = sessions[i];
                        return _SessionTile(
                          session: s,
                          onDelete: () async {
                            await ref.read(sessionListProvider.notifier).delete(s.id);
                          },
                          onArchive: () async {
                            await ref.read(sessionListProvider.notifier).archive(s.id);
                          },
                          onRename: () => _renameSession(s),
                          onTap: () {
                            Navigator.push(
                              context,
                              MaterialPageRoute(builder: (_) => ChatView(sessionId: s.id)),
                            );
                          },
                        );
                      },
                    ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: _createAndNavigate,
        icon: const Icon(Icons.add),
        label: const Text('New Session'),
        backgroundColor: Theme.of(context).primaryColor,
      ),
    );
  }
}

class _EmptyState extends StatelessWidget {
  final VoidCallback onCreate;
  const _EmptyState({required this.onCreate});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.chat_bubble_outline, size: 80, color: Colors.grey[300]),
          const SizedBox(height: 16),
          Text('No conversations yet', style: TextStyle(fontSize: 18, color: Colors.grey[600])),
          const SizedBox(height: 8),
          Text('Start a new conversation with the agent', style: TextStyle(color: Colors.grey[500])),
          const SizedBox(height: 24),
          FilledButton.icon(
            onPressed: onCreate,
            icon: const Icon(Icons.add),
            label: const Text('New Conversation'),
          ),
        ],
      ),
    );
  }
}

class _SessionTile extends StatelessWidget {
  final SessionMeta session;
  final VoidCallback onDelete;
  final VoidCallback onArchive;
  final VoidCallback onRename;
  final VoidCallback onTap;

  const _SessionTile({
    required this.session,
    required this.onDelete,
    required this.onArchive,
    required this.onRename,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: CircleAvatar(
        backgroundColor: Colors.blue[50],
        child: Icon(Icons.chat_bubble_outline, color: Colors.blue[700]),
      ),
      title: Text(
        session.conciseTitle,
        style: const TextStyle(fontWeight: FontWeight.w500),
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
      ),
      subtitle: Text(SessionMeta.formatTime(session.updatedAt)),
      trailing: PopupMenuButton<String>(
        icon: Icon(Icons.more_vert, color: Colors.grey[600]),
        onSelected: (v) async {
          if (v == 'delete') await onDelete();
          else if (v == 'archive') await onArchive();
          else if (v == 'rename') await onRename();
        },
        itemBuilder: (_) => [
          const PopupMenuItem(value: 'rename', child: Text('Rename')),
          const PopupMenuItem(value: 'archive', child: Text('Archive')),
          const PopupMenuItem(value: 'delete', child: Text('Delete')),
        ],
      ),
      onTap: onTap,
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
    );
  }
}
