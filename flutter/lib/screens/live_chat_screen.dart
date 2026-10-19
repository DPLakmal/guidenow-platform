import 'package:flutter/material.dart';
import '../models/guide.dart';
import '../models/message.dart';

class LiveChatScreen extends StatefulWidget {
  final Guide guide;

  const LiveChatScreen({super.key, required this.guide});

  @override
  State<LiveChatScreen> createState() => _LiveChatScreenState();
}

class _LiveChatScreenState extends State<LiveChatScreen> {
  final TextEditingController _textController = TextEditingController();
  final List<ChatMessage> _messages = [];

  @override
  void initState() {
    super.initState();
    _messages.add(
      ChatMessage(
        id: 'msg-1',
        roomId: 'room-${widget.guide.id}',
        senderId: widget.guide.id,
        senderName: widget.guide.name,
        senderRole: 'guide',
        content: 'Ayubowan! 🙏 Welcome to Sri Lanka. How can I assist you with your tour in ${widget.guide.city}?',
        timestamp: DateTime.now(),
      ),
    );
  }

  void _sendMessage() {
    final text = _textController.text.trim();
    if (text.isEmpty) return;

    setState(() {
      _messages.add(
        ChatMessage(
          id: 'client-${DateTime.now().millisecondsSinceEpoch}',
          roomId: 'room-${widget.guide.id}',
          senderId: 'tourist-user-1',
          senderName: 'You',
          senderRole: 'tourist',
          content: text,
          timestamp: DateTime.now(),
        ),
      );
      _textController.clear();
    });

    // Mock guide reply
    Future.delayed(const Duration(milliseconds: 1000), () {
      if (mounted) {
        setState(() {
          _messages.add(
            ChatMessage(
              id: 'reply-${DateTime.now().millisecondsSinceEpoch}',
              roomId: 'room-${widget.guide.id}',
              senderId: widget.guide.id,
              senderName: widget.guide.name,
              senderRole: 'guide',
              content: 'Got your message! I am ready to guide you. Tap "Request Tour" when ready.',
              timestamp: DateTime.now(),
            ),
          );
        });
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        backgroundColor: const Color(0xFF0F172A),
        foregroundColor: Colors.white,
        elevation: 0,
        title: Row(
          children: [
            CircleAvatar(
              radius: 16,
              backgroundColor: const Color(0xFF1E293B),
              child: Text(widget.guide.name[0], style: const TextStyle(color: Colors.white, fontSize: 12)),
            ),
            const SizedBox(width: 8),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(widget.guide.name, style: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold)),
                Text(
                  '${widget.guide.city} • \$${widget.guide.hourlyRate.toInt()}/hr',
                  style: const TextStyle(fontSize: 10, color: Color(0xFF10B981)),
                ),
              ],
            ),
          ],
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _messages.length,
              itemBuilder: (context, index) {
                final msg = _messages[index];
                final isMe = msg.senderRole == 'tourist';

                return Align(
                  alignment: isMe ? Alignment.centerRight : Alignment.centerLeft,
                  child: Container(
                    margin: const EdgeInsets.only(bottom: 8),
                    padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                    constraints: BoxConstraints(maxWidth: MediaQuery.of(context).size.width * 0.75),
                    decoration: BoxDecoration(
                      color: isMe ? const Color(0xFF10B981) : Colors.white,
                      borderRadius: BorderRadius.only(
                        topLeft: const Radius.circular(16),
                        topRight: const Radius.circular(16),
                        bottomLeft: isMe ? const Radius.circular(16) : Radius.zero,
                        bottomRight: isMe ? Radius.zero : const Radius.circular(16),
                      ),
                      border: isMe ? null : Border.all(color: const Color(0xFFE2E8F0)),
                      boxShadow: const [
                        BoxShadow(color: Colors.black12, blurRadius: 2, offset: Offset(0, 1)),
                      ],
                    ),
                    child: Column(
                      crossAxisAlignment: isMe ? CrossAxisAlignment.end : CrossAxisAlignment.start,
                      children: [
                        Text(
                          msg.content,
                          style: TextStyle(
                            fontSize: 13,
                            color: isMe ? Colors.white : const Color(0xFF0F172A),
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          '${msg.timestamp.hour.toString().padLeft(2, '0')}:${msg.timestamp.minute.toString().padLeft(2, '0')}',
                          style: TextStyle(
                            fontSize: 9,
                            color: isMe ? Colors.white70 : const Color(0xFF94A3B8),
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),
          ),
          Container(
            padding: const EdgeInsets.all(12),
            decoration: const BoxDecoration(
              color: Colors.white,
              border: Border(top: BorderSide(color: Color(0xFFE2E8F0))),
            ),
            child: SafeArea(
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _textController,
                      decoration: InputDecoration(
                        hintText: 'Type your message...',
                        hintStyle: const TextStyle(fontSize: 13, color: Color(0xFF94A3B8)),
                        filled: true,
                        fillColor: const Color(0xFFF1F5F9),
                        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(24),
                          borderSide: BorderSide.none,
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  IconButton(
                    onPressed: _sendMessage,
                    icon: const Icon(Icons.send_rounded),
                    color: const Color(0xFF10B981),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
