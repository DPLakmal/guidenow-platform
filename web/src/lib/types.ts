export interface Guide {
  id: string;
  name: string;
  bio: string;
  avatar: string;
  phone: string;
  email: string;
  city: string;
  hourly_rate: number;
  rating: number;
  review_count: number;
  languages: string[];
  specialties: string[];
  current_lat: number;
  current_lng: number;
  distance_km?: number;
  is_available: boolean;
  verified_status: boolean;
  license_number: string;
}

export interface AISearchFilter {
  original_query: string;
  location?: string;
  latitude?: number;
  longitude?: number;
  radius_km?: number;
  languages?: string[];
  specialties?: string[];
  max_rate?: number;
  min_rating?: number;
  immediate_only: boolean;
  confidence: number;
  explanation?: string;
}

export interface TourRequest {
  id: string;
  tourist_id: string;
  tourist_name: string;
  tourist_phone?: string;
  guide_id: string;
  guide_name: string;
  status: 'REQUESTED' | 'ACCEPTED' | 'DECLINED' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELLED';
  meeting_point: string;
  destination: string;
  duration_hours: number;
  hourly_rate: number;
  total_estimated: number;
  special_notes?: string;
  created_at: string;
}

export interface ChatMessage {
  id: string;
  room_id: string;
  sender_id: string;
  sender_name: string;
  sender_role: string;
  content: string;
  timestamp: string;
  read?: boolean;
}
