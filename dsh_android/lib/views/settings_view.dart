import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/session_provider.dart';

class SettingsView extends ConsumerStatefulWidget {
  const SettingsView({super.key});

  @override
  ConsumerState<SettingsView> createState() => _SettingsViewState();
}

class _SettingsViewState extends ConsumerState<SettingsView> {
  late int _maxTurns;
  late int _toolTimeoutMs;
  late bool _showTimeline;

  @override
  void initState() {
    super.initState();
    final cfg = ref.read(configProvider);
    _maxTurns = cfg.maxTurns;
    _toolTimeoutMs = cfg.toolTimeoutMs;
    _showTimeline = cfg.showTimeline;
  }

  void _save() {
    ref.read(configProvider.notifier).state =
        ref.read(configProvider).copyWith(
      maxTurns: _maxTurns,
      toolTimeoutMs: _toolTimeoutMs,
      showTimeline: _showTimeline,
    );
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Settings saved')),
    );
    Navigator.pop(context);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Settings')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          ListTile(
            title: const Text('Max Turns'),
            subtitle: Text('Maximum agent turns per request ($_maxTurns)'),
            trailing: SizedBox(
              width: 120,
              child: Slider(
                value: _maxTurns.toDouble(),
                min: 5,
                max: 200,
                divisions: 39,
                label: _maxTurns.toString(),
                onChanged: (v) => setState(() => _maxTurns = v.toInt()),
              ),
            ),
          ),
          ListTile(
            title: const Text('Tool Timeout (ms)'),
            subtitle: Text('$_toolTimeoutMs ms'),
            trailing: SizedBox(
              width: 120,
              child: Slider(
                value: _toolTimeoutMs.toDouble(),
                min: 5000,
                max: 120000,
                divisions: 23,
                label: '$_toolTimeoutMs',
                onChanged: (v) =>
                    setState(() => _toolTimeoutMs = (v / 5000).round() * 5000),
              ),
            ),
          ),
          SwitchListTile(
            title: const Text('Show Timeline'),
            subtitle: const Text('Display execution timeline panel'),
            value: _showTimeline,
            onChanged: (v) => setState(() => _showTimeline = v),
          ),
          const SizedBox(height: 24),
          FilledButton(
            onPressed: _save,
            child: const Text('Save Settings'),
          ),
        ],
      ),
    );
  }
}
