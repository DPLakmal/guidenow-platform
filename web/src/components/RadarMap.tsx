'use client';

import React from 'react';
import { Guide } from '@/lib/types';
import { Compass, Navigation } from 'lucide-react';

interface RadarMapProps {
  guides: Guide[];
  selectedGuide: Guide | null;
  onSelectGuide: (guide: Guide) => void;
  currentHub: string;
}

export default function RadarMap({ guides, selectedGuide, onSelectGuide, currentHub }: RadarMapProps) {
  return (
    <div className="relative w-full aspect-[4/3] md:aspect-[16/9] max-h-[460px] bg-slate-900 rounded-3xl overflow-hidden border border-slate-800 shadow-2xl flex items-center justify-center">
      {/* Radar Grid Lines */}
      <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
        {/* Concentric distance rings */}
        <div className="w-[85%] h-[85%] rounded-full border border-emerald-500/20" />
        <div className="w-[60%] h-[60%] rounded-full border border-emerald-500/25" />
        <div className="w-[35%] h-[35%] rounded-full border border-emerald-500/30" />
        <div className="w-[12%] h-[12%] rounded-full border border-emerald-500/40 bg-emerald-500/10" />

        {/* Crosshair lines */}
        <div className="absolute w-full h-[1px] bg-emerald-500/15" />
        <div className="absolute h-full w-[1px] bg-emerald-500/15" />

        {/* Rotating Radar Sweep */}
        <div className="absolute w-[85%] h-[85%] rounded-full pointer-events-none radar-sweep">
          <div className="w-1/2 h-1/2 bg-gradient-to-br from-emerald-500/20 to-transparent origin-bottom-right rounded-tl-full" />
        </div>
      </div>

      {/* Center User Pin */}
      <div className="relative z-10 flex flex-col items-center">
        <div className="w-5 h-5 rounded-full bg-emerald-500 ring-4 ring-emerald-500/30 flex items-center justify-center animate-pulse">
          <div className="w-2 h-2 rounded-full bg-white" />
        </div>
        <span className="mt-1 text-[11px] font-semibold text-emerald-400 bg-slate-900/90 px-2 py-0.5 rounded-full border border-emerald-500/30 shadow">
          You ({currentHub})
        </span>
      </div>

      {/* Guide Radar Blips */}
      {guides.map((guide, idx) => {
        // Distribute blips visually on the radar based on relative coordinates / index
        const angle = (idx * (360 / Math.max(guides.length, 1)) + 45) * (Math.PI / 180);
        const radiusPercent = Math.min(25 + idx * 8, 42); // Between 25% and 42% of container
        const leftPercent = 50 + radiusPercent * Math.cos(angle);
        const topPercent = 50 + radiusPercent * Math.sin(angle);

        const isSelected = selectedGuide?.id === guide.id;

        return (
          <button
            key={guide.id}
            onClick={() => onSelectGuide(guide)}
            style={{ left: `${leftPercent}%`, top: `${topPercent}%` }}
            className={`absolute z-20 -translate-x-1/2 -translate-y-1/2 group transition-all duration-300 focus:outline-none`}
          >
            {/* Ping animation */}
            <div className={`absolute -inset-2 rounded-full ${isSelected ? 'bg-amber-400/40 animate-ping' : 'bg-emerald-400/20 group-hover:animate-ping'}`} />

            <div className={`relative flex items-center gap-2 px-2.5 py-1.5 rounded-2xl border transition-all ${
              isSelected
                ? 'bg-amber-500 text-slate-950 font-bold border-amber-300 scale-110 shadow-lg shadow-amber-500/30'
                : 'bg-slate-800/90 hover:bg-slate-800 text-slate-200 border-slate-700/80 shadow-md hover:scale-105'
            }`}>
              <div className="w-6 h-6 rounded-full overflow-hidden border border-white/20 flex-shrink-0">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src={guide.avatar} alt={guide.name} className="w-full h-full object-cover" />
              </div>
              <div className="text-left text-[11px] leading-tight">
                <div className="font-semibold truncate max-w-[100px]">{guide.name.split(' ')[0]}</div>
                <div className="text-[10px] opacity-80">${guide.hourly_rate}/hr • {guide.distance_km || 0.8}km</div>
              </div>
            </div>
          </button>
        );
      })}

      {/* Radar HUD Overlay Controls */}
      <div className="absolute top-4 left-4 z-20 flex items-center gap-2 bg-slate-900/80 backdrop-blur-md px-3 py-1.5 rounded-xl border border-slate-800 text-xs text-slate-300">
        <Compass className="w-4 h-4 text-emerald-400 animate-spin" style={{ animationDuration: '20s' }} />
        <span>Live Geospatial Radar</span>
        <span className="w-2 h-2 rounded-full bg-emerald-500 animate-ping" />
      </div>

      <div className="absolute bottom-4 right-4 z-20 flex items-center gap-2 bg-slate-900/80 backdrop-blur-md px-3 py-1.5 rounded-xl border border-slate-800 text-[11px] text-slate-400">
        <Navigation className="w-3.5 h-3.5 text-emerald-400" />
        <span>Radius: 25 km • {guides.length} Guides Online</span>
      </div>
    </div>
  );
}
