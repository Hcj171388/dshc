import 'package:flutter/material.dart';

class TimelinePanel extends StatelessWidget {
  final String sessionId;
  final bool isLive;

  const TimelinePanel({super.key, required this.sessionId, this.isLive = false});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.all(8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.grey[100],
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.grey[300]!),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.timeline, size: 18),
              const SizedBox(width: 8),
              Text(isLive ? 'Live Execution' : 'Execution History',
                  style: const TextStyle(fontWeight: FontWeight.bold)),
              if (isLive)
                Padding(
                  padding: const EdgeInsets.only(left: 8),
                  child: SizedBox(
                    width: 8,
                    height: 8,
                    child: RawMaterialButton(
                      onPressed: null,
                      fillColor: Colors.green,
                      shape: const CircleBorder(),
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 8),
          const Text('No active execution',
              style: TextStyle(color: Colors.grey, fontSize: 13)),
        ],
      ),
    );
  }
}
