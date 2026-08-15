import 'package:flutter/material.dart';

enum ToolCallStatus { running, done, error }

class ToolCallCard extends StatefulWidget {
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
  State<ToolCallCard> createState() => _ToolCallCardState();
}

class _ToolCallCardState extends State<ToolCallCard> {
  bool _expanded = false;

  IconData _getIcon() {
    switch (widget.toolName.toLowerCase()) {
      case 'bash':
      case 'terminal':
        return Icons.terminal;
      case 'read_file':
      case 'write_file':
      case 'list_dir':
      case 'search_files':
        return Icons.insert_drive_file;
      case 'web_search':
      case 'web_fetch':
        return Icons.web;
      case 'git_command':
        return Icons.git_merge;
      case 'curl':
        return Icons.http;
      case 'code_runtime':
        return Icons.code;
      case 'todo_write':
        return Icons.checklist;
      case 'image_process':
        return Icons.image;
      case 'video_process':
        return Icons.videocam;
      default:
        return Icons.build;
    }
  }

  Color _getStatusColor() {
    switch (widget.status) {
      case ToolCallStatus.running:
        return Colors.orange;
      case ToolCallStatus.done:
        return Colors.green;
      case ToolCallStatus.error:
        return Colors.red;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.grey[200]!),
        boxShadow: [
          BoxShadow(color: Colors.black.withOpacity(0.03), blurRadius: 4, offset: const Offset(0, 2)),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          InkWell(
            onTap: () => setState(() => _expanded = !_expanded),
            borderRadius: BorderRadius.circular(12),
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Row(
                children: [
                  _StatusDot(color: _getStatusColor(), isRunning: widget.status == ToolCallStatus.running),
                  const SizedBox(width: 12),
                  Icon(_getIcon(), size: 20, color: _getStatusColor()),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Text(
                      widget.toolName,
                      style: const TextStyle(fontWeight: FontWeight.w500),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  if (widget.status == ToolCallStatus.running)
                    const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  else if (widget.status == ToolCallStatus.done)
                    const Icon(Icons.check, color: Colors.green, size: 18)
                  else
                    const Icon(Icons.error, color: Colors.red, size: 18),
                  const SizedBox(width: 8),
                  Icon(
                    _expanded ? Icons.expand_less : Icons.expand_more,
                    color: Colors.grey[600],
                    size: 20,
                  ),
                ],
              ),
            ),
          ),
          if (_expanded)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (widget.args.isNotEmpty) ...[
                    _InfoSection(title: 'Arguments', content: _argsToString()),
                    const SizedBox(height: 8),
                  ],
                  if (widget.status == ToolCallStatus.done && widget.output != null)
                    _InfoSection(title: 'Output', content: widget.output!),
                  if (widget.status == ToolCallStatus.error && widget.error != null)
                    _InfoSection(title: 'Error', content: widget.error!, isError: true),
                  const SizedBox(height: 8),
                ],
              ),
            ),
        ],
      ),
    );
  }

  String _argsToString() {
    return widget.args.entries
        .map((e) => '${e.key}: ${e.value}')
        .join(', ');
  }
}

class _StatusDot extends StatelessWidget {
  final Color color;
  final bool isRunning;
  const _StatusDot({required this.color, this.isRunning = false});

  @override
  Widget build(BuildContext context) {
    return AnimatedContainer(
      duration: const Duration(milliseconds: 300),
      width: 10,
      height: 10,
      decoration: BoxDecoration(
        color: color,
        shape: BoxShape.circle,
        boxShadow: isRunning
            ? [BoxShadow(color: color, blurRadius: 6, spreadPercentage: 0.3)]
            : null,
      ),
    );
  }
}

class _InfoSection extends StatelessWidget {
  final String title;
  final String content;
  final bool isError;
  const _InfoSection({required this.title, required this.content, this.isError = false});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: TextStyle(
            fontSize: 11,
            fontWeight: FontWeight.w600,
            color: isError ? Colors.red[700] : Colors.grey[600],
            letterSpacing: 0.5,
          ),
        ),
        const SizedBox(height: 4),
        Container(
          width: double.infinity,
          padding: const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color: isError ? Colors.red[50] : Colors.grey[100],
            borderRadius: BorderRadius.circular(8),
          ),
          child: SelectableText(
            content,
            style: TextStyle(
              fontSize: 12,
              fontFamily: 'monospace',
              color: isError ? Colors.red[900] : Colors.grey[900],
            ),
          ),
        ),
      ],
    );
  }
}
