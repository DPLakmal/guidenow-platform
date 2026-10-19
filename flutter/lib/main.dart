import 'package:flutter/material.dart';
import 'screens/tourist_home_screen.dart';

void main() {
  runApp(const GuideNowApp());
}

class GuideNowApp extends StatelessWidget {
  const GuideNowApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'GuideNow',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFF10B981),
          primary: const Color(0xFF10B981),
        ),
        useMaterial3: true,
        scaffoldBackgroundColor: const Color(0xFFF8FAFC),
        fontFamily: 'Roboto',
      ),
      home: const TouristHomeScreen(),
    );
  }
}
