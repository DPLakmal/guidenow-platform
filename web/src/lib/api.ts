import { Guide, AISearchFilter, TourRequest, ChatMessage } from './types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function fetchNearbyGuides(lat = 6.0329, lng = 80.2168, radius = 25): Promise<Guide[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/matching/nearby?lat=${lat}&lng=${lng}&radius=${radius}`);
    if (!res.ok) throw new Error('API request failed');
    const json = await res.json();
    return json.data?.guides || [];
  } catch (err) {
    // Fallback static guides for standalone dev
    return [
      {
        id: 'guide-galle-01',
        name: 'Chaminda Jayasuriya',
        bio: 'Born and raised inside Galle Fort. Specialized in 17th-century Dutch colonial ramparts and hidden gem cafes.',
        avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=400&q=80',
        phone: '+94 76 890 1234',
        email: 'chaminda.galle@guidenow.lk',
        city: 'Galle',
        hourly_rate: 25.0,
        rating: 4.98,
        review_count: 210,
        languages: ['English', 'French', 'Sinhala'],
        specialties: ['Galle Fort History', 'Dutch Architecture', 'Maritime Heritage'],
        current_lat: 6.0329,
        current_lng: 80.2168,
        distance_km: 0.4,
        is_available: true,
        verified_status: true,
        license_number: 'SLTDA/NTG/2016/0891',
      },
      {
        id: 'guide-galle-02',
        name: 'Nadeesha Fernando',
        bio: 'Certified naturalist guide for southern coastal wetlands, cinnamon plantations, and artisanal lace-making.',
        avatar: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=crop&w=400&q=80',
        phone: '+94 77 987 6543',
        email: 'nadeesha.f@guidenow.lk',
        city: 'Galle',
        hourly_rate: 20.0,
        rating: 4.91,
        review_count: 76,
        languages: ['English', 'German', 'Sinhala'],
        specialties: ['Artisan Crafts', 'Cinnamon Plantations', 'Coastal Ecology'],
        current_lat: 6.0354,
        current_lng: 80.2195,
        distance_km: 1.2,
        is_available: true,
        verified_status: true,
        license_number: 'SLTDA/CG/2021/3342',
      },
      {
        id: 'guide-weli-01',
        name: 'Tharindu Mendis',
        bio: 'ISA certified surf coach and marine wildlife guide. Expert in Weligama bay breaks and Mirissa blue whales.',
        avatar: 'https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?auto=format&fit=crop&w=400&q=80',
        phone: '+94 77 654 3210',
        email: 'tharindu.surf@guidenow.lk',
        city: 'Weligama',
        hourly_rate: 30.0,
        rating: 4.96,
        review_count: 195,
        languages: ['English', 'Spanish', 'Sinhala'],
        specialties: ['Surfing Lessons', 'Whale Watching', 'Ocean Ecology'],
        current_lat: 5.9723,
        current_lng: 80.4287,
        distance_km: 18.5,
        is_available: true,
        verified_status: true,
        license_number: 'SLTDA/CG/2019/1890',
      },
    ];
  }
}

export async function parseAIQuery(query: string): Promise<AISearchFilter> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/ai/parse-query`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query }),
    });
    if (!res.ok) throw new Error('AI parse failed');
    const json = await res.json();
    return json.data;
  } catch (err) {
    return {
      original_query: query,
      location: 'Galle',
      languages: ['English'],
      specialties: ['History & Culture'],
      max_rate: 30,
      immediate_only: true,
      confidence: 0.95,
      explanation: 'Extracted: Galle area, English speaking, under $30/hr, immediate availability.',
    };
  }
}

export async function createTourRequest(req: Partial<TourRequest>): Promise<TourRequest> {
  const res = await fetch(`${API_BASE}/api/v1/bookings/request`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  });
  const json = await res.json();
  return json.data;
}

export async function fetchChatHistory(roomId: string): Promise<ChatMessage[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/chat/rooms/${roomId}/messages`);
    const json = await res.json();
    return json.data || [];
  } catch {
    return [];
  }
}
