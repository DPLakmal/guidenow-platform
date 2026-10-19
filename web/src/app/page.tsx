'use client';

import React, { useState, useEffect } from 'react';
import { Guide, AISearchFilter } from '@/lib/types';
import { fetchNearbyGuides } from '@/lib/api';
import RadarMap from '@/components/RadarMap';
import AISearchBar from '@/components/AISearchBar';
import GuideCard from '@/components/GuideCard';
import ChatDrawer from '@/components/ChatDrawer';
import BookingModal from '@/components/BookingModal';
import GuidePortal from '@/components/GuidePortal';
import { Compass, Sparkles, MapPin, Radio, ShieldCheck, Layers, ExternalLink } from 'lucide-react';

const TOURISM_HUBS = [
  { name: 'Galle Fort', lat: 6.0329, lng: 80.2168 },
  { name: 'Colombo', lat: 6.9271, lng: 79.8612 },
  { name: 'Weligama & Mirissa', lat: 5.9723, lng: 80.4287 },
  { name: 'Kandy', lat: 7.2906, lng: 80.6337 },
  { name: 'Sigiriya & Dambulla', lat: 7.9570, lng: 80.7603 },
  { name: 'Ella', lat: 6.8667, lng: 81.0466 },
];

export default function Home() {
  const [guides, setGuides] = useState<Guide[]>([]);
  const [selectedHub, setSelectedHub] = useState(TOURISM_HUBS[0]);
  const [selectedGuide, setSelectedGuide] = useState<Guide | null>(null);
  const [chatGuide, setChatGuide] = useState<Guide | null>(null);
  const [bookingGuide, setBookingGuide] = useState<Guide | null>(null);
  const [activeTab, setActiveTab] = useState<'tourist' | 'guide'>('tourist');
  const [aiFilter, setAiFilter] = useState<AISearchFilter | null>(null);

  useEffect(() => {
    async function load() {
      const data = await fetchNearbyGuides(selectedHub.lat, selectedHub.lng, 25);
      setGuides(data);
      if (data.length > 0) {
        setSelectedGuide(data[0]);
      }
    }
    load();
  }, [selectedHub]);

  const handleFilterApplied = (filter: AISearchFilter) => {
    setAiFilter(filter);

    // If AI extracted a location, switch hub automatically
    if (filter.location) {
      const found = TOURISM_HUBS.find((h) =>
        h.name.toLowerCase().includes(filter.location!.toLowerCase())
      );
      if (found) {
        setSelectedHub(found);
      }
    }
  };

  // Filter guides based on AI criteria or city
  const filteredGuides = guides.filter((g) => {
    if (aiFilter?.max_rate && g.hourly_rate > aiFilter.max_rate) return false;
    if (aiFilter?.immediate_only && !g.is_available) return false;
    if (aiFilter?.languages && aiFilter.languages.length > 0) {
      const hasLang = aiFilter.languages.some((l) =>
        g.languages.some((gl) => gl.toLowerCase() === l.toLowerCase())
      );
      if (!hasLang) return false;
    }
    return true;
  });

  return (
    <main className="min-h-screen bg-slate-50 text-slate-900 pb-16">
      {/* Top Navigation Bar */}
      <header className="sticky top-0 z-40 bg-white/90 backdrop-blur-md border-b border-slate-200 shadow-xs">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-2xl bg-gradient-to-tr from-emerald-600 to-teal-500 text-white flex items-center justify-center shadow-md shadow-emerald-600/20">
              <Compass className="w-6 h-6 animate-spin" style={{ animationDuration: '30s' }} />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h1 className="font-extrabold text-lg tracking-tight text-slate-900">GuideNow</h1>
                <span className="text-[10px] font-bold uppercase tracking-wider bg-emerald-100 text-emerald-800 px-2 py-0.5 rounded-full border border-emerald-200">
                  Sri Lanka
                </span>
              </div>
              <p className="text-[11px] text-slate-500 hidden sm:block">
                Real-Time Tour Guide Matching Platform
              </p>
            </div>
          </div>

          {/* Mode Switcher */}
          <div className="flex items-center gap-2">
            <div className="bg-slate-100 p-1 rounded-2xl border border-slate-200 flex text-xs font-semibold">
              <button
                onClick={() => setActiveTab('tourist')}
                className={`px-3 py-1.5 rounded-xl transition-all cursor-pointer ${
                  activeTab === 'tourist'
                    ? 'bg-white text-emerald-700 shadow-sm'
                    : 'text-slate-600 hover:text-slate-900'
                }`}
              >
                Tourist Explorer
              </button>
              <button
                onClick={() => setActiveTab('guide')}
                className={`px-3 py-1.5 rounded-xl transition-all cursor-pointer ${
                  activeTab === 'guide'
                    ? 'bg-white text-emerald-700 shadow-sm'
                    : 'text-slate-600 hover:text-slate-900'
                }`}
              >
                Guide Portal
              </button>
            </div>

            <a
              href="https://github.com/DPLakmal/guidenow-platform"
              target="_blank"
              rel="noreferrer"
              className="p-2 text-slate-500 hover:text-slate-900 bg-slate-100 hover:bg-slate-200 rounded-xl transition-colors hidden md:flex items-center gap-1 text-xs font-semibold"
            >
              <ExternalLink className="w-3.5 h-3.5" />
              <span>GitHub</span>
            </a>
          </div>
        </div>
      </header>

      {/* Main Content Area */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-6">
        {activeTab === 'guide' ? (
          <GuidePortal guides={guides} />
        ) : (
          <div className="space-y-6">
            {/* Hub Selector Pills */}
            <div className="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-none">
              <span className="text-xs font-bold uppercase tracking-wider text-slate-400 mr-1 flex items-center gap-1">
                <MapPin className="w-3.5 h-3.5 text-emerald-600" /> Hubs:
              </span>
              {TOURISM_HUBS.map((hub) => (
                <button
                  key={hub.name}
                  onClick={() => setSelectedHub(hub)}
                  className={`px-3.5 py-1.5 rounded-2xl text-xs font-semibold whitespace-nowrap transition-all cursor-pointer ${
                    selectedHub.name === hub.name
                      ? 'bg-emerald-600 text-white shadow-md shadow-emerald-600/20'
                      : 'bg-white text-slate-700 hover:bg-slate-100 border border-slate-200'
                  }`}
                >
                  {hub.name}
                </button>
              ))}
            </div>

            {/* AI Natural Language Search Header */}
            <AISearchBar onFilterApplied={handleFilterApplied} />

            {/* Split View: Radar Map HUD & Guide List */}
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
              {/* Radar Map (5 Columns on Desktop) */}
              <div className="lg:col-span-5 sticky top-20 space-y-4">
                <RadarMap
                  guides={filteredGuides}
                  selectedGuide={selectedGuide}
                  onSelectGuide={(g) => setSelectedGuide(g)}
                  currentHub={selectedHub.name}
                />

                {/* Microservice Architecture Pill */}
                <div className="p-4 bg-white rounded-2xl border border-slate-200 text-xs text-slate-600 space-y-2 shadow-xs">
                  <div className="flex items-center justify-between font-bold text-slate-800">
                    <span className="flex items-center gap-1.5">
                      <Layers className="w-4 h-4 text-emerald-600" />
                      Live Go Microservices Architecture
                    </span>
                    <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                  </div>
                  <div className="grid grid-cols-2 gap-2 text-[11px] pt-1">
                    <div className="p-2 bg-slate-50 rounded-xl border border-slate-100">
                      <div className="text-slate-400">Gateway</div>
                      <div className="font-semibold text-slate-800">:8080 (Reverse Proxy)</div>
                    </div>
                    <div className="p-2 bg-slate-50 rounded-xl border border-slate-100">
                      <div className="text-slate-400">Matching & Geo</div>
                      <div className="font-semibold text-slate-800">:8082 (Redis GEOSEARCH)</div>
                    </div>
                    <div className="p-2 bg-slate-50 rounded-xl border border-slate-100">
                      <div className="text-slate-400">Chat & Presence</div>
                      <div className="font-semibold text-slate-800">:8083 (WebSocket)</div>
                    </div>
                    <div className="p-2 bg-slate-50 rounded-xl border border-slate-100">
                      <div className="text-slate-400">AI Query Parser</div>
                      <div className="font-semibold text-slate-800">:8084 (Gemini / NLP)</div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Guide Cards Grid (7 Columns on Desktop) */}
              <div className="lg:col-span-7 space-y-4">
                <div className="flex items-center justify-between">
                  <h2 className="text-sm font-bold uppercase tracking-wider text-slate-700 flex items-center gap-1.5">
                    <span>Verified Nearby Guides</span>
                    <span className="text-xs font-normal text-slate-400">({filteredGuides.length} active)</span>
                  </h2>

                  {aiFilter && (
                    <button
                      onClick={() => setAiFilter(null)}
                      className="text-xs text-emerald-600 hover:text-emerald-800 font-semibold cursor-pointer"
                    >
                      Reset AI Filters
                    </button>
                  )}
                </div>

                {filteredGuides.length === 0 ? (
                  <div className="bg-white rounded-3xl p-12 text-center border border-slate-200">
                    <Compass className="w-12 h-12 text-slate-300 mx-auto mb-3" />
                    <h3 className="font-bold text-slate-800">No Guides Match Current Filter</h3>
                    <p className="text-xs text-slate-500 mt-1 max-w-sm mx-auto">
                      Try broadening your search criteria or switching to a nearby hub such as Galle Fort or Colombo.
                    </p>
                  </div>
                ) : (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {filteredGuides.map((guide) => (
                      <GuideCard
                        key={guide.id}
                        guide={guide}
                        isSelected={selectedGuide?.id === guide.id}
                        onSelect={() => setSelectedGuide(guide)}
                        onOpenChat={(g) => setChatGuide(g)}
                        onRequestBooking={(g) => setBookingGuide(g)}
                      />
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Slide-over WebSocket Chat Drawer */}
      <ChatDrawer
        guide={chatGuide}
        isOpen={!!chatGuide}
        onClose={() => setChatGuide(null)}
      />

      {/* Quick Tour Request Booking Modal */}
      <BookingModal
        guide={bookingGuide}
        isOpen={!!bookingGuide}
        onClose={() => setBookingGuide(null)}
      />
    </main>
  );
}
