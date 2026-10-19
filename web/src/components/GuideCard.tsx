'use client';

import React from 'react';
import { Guide } from '@/lib/types';
import { Star, ShieldCheck, MessageSquare, Compass, Languages, Clock } from 'lucide-react';

interface GuideCardProps {
  guide: Guide;
  isSelected?: boolean;
  onSelect: () => void;
  onOpenChat: (guide: Guide) => void;
  onRequestBooking: (guide: Guide) => void;
}

export default function GuideCard({
  guide,
  isSelected,
  onSelect,
  onOpenChat,
  onRequestBooking,
}: GuideCardProps) {
  return (
    <div
      onClick={onSelect}
      className={`relative bg-white rounded-3xl p-5 border transition-all duration-300 cursor-pointer ${
        isSelected
          ? 'border-emerald-500 shadow-xl shadow-emerald-500/10 ring-2 ring-emerald-500/20'
          : 'border-slate-200/80 hover:border-slate-300 hover:shadow-lg'
      }`}
    >
      <div className="flex items-start gap-4">
        {/* Guide Avatar */}
        <div className="relative w-16 h-16 rounded-2xl overflow-hidden bg-slate-100 flex-shrink-0 border border-slate-200">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={guide.avatar}
            alt={guide.name}
            className="w-full h-full object-cover"
          />
          {guide.is_available && (
            <div className="absolute top-1 right-1 w-3.5 h-3.5 bg-emerald-500 border-2 border-white rounded-full shadow" title="Available Now" />
          )}
        </div>

        {/* Info Header */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between gap-2">
            <h3 className="font-bold text-slate-900 text-base truncate flex items-center gap-1.5">
              <span>{guide.name}</span>
              {guide.verified_status && (
                <span title={`Verified: ${guide.license_number}`} className="flex items-center">
                  <ShieldCheck className="w-4 h-4 text-emerald-600 flex-shrink-0" />
                </span>
              )}
            </h3>
            <div className="text-right">
              <span className="text-lg font-extrabold text-emerald-700">${guide.hourly_rate}</span>
              <span className="text-[11px] text-slate-400 font-medium">/hr</span>
            </div>
          </div>

          <div className="flex items-center gap-3 text-xs text-slate-500 mt-1">
            <div className="flex items-center gap-1 text-amber-500 font-bold">
              <Star className="w-3.5 h-3.5 fill-amber-400 text-amber-400" />
              <span>{guide.rating.toFixed(2)}</span>
              <span className="text-slate-400 font-normal">({guide.review_count})</span>
            </div>
            <span>•</span>
            <span className="font-medium text-slate-700">{guide.city}</span>
            {guide.distance_km !== undefined && (
              <>
                <span>•</span>
                <span className="text-emerald-600 font-semibold">{guide.distance_km} km away</span>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Bio */}
      <p className="mt-3 text-xs text-slate-600 line-clamp-2 leading-relaxed">
        {guide.bio}
      </p>

      {/* Languages & Specialties */}
      <div className="mt-3 pt-3 border-t border-slate-100 flex flex-wrap gap-1.5">
        <div className="flex items-center gap-1 text-[11px] text-slate-500 mr-1">
          <Languages className="w-3 h-3 text-slate-400" />
          <span>{guide.languages.join(', ')}</span>
        </div>
      </div>

      <div className="mt-2 flex flex-wrap gap-1.5">
        {guide.specialties.map((spec, i) => (
          <span
            key={i}
            className="text-[10px] font-semibold text-slate-600 bg-slate-100 px-2 py-0.5 rounded-md"
          >
            {spec}
          </span>
        ))}
      </div>

      {/* Actions */}
      <div className="mt-4 pt-3 border-t border-slate-100 flex items-center gap-2">
        <button
          onClick={(e) => {
            e.stopPropagation();
            onOpenChat(guide);
          }}
          className="flex-1 py-2.5 px-3 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs rounded-xl flex items-center justify-center gap-1.5 transition-colors cursor-pointer"
        >
          <MessageSquare className="w-3.5 h-3.5 text-slate-500" />
          <span>Live Chat</span>
        </button>

        <button
          onClick={(e) => {
            e.stopPropagation();
            onRequestBooking(guide);
          }}
          className="flex-1 py-2.5 px-3 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs rounded-xl flex items-center justify-center gap-1.5 shadow-sm shadow-emerald-600/20 transition-all cursor-pointer"
        >
          <Compass className="w-3.5 h-3.5" />
          <span>Request Tour</span>
        </button>
      </div>
    </div>
  );
}
