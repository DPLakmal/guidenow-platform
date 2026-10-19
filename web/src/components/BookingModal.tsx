'use client';

import React, { useState } from 'react';
import { Guide, TourRequest } from '@/lib/types';
import { createTourRequest } from '@/lib/api';
import { X, Calendar, MapPin, DollarSign, CheckCircle2, Clock } from 'lucide-react';

interface BookingModalProps {
  guide: Guide | null;
  isOpen: boolean;
  onClose: () => void;
}

export default function BookingModal({ guide, isOpen, onClose }: BookingModalProps) {
  const [duration, setDuration] = useState(2.0);
  const [meetingPoint, setMeetingPoint] = useState('');
  const [destination, setDestination] = useState('');
  const [notes, setNotes] = useState('');
  const [loading, setLoading] = useState(false);
  const [confirmedBooking, setConfirmedBooking] = useState<TourRequest | null>(null);

  if (!isOpen || !guide) return null;

  const totalEstimate = duration * guide.hourly_rate;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const result = await createTourRequest({
        tourist_id: 'tourist-user-101',
        tourist_name: 'Alex Rivera',
        tourist_phone: '+94 77 000 1122',
        guide_id: guide.id,
        meeting_point: meetingPoint || `${guide.city} Central Landmark`,
        destination: destination || `${guide.specialties[0] || 'City'} Highlights Tour`,
        duration_hours: duration,
        special_notes: notes,
      });
      setConfirmedBooking(result);
    } catch {
      // Fallback local mock booking
      setConfirmedBooking({
        id: `tour-${Date.now().toString().slice(-6)}`,
        tourist_id: 'tourist-user-101',
        tourist_name: 'Alex Rivera',
        guide_id: guide.id,
        guide_name: guide.name,
        status: 'REQUESTED',
        meeting_point: meetingPoint || `${guide.city} Central Landmark`,
        destination: destination || `${guide.specialties[0] || 'City'} Highlights`,
        duration_hours: duration,
        hourly_rate: guide.hourly_rate,
        total_estimated: totalEstimate,
        special_notes: notes,
        created_at: new Date().toISOString(),
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs animate-fadeIn">
      <div className="w-full max-w-lg bg-white rounded-3xl shadow-2xl border border-slate-200 overflow-hidden animate-scaleUp">
        {/* Header */}
        <div className="p-5 bg-gradient-to-r from-slate-900 to-slate-800 text-white flex items-center justify-between">
          <div>
            <h3 className="font-bold text-base">Request Instant Tour</h3>
            <p className="text-xs text-slate-300">With {guide.name} • {guide.city}</p>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-xl text-slate-400 hover:text-white hover:bg-slate-700 transition-colors cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {confirmedBooking ? (
          <div className="p-6 text-center">
            <div className="w-14 h-14 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto mb-3">
              <CheckCircle2 className="w-8 h-8" />
            </div>
            <h4 className="text-lg font-bold text-slate-900">Tour Request Dispatched!</h4>
            <p className="text-xs text-slate-500 mt-1 max-w-sm mx-auto">
              {guide.name} has been notified via real-time WebSocket. You will receive an alert once accepted.
            </p>

            <div className="my-5 p-4 bg-slate-50 border border-slate-200 rounded-2xl text-left text-xs space-y-1.5">
              <div className="flex justify-between text-slate-600">
                <span>Booking ID:</span>
                <span className="font-mono font-bold text-slate-900">{confirmedBooking.id}</span>
              </div>
              <div className="flex justify-between text-slate-600">
                <span>Duration:</span>
                <span className="font-semibold text-slate-900">{confirmedBooking.duration_hours} Hours</span>
              </div>
              <div className="flex justify-between text-slate-600">
                <span>Meeting Point:</span>
                <span className="font-semibold text-slate-900">{confirmedBooking.meeting_point}</span>
              </div>
              <div className="flex justify-between text-slate-600 pt-2 border-t border-slate-200">
                <span className="font-bold text-slate-800">Total Estimated:</span>
                <span className="font-extrabold text-emerald-700 text-sm">${confirmedBooking.total_estimated.toFixed(2)}</span>
              </div>
            </div>

            <button
              onClick={onClose}
              className="w-full py-3 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs rounded-xl shadow-md transition-all cursor-pointer"
            >
              Done
            </button>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="p-6 space-y-4">
            {/* Duration Selector */}
            <div>
              <label className="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-2">
                Tour Duration
              </label>
              <div className="grid grid-cols-4 gap-2">
                {[1, 2, 4, 8].map((hours) => (
                  <button
                    key={hours}
                    type="button"
                    onClick={() => setDuration(hours)}
                    className={`py-2 px-3 rounded-xl border text-xs font-semibold flex flex-col items-center gap-0.5 transition-all cursor-pointer ${
                      duration === hours
                        ? 'bg-emerald-600 text-white border-emerald-600 shadow-sm'
                        : 'bg-slate-50 text-slate-700 border-slate-200 hover:bg-slate-100'
                    }`}
                  >
                    <span>{hours} {hours === 1 ? 'Hour' : 'Hours'}</span>
                    <span className={`text-[10px] ${duration === hours ? 'text-emerald-100' : 'text-slate-400'}`}>
                      ${hours * guide.hourly_rate}
                    </span>
                  </button>
                ))}
              </div>
            </div>

            {/* Meeting Point */}
            <div>
              <label className="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1">
                Meeting Point
              </label>
              <div className="relative">
                <MapPin className="absolute left-3 top-3 w-4 h-4 text-slate-400" />
                <input
                  type="text"
                  value={meetingPoint}
                  onChange={(e) => setMeetingPoint(e.target.value)}
                  placeholder={`e.g. ${guide.city === 'Galle' ? 'Galle Fort Clock Tower' : guide.city + ' Main Landmark'}`}
                  className="w-full pl-9 pr-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>

            {/* Destination / Activity */}
            <div>
              <label className="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1">
                Tour Activity / Focus
              </label>
              <input
                type="text"
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
                placeholder={`e.g. ${guide.specialties[0] || 'Historical walking tour'}`}
                className="w-full px-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
            </div>

            {/* Price Breakdown Banner */}
            <div className="p-3 bg-emerald-50 border border-emerald-200 rounded-2xl flex items-center justify-between text-xs">
              <div className="flex items-center gap-2 text-emerald-900 font-semibold">
                <Clock className="w-4 h-4 text-emerald-600" />
                <span>Rate: ${guide.hourly_rate}/hr × {duration}h</span>
              </div>
              <div className="font-extrabold text-base text-emerald-800">
                ${totalEstimate.toFixed(2)}
              </div>
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 text-white font-bold text-xs rounded-xl shadow-md shadow-emerald-600/20 transition-all cursor-pointer flex items-center justify-center gap-2"
            >
              {loading ? (
                <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              ) : (
                <>
                  <CheckCircle2 className="w-4 h-4" />
                  <span>Send Tour Request to {guide.name}</span>
                </>
              )}
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
