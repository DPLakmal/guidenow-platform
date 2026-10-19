'use client';

import React, { useState } from 'react';
import { Sparkles, Search, ArrowRight, CheckCircle2, DollarSign, MapPin, Globe } from 'lucide-react';
import { AISearchFilter } from '@/lib/types';
import { parseAIQuery } from '@/lib/api';

interface AISearchBarProps {
  onFilterApplied: (filter: AISearchFilter) => void;
}

const PRESET_QUERIES = [
  '🏰 Galle Fort history & Dutch architecture under $25/hr',
  '🏄 Surfing coach in Weligama under $35 available right now',
  '🦁 Sigiriya Lion Rock fortress guide speaking French',
  '🫖 Ella Nine Arch Bridge & tea plantation hiking guide',
  '🥘 Colombo street food & Pettah walking tour with Tamil/English',
];

export default function AISearchBar({ onFilterApplied }: AISearchBarProps) {
  const [query, setQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [activeFilter, setActiveFilter] = useState<AISearchFilter | null>(null);

  const handleSearch = async (textToSearch?: string) => {
    const q = textToSearch || query;
    if (!q.trim()) return;

    setLoading(true);
    try {
      const filter = await parseAIQuery(q);
      setActiveFilter(filter);
      onFilterApplied(filter);
    } finally {
      setLoading(false);
    }
  };

  const handlePresetClick = (preset: string) => {
    // Strip emoji prefix
    const cleanText = preset.replace(/^[^\w]+/, '').trim();
    setQuery(cleanText);
    handleSearch(cleanText);
  };

  return (
    <div className="w-full bg-white rounded-3xl p-6 border border-slate-200/80 shadow-xl shadow-slate-200/50">
      <div className="flex items-center gap-2 mb-3">
        <div className="p-1.5 rounded-xl bg-gradient-to-tr from-emerald-500 to-teal-400 text-white shadow-sm">
          <Sparkles className="w-4 h-4" />
        </div>
        <h2 className="text-sm font-bold uppercase tracking-wider text-slate-800">
          AI Natural-Language Guide Finder
        </h2>
        <span className="ml-auto text-[11px] font-semibold text-emerald-700 bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-200">
          Go AI Service (:8084)
        </span>
      </div>

      {/* Main Input Field */}
      <form
        onSubmit={(e) => {
          e.preventDefault();
          handleSearch();
        }}
        className="relative flex items-center"
      >
        <Search className="absolute left-4 w-5 h-5 text-slate-400" />
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="e.g. Need an English or German guide in Galle Fort for architecture under $25/hr available now..."
          className="w-full pl-12 pr-32 py-4 bg-slate-50 border border-slate-200 rounded-2xl text-slate-800 placeholder-slate-400 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:bg-white transition-all shadow-inner"
        />
        <button
          type="submit"
          disabled={loading || !query.trim()}
          className="absolute right-2 px-5 py-2.5 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 text-white rounded-xl text-xs font-semibold flex items-center gap-1.5 shadow-md shadow-emerald-600/20 transition-all cursor-pointer"
        >
          {loading ? (
            <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
          ) : (
            <>
              <span>Extract</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </>
          )}
        </button>
      </form>

      {/* Quick Suggestion Pills */}
      <div className="mt-3 flex flex-wrap gap-1.5">
        <span className="text-[11px] text-slate-400 self-center mr-1">Try:</span>
        {PRESET_QUERIES.map((preset, i) => (
          <button
            key={i}
            type="button"
            onClick={() => handlePresetClick(preset)}
            className="text-[11px] text-slate-600 hover:text-emerald-700 bg-slate-100 hover:bg-emerald-50 hover:border-emerald-200 border border-slate-200/80 px-2.5 py-1 rounded-xl transition-all cursor-pointer"
          >
            {preset}
          </button>
        ))}
      </div>

      {/* Structured Filter Output Badge */}
      {activeFilter && (
        <div className="mt-4 p-3.5 bg-emerald-50/70 border border-emerald-200/90 rounded-2xl text-xs text-slate-700 animate-fadeIn">
          <div className="flex items-center gap-1.5 font-semibold text-emerald-900 mb-2">
            <CheckCircle2 className="w-4 h-4 text-emerald-600" />
            <span>AI Intent Parsed ({Math.round(activeFilter.confidence * 100)}% Confidence):</span>
          </div>
          <div className="flex flex-wrap gap-2 items-center">
            {activeFilter.location && (
              <span className="inline-flex items-center gap-1 bg-white border border-emerald-300 text-emerald-800 font-medium px-2.5 py-1 rounded-lg shadow-xs">
                <MapPin className="w-3 h-3 text-emerald-600" />
                {activeFilter.location}
              </span>
            )}
            {activeFilter.max_rate ? (
              <span className="inline-flex items-center gap-1 bg-white border border-emerald-300 text-emerald-800 font-medium px-2.5 py-1 rounded-lg shadow-xs">
                <DollarSign className="w-3 h-3 text-emerald-600" />
                Max ${activeFilter.max_rate}/hr
              </span>
            ) : null}
            {activeFilter.languages && activeFilter.languages.length > 0 && (
              <span className="inline-flex items-center gap-1 bg-white border border-emerald-300 text-emerald-800 font-medium px-2.5 py-1 rounded-lg shadow-xs">
                <Globe className="w-3 h-3 text-emerald-600" />
                {activeFilter.languages.join(', ')}
              </span>
            )}
            {activeFilter.specialties && activeFilter.specialties.map((spec, i) => (
              <span key={i} className="bg-emerald-100/90 text-emerald-800 font-medium px-2.5 py-1 rounded-lg">
                {spec}
              </span>
            ))}
            {activeFilter.immediate_only && (
              <span className="bg-amber-100 text-amber-800 font-medium px-2.5 py-1 rounded-lg">
                ⚡ Available Now
              </span>
            )}
          </div>
          {activeFilter.explanation && (
            <p className="mt-2 text-[11px] text-slate-500 italic">
              &quot;{activeFilter.explanation}&quot;
            </p>
          )}
        </div>
      )}
    </div>
  );
}
