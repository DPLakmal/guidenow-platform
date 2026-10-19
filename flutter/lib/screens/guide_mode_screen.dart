import 'package:flutter/material.dart';

class GuideModeScreen extends StatefulWidget {
  const GuideModeScreen({super.key});

  @override
  State<GuideModeScreen> createState() => _GuideModeScreenState();
}

class _GuideModeScreenState extends State<GuideModeScreen> {
  bool _isAvailable = true;
  int _gpsPings = 18;

  final List<Map<String, dynamic>> _requests = [
    {
      'id': 'tour-mob-101',
      'tourist': 'Sarah Connor',
      'destination': 'Galle Fort Ramparts Heritage Walk',
      'meetingPoint': 'Old Dutch Hospital Entrance',
      'duration': '2h',
      'rate': '\$50',
      'status': 'PENDING',
    },
  ];

  void _accept(int index) {
    setState(() {
      _requests[index]['status'] = 'ACCEPTED';
    });
  }

  void _decline(int index) {
    setState(() {
      _requests[index]['status'] = 'DECLINED';
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        backgroundColor: const Color(0xFF0F172A),
        foregroundColor: Colors.white,
        elevation: 0,
        title: const Text('Guide Dispatch Mode', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Status Card
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: const Color(0xFF0F172A),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text('Current Availability', style: TextStyle(color: Color(0xFF94A3B8), fontSize: 12)),
                          const SizedBox(height: 4),
                          Row(
                            children: [
                              Container(
                                width: 8,
                                height: 8,
                                decoration: BoxDecoration(
                                  color: _isAvailable ? const Color(0xFF10B981) : Colors.redAccent,
                                  shape: BoxShape.circle,
                                ),
                              ),
                              const SizedBox(width: 6),
                              Text(
                                _isAvailable ? 'Online & Available' : 'Offline / Busy',
                                style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 15),
                              ),
                            ],
                          ),
                        ],
                      ),
                      Switch(
                        value: _isAvailable,
                        activeColor: const Color(0xFF10B981),
                        onChanged: (val) {
                          setState(() {
                            _isAvailable = val;
                          });
                        },
                      ),
                    ],
                  ),
                  const Divider(color: Color(0xFF1E293B), height: 24),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Row(
                        children: [
                          const Icon(Icons.satellite_alt_rounded, color: Color(0xFF10B981), size: 18),
                          const SizedBox(width: 8),
                          Text('Redis GPS Telemetry: $_gpsPings pings', style: const TextStyle(color: Color(0xFFE2E8F0), fontSize: 12)),
                        ],
                      ),
                      TextButton.icon(
                        onPressed: () {
                          setState(() {
                            _gpsPings++;
                          });
                        },
                        icon: const Icon(Icons.refresh, size: 14),
                        label: const Text('Ping Now'),
                        style: TextButton.styleFrom(
                          foregroundColor: const Color(0xFF10B981),
                          textStyle: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            const Text(
              'INCOMING TOUR REQUESTS',
              style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: Color(0xFF64748B), letterSpacing: 0.5),
            ),
            const SizedBox(height: 12),

            ..._requests.asMap().entries.map((entry) {
              final idx = entry.key;
              final req = entry.value;

              return Container(
                margin: const EdgeInsets.only(bottom: 12),
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(color: const Color(0xFFE2E8F0)),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(req['tourist'], style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Color(0xFF0F172A))),
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                          decoration: BoxDecoration(
                            color: req['status'] == 'ACCEPTED' ? const Color(0xFFDCFCE7) : const Color(0xFFFEF3C7),
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Text(
                            req['status'],
                            style: TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.bold,
                              color: req['status'] == 'ACCEPTED' ? const Color(0xFF166534) : const Color(0xFF92400E),
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 6),
                    Text(req['destination'], style: const TextStyle(fontSize: 13, color: Color(0xFF334155), fontWeight: FontWeight.w600)),
                    const SizedBox(height: 2),
                    Text('Meeting: ${req['meetingPoint']}', style: const TextStyle(fontSize: 12, color: Color(0xFF64748B))),
                    const SizedBox(height: 4),
                    Text('Duration: ${req['duration']} • Total: ${req['rate']}', style: const TextStyle(fontSize: 12, color: Color(0xFF047857), fontWeight: FontWeight.bold)),
                    if (req['status'] == 'PENDING') ...[
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: ElevatedButton(
                              onPressed: () => _accept(idx),
                              style: ElevatedButton.styleFrom(
                                backgroundColor: const Color(0xFF10B981),
                                foregroundColor: Colors.white,
                                elevation: 0,
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
                              ),
                              child: const Text('Accept Tour', style: TextStyle(fontWeight: FontWeight.bold)),
                            ),
                          ),
                          const SizedBox(width: 8),
                          Expanded(
                            child: OutlinedButton(
                              onPressed: () => _decline(idx),
                              style: OutlinedButton.styleFrom(
                                foregroundColor: Colors.redAccent,
                                side: const BorderSide(color: Color(0xFFCBD5E1)),
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
                              ),
                              child: const Text('Decline'),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ],
                ),
              );
            }),
          ],
        ),
      ),
    );
  }
}
