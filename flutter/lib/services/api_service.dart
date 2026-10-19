import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/guide.dart';
import '../models/booking.dart';

class ApiService {
  static const String baseUrl = 'http://10.0.2.2:8080'; // Android emulator localhost alias

  static List<Guide> getFallbackGuides() {
    return [
      Guide(
        id: 'guide-galle-01',
        name: 'Chaminda Jayasuriya',
        bio: 'Born and raised inside Galle Fort. Specialized in 17th-century Dutch colonial ramparts.',
        avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=400&q=80',
        phone: '+94 76 890 1234',
        city: 'Galle',
        hourlyRate: 25.0,
        rating: 4.98,
        reviewCount: 210,
        languages: ['English', 'French', 'Sinhala'],
        specialties: ['Galle Fort History', 'Dutch Architecture'],
        currentLat: 6.0329,
        currentLng: 80.2168,
        distanceKm: 0.4,
        isAvailable: true,
        verifiedStatus: true,
        licenseNumber: 'SLTDA/NTG/2016/0891',
      ),
      Guide(
        id: 'guide-weli-01',
        name: 'Tharindu Mendis',
        bio: 'ISA certified surf coach and marine wildlife guide in Weligama and Mirissa.',
        avatar: 'https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?auto=format&fit=crop&w=400&q=80',
        phone: '+94 77 654 3210',
        city: 'Weligama',
        hourlyRate: 30.0,
        rating: 4.96,
        reviewCount: 195,
        languages: ['English', 'Spanish', 'Sinhala'],
        specialties: ['Surfing Lessons', 'Whale Watching'],
        currentLat: 5.9723,
        currentLng: 80.4287,
        distanceKm: 18.2,
        isAvailable: true,
        verifiedStatus: true,
        licenseNumber: 'SLTDA/CG/2019/1890',
      ),
    ];
  }

  static Future<List<Guide>> fetchNearbyGuides({double lat = 6.0329, double lng = 80.2168, double radius = 25}) async {
    try {
      final res = await http.get(Uri.parse('$baseUrl/api/v1/matching/nearby?lat=$lat&lng=$lng&radius=$radius'))
          .timeout(const Duration(seconds: 3));
      if (res.statusCode == 200) {
        final data = json.decode(res.body);
        final list = data['data']['guides'] as List;
        return list.map((item) => Guide.fromJson(item)).toList();
      }
    } catch (_) {
      // Fallback
    }
    return getFallbackGuides();
  }

  static Future<TourRequest?> createTourRequest({
    required String guideId,
    required String touristName,
    required String meetingPoint,
    required String destination,
    required double durationHours,
  }) async {
    try {
      final res = await http.post(
        Uri.parse('$baseUrl/api/v1/bookings/request'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'tourist_id': 'tourist-user-mobile-1',
          'tourist_name': touristName,
          'guide_id': guideId,
          'meeting_point': meetingPoint,
          'destination': destination,
          'duration_hours': durationHours,
        }),
      ).timeout(const Duration(seconds: 3));

      if (res.statusCode == 200) {
        final data = json.decode(res.body);
        return TourRequest.fromJson(data['data']);
      }
    } catch (_) {
      // Return local mock
    }
    return TourRequest(
      id: 'tour-mob-${DateTime.now().millisecondsSinceEpoch % 10000}',
      touristId: 'tourist-user-mobile-1',
      touristName: touristName,
      guideId: guideId,
      guideName: 'Chaminda Jayasuriya',
      status: 'REQUESTED',
      meetingPoint: meetingPoint,
      destination: destination,
      durationHours: durationHours,
      hourlyRate: 25.0,
      totalEstimated: durationHours * 25.0,
    );
  }
}
