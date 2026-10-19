import 'package:flutter/material.dart';
import '../models/guide.dart';
import '../services/api_service.dart';
import '../widgets/radar_map_widget.dart';
import '../widgets/guide_tile.dart';
import 'live_chat_screen.dart';
import 'guide_mode_screen.dart';

class TouristHomeScreen extends StatefulWidget {
  const TouristHomeScreen({super.key});

  @override
  State<TouristHomeScreen> createState() => _TouristHomeScreenState();
}

class _TouristHomeScreenState extends State<TouristHomeScreen> {
  final List<Map<String, dynamic>> _hubs = [
    {'name': 'Galle Fort', 'lat': 6.0329, 'lng': 80.2168},
    {'name': 'Colombo', 'lat': 6.9271, 'lng': 79.8612},
    {'name': 'Weligama', 'lat': 5.9723, 'lng': 80.4287},
    {'name': 'Kandy', 'lat': 7.2906, 'lng': 80.6337},
    {'name': 'Sigiriya', 'lat': 7.9570, 'lng': 80.7603},
    {'name': 'Ella', 'lat': 6.8667, 'lng': 81.0466},
  ];

  late Map<String, dynamic> _selectedHub;
  List<Guide> _guides = [];
  Guide? _selectedGuide;
  bool _loading = false;

  @override
  void initState() {
    super.initState();
    _selectedHub = _hubs[0];
    _loadGuides();
  }

  Future<void> _loadGuides() async {
    setState(() => _loading = true);
    final data = await ApiService.fetchNearbyGuides(
      lat: _selectedHub['lat'],
      lng: _selectedHub['lng'],
    );
    setState(() {
      _guides = data;
      if (data.isNotEmpty) {
        _selectedGuide = data[0];
      }
      _loading = false;
    });
  }

  void _requestTour(Guide guide) {
    showDialog(
      context: context,
      builder: (context) {
        double hours = 2.0;
        return StatefulBuilder(
          builder: (context, setDialogState) {
            return AlertDialog(
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
              title: Text('Request Tour with ${guide.name}', style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
              content: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('City: ${guide.city} • Rate: \$${guide.hourlyRate.toInt()}/hr', style: const TextStyle(fontSize: 12, color: Color(0xFF64748B))),
                  const SizedBox(height: 16),
                  const Text('Select Duration:', style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [1, 2, 4, 8].map((h) {
                      final isSel = hours == h.toDouble();
                      return ChoiceChip(
                        label: Text('${h}h'),
                        selected: isSel,
                        selectedColor: const Color(0xFF10B981),
                        labelStyle: TextStyle(color: isSel ? Colors.white : Colors.black87, fontSize: 12),
                        onSelected: (selected) {
                          if (selected) setDialogState(() => hours = h.toDouble());
                        },
                      );
                    }).toList(),
                  ),
                  const SizedBox(height: 16),
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: const Color(0xFFECFDF5),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: const Color(0xFFA7F3D0)),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text('Estimated Total:', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Color(0xFF065F46))),
                        Text('\$${(hours * guide.hourlyRate).toInt()}', style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800, color: Color(0xFF047857))),
                      ],
                    ),
                  ),
                ],
              ),
              actions: [
                TextButton(
                  onPressed: () => Navigator.pop(context),
                  child: const Text('Cancel'),
                ),
                ElevatedButton(
                  onPressed: () async {
                    Navigator.pop(context);
                    await ApiService.createTourRequest(
                      guideId: guide.id,
                      touristName: 'Alex Rivera',
                      meetingPoint: '${guide.city} Dutch Hospital',
                      destination: '${guide.specialties.first} Highlights',
                      durationHours: hours,
                    );
                    if (mounted) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          content: Text('Tour request sent to ${guide.name}!'),
                          backgroundColor: const Color(0xFF10B981),
                        ),
                      );
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF10B981),
                    foregroundColor: Colors.white,
                  ),
                  child: const Text('Dispatch Request'),
                ),
              ],
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        title: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(6),
              decoration: BoxDecoration(
                color: const Color(0xFF10B981),
                borderRadius: BorderRadius.circular(10),
              ),
              child: const Icon(Icons.explore, color: Colors.white, size: 18),
            ),
            const SizedBox(width: 8),
            const Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('GuideNow', style: TextStyle(fontWeight: FontWeight.w900, fontSize: 16, color: Color(0xFF0F172A))),
                Text('Sri Lanka Realtime Tour Matching', style: TextStyle(fontSize: 10, color: Color(0xFF64748B))),
              ],
            ),
          ],
        ),
        actions: [
          IconButton(
            onPressed: () {
              Navigator.push(context, MaterialPageRoute(builder: (_) => const GuideModeScreen()));
            },
            icon: const Icon(Icons.person_pin, color: Color(0xFF0F172A)),
            tooltip: 'Guide Mode',
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Hub Chips
            SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(
                children: _hubs.map((hub) {
                  final isSelected = _selectedHub['name'] == hub['name'];
                  return Padding(
                    padding: const EdgeInsets.only(right: 8),
                    child: ChoiceChip(
                      label: Text(hub['name']),
                      selected: isSelected,
                      selectedColor: const Color(0xFF10B981),
                      labelStyle: TextStyle(
                        color: isSelected ? Colors.white : const Color(0xFF334155),
                        fontWeight: FontWeight.w600,
                        fontSize: 12,
                      ),
                      onSelected: (val) {
                        if (val) {
                          setState(() {
                            _selectedHub = hub;
                          });
                          _loadGuides();
                        }
                      },
                    ),
                  );
                }).toList(),
              ),
            ),
            const SizedBox(height: 16),

            // Radar Map HUD
            RadarMapWidget(
              guides: _guides,
              selectedGuide: _selectedGuide,
              onSelectGuide: (g) => setState(() => _selectedGuide = g),
              locationName: _selectedHub['name'],
            ),
            const SizedBox(height: 20),

            // AI Search Banner
            Container(
              padding: const EdgeInsets.all(14),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: const Color(0xFFE2E8F0)),
              ),
              child: Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: const Color(0xFFECFDF5),
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: const Icon(Icons.auto_awesome, color: Color(0xFF10B981), size: 18),
                  ),
                  const SizedBox(width: 12),
                  const Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('AI Natural-Language Filter', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Color(0xFF0F172A))),
                        Text('e.g. "Surfing Weligama under \$35 speaking English"', style: TextStyle(fontSize: 11, color: Color(0xFF64748B))),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),

            // Guides List Header
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'NEARBY GUIDES (${_guides.length})',
                  style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 12, color: Color(0xFF64748B), letterSpacing: 0.5),
                ),
                const Row(
                  children: [
                    Icon(Icons.bolt, color: Color(0xFF10B981), size: 14),
                    SizedBox(width: 2),
                    Text('Live Presence', style: TextStyle(fontSize: 11, color: Color(0xFF10B981), fontWeight: FontWeight.w600)),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 12),

            if (_loading)
              const Center(child: CircularProgressIndicator(color: Color(0xFF10B981)))
            else
              ..._guides.map((g) {
                return GuideTile(
                  guide: g,
                  isSelected: _selectedGuide?.id == g.id,
                  onTap: () => setState(() => _selectedGuide = g),
                  onChat: () {
                    Navigator.push(context, MaterialPageRoute(builder: (_) => LiveChatScreen(guide: g)));
                  },
                  onRequest: () => _requestTour(g),
                );
              }),
          ],
        ),
      ),
    );
  }
}
