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

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(sessionListProvider);
    return Scaffold(
      appBar: AppBar(
        title: const Text('DeepSeek Harness'),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () => Navigator.pushNamed(context, '/settings'),
          ),
        ],
      ),
      body: state.isLoading
          ? const Center(child: CircularProgressIndicator())
          : state.error != null
              ? Center(child: Text(state.error!))
              : ListView.builder(
                  itemCount: state.sessions.length,
                  itemBuilder: (ctx, i) {
                    final s = state.sessions[i];
                    return ListTile(
                      leading: const Icon(Icons.chat_bubble_outline),
                      title: Text(s.conciseTitle),
                      subtitle: Text(SessionMeta.formatTime(s.updatedAt)),
                      trailing: PopupMenuButton<String>(
                        onSelected: (v) async {
                          if (v == 'delete') {
                            await ref
                                .read(sessionListProvider.notifier)
                                .delete(s.id);
                          } else if (v == 'archive') {
                            await ref
                                .read(sessionListProvider.notifier)
                                .archive(s.id);
                          }
                        },
                        itemBuilder: (_) => const [
                          PopupMenuItem(value: 'delete', child: Text('Delete')),
                          PopupMenuItem(value: 'archive', child: Text('Archive')),
                        ],
                      ),
                      onTap: () => Navigator.push(
                        context,
                        MaterialPageRoute(
                            builder: (_) => ChatView(sessionId: s.id)),
                      ),
                    );
                  },
                ),
      floatingActionButton: FloatingActionButton(
        onPressed: _createAndNavigate,
        child: const Icon(Icons.add),
        tooltip: 'New Session',
      ),
    );
  }
}
