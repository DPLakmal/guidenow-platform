class TourRequest {
  final String id;
  final String touristId;
  final String touristName;
  final String guideId;
  final String guideName;
  final String status;
  final String meetingPoint;
  final String destination;
  final double durationHours;
  final double hourlyRate;
  final double totalEstimated;
  final String? specialNotes;

  TourRequest({
    required this.id,
    required this.touristId,
    required this.touristName,
    required this.guideId,
    required this.guideName,
    required this.status,
    required this.meetingPoint,
    required this.destination,
    required this.durationHours,
    required this.hourlyRate,
    required this.totalEstimated,
    this.specialNotes,
  });

  factory TourRequest.fromJson(Map<String, dynamic> json) {
    return TourRequest(
      id: json['id'] ?? '',
      touristId: json['tourist_id'] ?? '',
      touristName: json['tourist_name'] ?? '',
      guideId: json['guide_id'] ?? '',
      guideName: json['guide_name'] ?? '',
      status: json['status'] ?? 'REQUESTED',
      meetingPoint: json['meeting_point'] ?? '',
      destination: json['destination'] ?? '',
      durationHours: (json['duration_hours'] as num?)?.toDouble() ?? 2.0,
      hourlyRate: (json['hourly_rate'] as num?)?.toDouble() ?? 20.0,
      totalEstimated: (json['total_estimated'] as num?)?.toDouble() ?? 40.0,
      specialNotes: json['special_notes'],
    );
  }
}
