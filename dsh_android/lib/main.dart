import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'views/home_view.dart';
import 'views/settings_view.dart';

void main() {
  runApp(const ProviderScope(child: DeepSeekHarness()));
}

class DeepSeekHarness extends StatelessWidget {
  const DeepSeekHarness({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'DeepSeek Harness',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
        useMaterial3: true,
      ),
      initialRoute: '/',
      routes: {
        '/': (context) => const HomeView(),
        '/settings': (context) => const SettingsView(),
      },
    );
  }
}
