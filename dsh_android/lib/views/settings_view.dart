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
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        title: const Text('Settings'),
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.black87),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          _SettingsSection(
            title: 'Agent Behavior',
            children: [
              _SliderSetting(
                title: 'Max Turns',
                subtitle: 'Maximum agent turns per request ($_maxTurns)',
                value: _maxTurns,
                min: 5,
                max: 200,
                divisions: 39,
                onChanged: (v) => setState(() => _maxTurns = v.toInt()),
              ),
              const SizedBox(height: 16),
              _SliderSetting(
                title: 'Tool Timeout',
                subtitle: '${_toolTimeoutMs / 1000}s timeout for tool execution',
                value: _toolTimeoutMs ~/ 1000,
                min: 5,
                max: 120,
                divisions: 23,
                labelSuffix: 's',
                onChanged: (v) => setState(() => _toolTimeoutMs = v * 1000),
              ),
            ],
          ),
          const SizedBox(height: 16),
          _SettingsSection(
            title: 'Display',
            children: [
              SwitchListTile(
                title: const Text('Show Timeline'),
                subtitle: const Text('Display execution timeline panel in chat'),
                value: _showTimeline,
                onChanged: (v) => setState(() => _showTimeline = v),
                contentPadding: EdgeInsets.zero,
              ),
            ],
          ),
          const SizedBox(height: 24),
          FilledButton(
            onPressed: _save,
            style: FilledButton.styleFrom(
              padding: const EdgeInsets.symmetric(vertical: 16),
              minimumSize: const Size.fromHeight(56),
            ),
            child: const Text('Save Settings', style: TextStyle(fontSize: 16)),
          ),
          const SizedBox(height: 32),
          _AboutSection(),
        ],
      ),
    );
  }
}

class _SettingsSection extends StatelessWidget {
  final String title;
  final List<Widget> children;
  const _SettingsSection({required this.title, required this.children});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 4, offset: const Offset(0, 2)),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
            child: Text(title, style: TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: Colors.grey[600])),
          ),
          ...children,
        ],
      ),
    );
  }
}

class _SliderSetting extends StatelessWidget {
  final String title;
  final String subtitle;
  final int value;
  final int min;
  final int max;
  final int divisions;
  final String labelSuffix;
  final ValueChanged<double> onChanged;

  const _SliderSetting({
    required this.title,
    required this.subtitle,
    required this.value,
    required this.min,
    required this.max,
    required this.divisions,
    this.labelSuffix = '',
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Text(title, style: const TextStyle(fontWeight: FontWeight.w500)),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          child: Text(subtitle, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Slider(
            value: value.toDouble(),
            min: min.toDouble(),
            max: max.toDouble(),
            divisions: divisions,
            label: '$value$labelSuffix',
            onChanged: onChanged,
          ),
        ),
      ],
    );
  }
}

class _AboutSection extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        children: [
          Text('DeepSeek Harness Android', style: TextStyle(color: Colors.grey[600])),
          const SizedBox(height: 4),
          Text('Version 1.2.0', style: TextStyle(fontSize: 12, color: Colors.grey[500])),
        ],
      ),
    );
  }
}
