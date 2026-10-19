import 'package:flutter_test/flutter_test.dart';
import 'package:guidenow/main.dart';

void main() {
  testWidgets('GuideNow app smoke test and title rendering', (WidgetTester tester) async {
    await tester.pumpWidget(const GuideNowApp());
    await tester.pump();

    // Verify GuideNow title renders
    expect(find.text('GuideNow'), findsOneWidget);
    expect(find.text('Sri Lanka Realtime Tour Matching'), findsOneWidget);
  });
}
