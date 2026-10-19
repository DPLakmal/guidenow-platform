import 'dart:math';
import 'package:flutter/material.dart';
import '../models/guide.dart';

class RadarMapWidget extends StatefulWidget {
  final List<Guide> guides;
  final Guide? selectedGuide;
  final Function(Guide) onSelectGuide;
  final String locationName;

  const RadarMapWidget({
    super.key,
    required this.guides,
    required this.selectedGuide,
    required this.onSelectGuide,
    required this.locationName,
  });

  @override
  State<RadarMapWidget> createState() => _RadarMapWidgetState();
}

class _RadarMapWidgetState extends State<RadarMapWidget> with SingleTickerProviderStateMixin {
  late AnimationController _sweepController;

  @override
  void initState() {
    super.initState();
    _sweepController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 4),
    )..repeat();
  }

  @override
  void dispose() {
    _sweepController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 240,
      width: double.infinity,
      decoration: BoxDecoration(
        color: const Color(0xFF0F172A),
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: const Color(0xFF1E293B)),
        boxShadow: const [
          BoxShadow(
            color: Colors.black26,
            blurRadius: 10,
            offset: Offset(0, 4),
          ),
        ],
      ),
      child: Stack(
        alignment: Alignment.center,
        children: [
          // Animated radar sweep
          AnimatedBuilder(
            animation: _sweepController,
            builder: (context, child) {
              return CustomPaint(
                size: const Size(200, 200),
                painter: _RadarGridPainter(angle: _sweepController.value * 2 * pi),
              );
            },
          ),

          // User Pin in Center
          Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 14,
                height: 14,
                decoration: BoxDecoration(
                  color: const Color(0xFF10B981),
                  shape: BoxShape.circle,
                  border: Border.all(color: Colors.white, width: 2),
                  boxShadow: const [
                    BoxShadow(color: Color(0x6610B981), blurRadius: 8, spreadRadius: 2),
                  ],
                ),
              ),
              const SizedBox(height: 4),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                decoration: BoxDecoration(
                  color: const Color(0xCC0F172A),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: const Color(0x4D10B981)),
                ),
                child: Text(
                  'You (${widget.locationName})',
                  style: const TextStyle(color: Color(0xFF10B981), fontSize: 10, fontWeight: FontWeight.bold),
                ),
              ),
            ],
          ),

          // Guide Blips
          ...widget.guides.asMap().entries.map((entry) {
            final idx = entry.key;
            final guide = entry.value;
            final isSelected = widget.selectedGuide?.id == guide.id;

            final angle = (idx * (2 * pi / max(widget.guides.length, 1))) + 0.6;
            final radius = 50.0 + (idx * 20.0).clamp(0.0, 40.0);
            final dx = radius * cos(angle);
            final dy = radius * sin(angle);

            return Transform.translate(
              offset: Offset(dx, dy),
              child: GestureDetector(
                onTap: () => widget.onSelectGuide(guide),
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: isSelected ? const Color(0xFFF59E0B) : const Color(0xFF1E293B),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: isSelected ? Colors.amberAccent : const Color(0xFF334155),
                      width: isSelected ? 2 : 1,
                    ),
                    boxShadow: const [
                      BoxShadow(color: Colors.black38, blurRadius: 4, offset: Offset(0, 2)),
                    ],
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      CircleAvatar(
                        radius: 8,
                        backgroundColor: Colors.white24,
                        child: Text(
                          guide.name[0],
                          style: TextStyle(
                            fontSize: 9,
                            fontWeight: FontWeight.bold,
                            color: isSelected ? Colors.black87 : Colors.white,
                          ),
                        ),
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '${guide.name.split(' ')[0]} \$${guide.hourlyRate.toInt()}',
                        style: TextStyle(
                          fontSize: 10,
                          fontWeight: FontWeight.bold,
                          color: isSelected ? Colors.black87 : Colors.white,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}

class _RadarGridPainter extends CustomPainter {
  final double angle;

  _RadarGridPainter({required this.angle});

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final paint = Paint()
      ..color = const Color(0x2610B981)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.0;

    // Rings
    canvas.drawCircle(center, size.width * 0.45, paint);
    canvas.drawCircle(center, size.width * 0.30, paint);
    canvas.drawCircle(center, size.width * 0.15, paint);

    // Crosshairs
    canvas.drawLine(Offset(0, center.dy), Offset(size.width, center.dy), paint);
    canvas.drawLine(Offset(center.dx, 0), Offset(center.dx, size.height), paint);

    // Sweep line
    final sweepPaint = Paint()
      ..color = const Color(0x6610B981)
      ..strokeWidth = 1.5;
    final endPoint = Offset(
      center.dx + (size.width * 0.45) * cos(angle),
      center.dy + (size.width * 0.45) * sin(angle),
    );
    canvas.drawLine(center, endPoint, sweepPaint);
  }

  @override
  bool shouldRepaint(covariant _RadarGridPainter oldDelegate) => oldDelegate.angle != angle;
}
