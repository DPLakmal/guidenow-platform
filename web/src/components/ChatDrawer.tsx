'use client';

import React, { useState, useEffect, useRef } from 'react';
import { Guide, ChatMessage } from '@/lib/types';
import { X, Send, ShieldCheck, Wifi, WifiOff } from 'lucide-react';

interface ChatDrawerProps {
  guide: Guide | null;
  isOpen: boolean;
  onClose: () => void;
}

export default function ChatDrawer({ guide, isOpen, onClose }: ChatDrawerProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [connected, setConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);
  const scrollRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!isOpen || !guide) {
      if (socketRef.current) {
        socketRef.current.close();
      }
      return;
    }

    const roomId = `room-${guide.id}`;
    const wsUrl = `ws://localhost:8080/ws/chat?room=${roomId}&userId=tourist-101&name=Traveler&role=tourist`;

    // Initialize messages with friendly greeting
    setMessages([
      {
        id: 'msg-init',
        room_id: roomId,
        sender_id: guide.id,
        sender_name: guide.name,
        sender_role: 'guide',
        content: `Ayubowan! 🙏 I'm currently available in ${guide.city}. How can I help you plan your tour?`,
        timestamp: new Date().toISOString(),
      },
    ]);

    try {
      const socket = new WebSocket(wsUrl);

      socket.onopen = () => {
        setConnected(true);
      };

      socket.onmessage = (event) => {
        try {
          const incoming: ChatMessage = JSON.parse(event.data);
          setMessages((prev) => [...prev, incoming]);
        } catch {
          // Non-JSON ping/pong or system notice
        }
      };

      socket.onclose = () => {
        setConnected(false);
      };

      socket.onerror = () => {
        setConnected(false);
      };

      socketRef.current = socket;
    } catch {
      setConnected(false);
    }

    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, [isOpen, guide]);

  useEffect(() => {
    scrollRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = () => {
    if (!input.trim() || !guide) return;

    const newMsg: ChatMessage = {
      id: `client-${Date.now()}`,
      room_id: `room-${guide.id}`,
      sender_id: 'tourist-101',
      sender_name: 'You (Traveler)',
      sender_role: 'tourist',
      content: input,
      timestamp: new Date().toISOString(),
    };

    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ content: input }));
    } else {
      // Local echo fallback if websocket is not connected
      setMessages((prev) => [...prev, newMsg]);

      // Mock guide response after 1 second
      setTimeout(() => {
        setMessages((prev) => [
          ...prev,
          {
            id: `reply-${Date.now()}`,
            room_id: `room-${guide.id}`,
            sender_id: guide.id,
            sender_name: guide.name,
            sender_role: 'guide',
            content: `Understood! I'll be ready at the meeting point. Feel free to send me your tour request!`,
            timestamp: new Date().toISOString(),
          },
        ]);
      }, 1000);
    }

    setInput('');
  };

  if (!isOpen || !guide) return null;

  return (
    <div className="fixed inset-0 z-50 flex justify-end bg-slate-900/40 backdrop-blur-xs transition-opacity animate-fadeIn">
      <div className="w-full max-w-md bg-white h-full shadow-2xl flex flex-col border-l border-slate-200 animate-slideLeft">
        {/* Header */}
        <div className="p-4 bg-slate-900 text-white flex items-center justify-between border-b border-slate-800">
          <div className="flex items-center gap-3">
            <div className="relative w-10 h-10 rounded-full overflow-hidden border border-white/20">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src={guide.avatar} alt={guide.name} className="w-full h-full object-cover" />
              <div className="absolute top-0 right-0 w-3 h-3 bg-emerald-500 rounded-full border-2 border-slate-900" />
            </div>
            <div>
              <div className="font-bold text-sm flex items-center gap-1.5">
                <span>{guide.name}</span>
                <ShieldCheck className="w-3.5 h-3.5 text-emerald-400" />
              </div>
              <div className="text-[11px] text-slate-400 flex items-center gap-1.5">
                {connected ? (
                  <span className="flex items-center gap-1 text-emerald-400">
                    <Wifi className="w-3 h-3" /> Live WebSocket
                  </span>
                ) : (
                  <span className="flex items-center gap-1 text-amber-400">
                    <WifiOff className="w-3 h-3" /> Standalone Echo
                  </span>
                )}
                <span>• ${guide.hourly_rate}/hr</span>
              </div>
            </div>
          </div>

          <button
            onClick={onClose}
            className="p-2 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 transition-colors cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Message Area */}
        <div className="flex-1 overflow-y-auto p-4 space-y-3 bg-slate-50">
          {messages.map((m) => {
            const isMe = m.sender_role === 'tourist';
            return (
              <div
                key={m.id}
                className={`flex flex-col ${isMe ? 'items-end' : 'items-start'}`}
              >
                <div
                  className={`max-w-[80%] p-3 rounded-2xl text-xs leading-relaxed shadow-xs ${
                    isMe
                      ? 'bg-emerald-600 text-white rounded-br-none'
                      : 'bg-white text-slate-800 border border-slate-200 rounded-bl-none'
                  }`}
                >
                  <div className={`font-semibold text-[10px] mb-1 ${isMe ? 'text-emerald-200' : 'text-slate-500'}`}>
                    {m.sender_name}
                  </div>
                  <div>{m.content}</div>
                </div>
                <span className="text-[10px] text-slate-400 mt-1 px-1">
                  {new Date(m.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                </span>
              </div>
            );
          })}
          <div ref={scrollRef} />
        </div>

        {/* Input Bar */}
        <div className="p-4 bg-white border-t border-slate-200">
          <form
            onSubmit={(e) => {
              e.preventDefault();
              handleSend();
            }}
            className="flex items-center gap-2"
          >
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder={`Message ${guide.name}...`}
              className="flex-1 px-4 py-2.5 bg-slate-100 border border-slate-200 rounded-xl text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:bg-white transition-all"
            />
            <button
              type="submit"
              disabled={!input.trim()}
              className="p-2.5 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 text-white rounded-xl shadow-sm transition-all cursor-pointer"
            >
              <Send className="w-4 h-4" />
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
