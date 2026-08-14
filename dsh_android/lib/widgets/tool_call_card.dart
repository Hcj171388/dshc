import 'package:flutter/material.dart';

enum ToolCallStatus { running, done, error }

class ToolCallCard extends StatelessWidget {
  final String toolName;
  final Map<String, dynamic> args;
  final String? output;
  final String? error;
  final ToolCallStatus status;

  const ToolCallCard({
    super.key,
    required this.toolName,
    this.args = const {},
    this.output,
    this.error,
    this.status = ToolCallStatus.running,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
      elevation: 1,
      child: ExpansionTile(
        leading: _statusIcon,
        title: Text(toolName, style: const TextStyle(fontWeight: FontWeight.w500)),
        subtitle: status == ToolCallStatus.running
            ? const Text('Executing...')
            : null,
        children: [
          if (args.isNotEmpty)
            Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Args:',
                      style: TextStyle(fontWeight: FontWeight.bold, fontSize: 12)),
                  const SizedBox(height: 4),
                  Text(_argsToString(),
                      style: const TextStyle(fontSize: 12, fontFamily: 'monospace')),
                ],
              ),
            ),
          if (status == ToolCallStatus.done && output != null)
            Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Output:',
                      style: TextStyle(fontWeight: FontWeight.bold, fontSize: 12)),
                  const SizedBox(height: 4),
                  Text(output!,
                      style: const TextStyle(fontSize: 12, fontFamily: 'monospace')),
                ],
              ),
            ),
          if (status == ToolCallStatus.error && error != null)
            Padding(
              padding: const EdgeInsets.all(12),
              child: Text('Error: $error',
                  style: const TextStyle(color: Colors.red, fontSize: 12)),
            ),
        ],
      ),
    );
  }

  Widget get _statusIcon {
    switch (status) {
      case ToolCallStatus.running:
        return const SizedBox(
          width: 20,
          height: 20,
          child: CircularProgressIndicator(strokeWidth: 2),
        );
      case ToolCallStatus.done:
        return const Icon(Icons.check_circle, color: Colors.green, size: 20);
      case ToolCallStatus.error:
        return const Icon(Icons.error, color: Colors.red, size: 20);
    }
  }

  String _argsToString() {
    return args.entries
        .map((e) => '${e.key}: ${e.value}')
        .join(', ');
  }
}
