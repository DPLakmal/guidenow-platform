class Guide {
  final String id;
  final String name;
  final String bio;
  final String avatar;
  final String phone;
  final String city;
  final double hourlyRate;
  final double rating;
  final int reviewCount;
  final List<String> languages;
  final List<String> specialties;
  final double currentLat;
  final double currentLng;
  final double? distanceKm;
  final bool isAvailable;
  final bool verifiedStatus;
  final String licenseNumber;

  Guide({
    required this.id,
    required this.name,
    required this.bio,
    required this.avatar,
    required this.phone,
    required this.city,
    required this.hourlyRate,
    required this.rating,
    required this.reviewCount,
    required this.languages,
    required this.specialties,
    required this.currentLat,
    required this.currentLng,
    this.distanceKm,
    required this.isAvailable,
    required this.verifiedStatus,
    required this.licenseNumber,
  });

  factory Guide.fromJson(Map<String, dynamic> json) {
    return Guide(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      bio: json['bio'] ?? '',
      avatar: json['avatar'] ?? '',
      phone: json['phone'] ?? '',
      city: json['city'] ?? '',
      hourlyRate: (json['hourly_rate'] as num?)?.toDouble() ?? 20.0,
      rating: (json['rating'] as num?)?.toDouble() ?? 5.0,
      reviewCount: json['review_count'] ?? 0,
      languages: List<String>.from(json['languages'] ?? []),
      specialties: List<String>.from(json['specialties'] ?? []),
      currentLat: (json['current_lat'] as num?)?.toDouble() ?? 6.0,
      currentLng: (json['current_lng'] as num?)?.toDouble() ?? 80.0,
      distanceKm: (json['distance_km'] as num?)?.toDouble(),
      isAvailable: json['is_available'] ?? true,
      verifiedStatus: json['verified_status'] ?? true,
      licenseNumber: json['license_number'] ?? '',
    );
  }
}
