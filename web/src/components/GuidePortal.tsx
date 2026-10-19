'use client';

import React, { useState } from 'react';
import { Guide, TourRequest } from '@/lib/types';
import { Radio, Power, Check, X, ShieldAlert, Navigation } from 'lucide-react';

interface GuidePortalProps {
  guides: Guide[];
}

export default function GuidePortal({ guides }: GuidePortalProps) {
  const [activeGuide, setActiveGuide] = useState<Guide>(guides[0] || {} as Guide);
  const [isAvailable, setIsAvailable] = useState(true);
  const [broadcastCount, setBroadcastCount] = useState(14);
  const [incomingRequests, setIncomingRequests] = useState<TourRequest[]>([
    {
      id: 'tour-galle-902',
      tourist_id: 'user-78',
      tourist_name: 'Emma Watson',
      guide_id: guides[0]?.id || 'guide-galle-01',
      guide_name: guides[0]?.name || 'Chaminda Jayasuriya',
      status: 'REQUESTED',
      meeting_point: 'Galle Fort Lighthouse',
      destination: 'Dutch Ramparts Sunset Walking Tour',
      duration_hours: 2,
      hourly_rate: guides[0]?.hourly_rate || 25,
      total_estimated: 50,
      special_notes: 'Interested in maritime history and photography spots.',
      created_at: new Date().toISOString(),
    },
  ]);

  const handleToggleAvailability = async () => {
    const nextState = !isAvailable;
    setIsAvailable(nextState);

    try {
      await fetch(`http://localhost:8080/api/v1/guides/${activeGuide.id}/availability`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ available: nextState }),
      });
    } catch {
      // Local fallback
    }
  };

  const handleBroadcastGPS = async () => {
    setBroadcastCount((c) => c + 1);
    try {
      await fetch(`http://localhost:8080/api/v1/matching/location`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          guide_id: activeGuide.id,
          latitude: activeGuide.current_lat + (Math.random() - 0.5) * 0.002,
          longitude: activeGuide.current_lng + (Math.random() - 0.5) * 0.002,
        }),
      });
    } catch {
      // Local fallback
    }
  };

  const handleAccept = (tourId: string) => {
    setIncomingRequests((prev) =>
      prev.map((r) => (r.id === tourId ? { ...r, status: 'ACCEPTED' } : r))
    );
  };

  const handleDecline = (tourId: string) => {
    setIncomingRequests((prev) =>
      prev.map((r) => (r.id === tourId ? { ...r, status: 'DECLINED' } : r))
    );
  };

  return (
    <div className="bg-slate-900 text-white rounded-3xl p-6 border border-slate-800 shadow-2xl">
      <div className="flex flex-wrap items-center justify-between gap-4 pb-5 border-b border-slate-800">
        <div>
          <span className="text-[10px] font-bold uppercase tracking-wider text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-full border border-emerald-500/20">
            Guide Dispatcher Mode
          </span>
          <h2 className="text-lg font-bold mt-1">Guide Availability & Dispatch Center</h2>
          <p className="text-xs text-slate-400">Manage real-time Redis presence and accept incoming tourist requests.</p>
        </div>

        {/* Guide Selector */}
        <select
          value={activeGuide.id}
          onChange={(e) => {
            const found = guides.find((g) => g.id === e.target.value);
            if (found) setActiveGuide(found);
          }}
          className="bg-slate-800 border border-slate-700 rounded-xl px-3 py-2 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
        >
          {guides.map((g) => (
            <option key={g.id} value={g.id}>
              {g.name} ({g.city})
            </option>
          ))}
        </select>
      </div>

      {/* Controls Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 my-6">
        {/* Availability Toggle */}
        <div className="bg-slate-800/60 p-4 rounded-2xl border border-slate-700/60 flex items-center justify-between">
          <div>
            <div className="text-xs font-semibold text-slate-300">Live Status</div>
            <div className="text-sm font-bold text-white flex items-center gap-1.5 mt-0.5">
              <span className={`w-2.5 h-2.5 rounded-full ${isAvailable ? 'bg-emerald-500 animate-pulse' : 'bg-rose-500'}`} />
              {isAvailable ? 'Online & Available' : 'Offline / On Tour'}
            </div>
          </div>
          <button
            onClick={handleToggleAvailability}
            className={`p-2.5 rounded-xl border transition-all cursor-pointer ${
              isAvailable
                ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40 hover:bg-emerald-500/30'
                : 'bg-rose-500/20 text-rose-300 border-rose-500/40 hover:bg-rose-500/30'
            }`}
          >
            <Power className="w-4 h-4" />
          </button>
        </div>

        {/* GPS Broadcast Simulator */}
        <div className="bg-slate-800/60 p-4 rounded-2xl border border-slate-700/60 flex items-center justify-between">
          <div>
            <div className="text-xs font-semibold text-slate-300">Redis Geospatial Pings</div>
            <div className="text-sm font-bold text-white flex items-center gap-1.5 mt-0.5">
              <Navigation className="w-3.5 h-3.5 text-emerald-400" />
              <span>{broadcastCount} Pings Sent</span>
            </div>
          </div>
          <button
            onClick={handleBroadcastGPS}
            className="px-3 py-2 bg-slate-700 hover:bg-slate-600 text-xs font-semibold rounded-xl text-slate-200 border border-slate-600 transition-all cursor-pointer flex items-center gap-1"
          >
            <Radio className="w-3.5 h-3.5 text-emerald-400 animate-pulse" />
            <span>Ping GPS</span>
          </button>
        </div>

        {/* License & Credentials */}
        <div className="bg-slate-800/60 p-4 rounded-2xl border border-slate-700/60">
          <div className="text-xs font-semibold text-slate-300">SLTDA Certification</div>
          <div className="text-xs font-mono font-bold text-emerald-400 mt-1">
            {activeGuide.license_number || 'SLTDA/NTG/2018/1042'}
          </div>
          <div className="text-[10px] text-slate-400 mt-0.5">Verified National Tourist Guide</div>
        </div>
      </div>

      {/* Incoming Requests */}
      <div>
        <h3 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">
          Incoming Tour Requests ({incomingRequests.length})
        </h3>

        <div className="space-y-3">
          {incomingRequests.map((req) => (
            <div
              key={req.id}
              className="bg-slate-800/90 p-4 rounded-2xl border border-slate-700 flex flex-wrap items-center justify-between gap-4"
            >
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <span className="font-bold text-sm text-white">{req.tourist_name}</span>
                  <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                    req.status === 'ACCEPTED'
                      ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/40'
                      : req.status === 'DECLINED'
                      ? 'bg-rose-500/20 text-rose-300 border border-rose-500/40'
                      : 'bg-amber-500/20 text-amber-300 border border-amber-500/40 animate-pulse'
                  }`}>
                    {req.status}
                  </span>
                </div>
                <div className="text-xs text-slate-300 font-medium">
                  {req.destination} • {req.duration_hours}h (${req.total_estimated})
                </div>
                <div className="text-[11px] text-slate-400">
                  Meeting Point: <span className="text-slate-200">{req.meeting_point}</span>
                </div>
              </div>

              {req.status === 'REQUESTED' && (
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => handleAccept(req.id)}
                    className="px-3.5 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl text-xs font-semibold flex items-center gap-1.5 transition-all cursor-pointer shadow-md"
                  >
                    <Check className="w-3.5 h-3.5" />
                    <span>Accept</span>
                  </button>
                  <button
                    onClick={() => handleDecline(req.id)}
                    className="px-3.5 py-2 bg-slate-700 hover:bg-rose-600 text-white rounded-xl text-xs font-semibold flex items-center gap-1.5 transition-all cursor-pointer"
                  >
                    <X className="w-3.5 h-3.5" />
                    <span>Decline</span>
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
