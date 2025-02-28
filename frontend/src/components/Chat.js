import React, { useState, useEffect, useRef } from 'react';
import {
  Box,
  TextField,
  Button,
  Paper,
  Typography,
  Container,
} from '@mui/material';

const Chat = () => {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [socket, setSocket] = useState(null);
  const messagesEndRef = useRef(null);
  const [userId, setUserId] = useState('');
  const [receiverId, setReceiverId] = useState('');
  const [isConnected, setIsConnected] = useState(false);

  const handleConnect = (e) => {
    e.preventDefault();
    if (!userId.trim() || !receiverId.trim()) return;

    const ws = new WebSocket(`ws://localhost:8080/ws?userId=${userId}`);

    ws.onopen = () => {
      console.log('Connected to WebSocket');
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages((prev) => [...prev, message]);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    setSocket(ws);
  };

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (!newMessage.trim()) return;

    const messageData = {
      type: 'message',
      content: newMessage,
      senderId: userId,
      receiverId: receiverId,
      timestamp: new Date().toISOString(),
    };

    if (socket) {
      socket.send(JSON.stringify(messageData));
      setNewMessage('');
    }
  };

  if (!isConnected) {
    return (
      <Container maxWidth="sm">
        <Paper elevation={3} sx={{ p: 3, mt: 2 }}>
          <Typography variant="h6" sx={{ mb: 2 }}>Connect to Chat</Typography>
          <Box component="form" onSubmit={handleConnect} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              fullWidth
              label="Your User ID"
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              required
            />
            <TextField
              fullWidth
              label="Recipient User ID"
              value={receiverId}
              onChange={(e) => setReceiverId(e.target.value)}
              required
            />
            <Button type="submit" variant="contained">
              Connect
            </Button>
          </Box>
        </Paper>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm">
      <Paper elevation={3} sx={{ height: '80vh', mt: 2, display: 'flex', flexDirection: 'column' }}>
        {/* Chat Header */}
        <Box sx={{ p: 2, backgroundColor: '#f5f5f5', borderBottom: '1px solid #e0e0e0' }}>
          <Typography variant="h6">Chat Room</Typography>
        </Box>

        {/* Chat Messages */}
        <Box sx={{ flex: 1, overflow: 'auto', p: 2 }}>
          {messages.map((message, index) => (
            <Box
              key={index}
              sx={{
                display: 'flex',
                justifyContent: message.senderId === userId ? 'flex-end' : 'flex-start',
                mb: 1,
              }}
            >
              <Paper
                sx={{
                  p: 1,
                  backgroundColor: message.senderId === userId ? '#e3f2fd' : '#f5f5f5',
                  maxWidth: '70%',
                }}
              >
                <Typography variant="body1">{message.content}</Typography>
              </Paper>
            </Box>
          ))}
          <div ref={messagesEndRef} />
        </Box>

        {/* Chat Input */}
        <Box
          component="form"
          onSubmit={handleSendMessage}
          sx={{ p: 2, borderTop: '1px solid #e0e0e0', display: 'flex', gap: 1 }}
        >
          <TextField
            fullWidth
            size="small"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type a message..."
          />
          <Button type="submit" variant="contained">
            Send
          </Button>
        </Box>
      </Paper>
    </Container>
  );
};

export default Chat;