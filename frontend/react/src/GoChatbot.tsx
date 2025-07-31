import React, { useState, useRef, useEffect } from 'react';

export interface ChatMessage {
  id: string;
  text: string;
  sender: 'user' | 'bot';
  timestamp: Date;
}

export interface GoChatbotProps {
  /** API endpoint for the chatbot */
  apiEndpoint?: string;
  /** Custom CSS classes */
  className?: string;
  /** Placeholder text for input */
  placeholder?: string;
  /** Maximum height of chat container */
  maxHeight?: string;
  /** Custom styles */
  style?: React.CSSProperties;
  /** Callback when message is sent */
  onMessageSent?: (message: string) => void;
  /** Callback when response is received */
  onResponseReceived?: (response: string) => void;
  /** Callback when error occurs */
  onError?: (error: string) => void;
  /** Initial messages */
  initialMessages?: ChatMessage[];
  /** Whether to show typing indicator */
  showTypingIndicator?: boolean;
  /** Custom header content */
  header?: React.ReactNode;
  /** Whether chat is disabled */
  disabled?: boolean;
}

const GoChatbot: React.FC<GoChatbotProps> = ({
  apiEndpoint = '/chat/',
  className = '',
  placeholder = 'Type your message...',
  maxHeight = '400px',
  style = {},
  onMessageSent,
  onResponseReceived,
  onError,
  initialMessages = [],
  showTypingIndicator = true,
  header,
  disabled = false,
}) => {
  const [messages, setMessages] = useState<ChatMessage[]>(initialMessages);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, isLoading]);

  useEffect(() => {
    if (isOpen && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isOpen]);

  const generateId = () => `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

  const sendMessage = async () => {
    if (!inputValue.trim() || isLoading || disabled) return;

    const userMessage: ChatMessage = {
      id: generateId(),
      text: inputValue.trim(),
      sender: 'user',
      timestamp: new Date(),
    };

    setMessages(prev => [...prev, userMessage]);
    setInputValue('');
    setIsLoading(true);

    onMessageSent?.(userMessage.text);

    try {
      const response = await fetch(apiEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ message: userMessage.text }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.error || 'Unknown error occurred');
      }

      const botMessage: ChatMessage = {
        id: generateId(),
        text: data.response,
        sender: 'bot',
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, botMessage]);
      onResponseReceived?.(data.response);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to send message';
      const botMessage: ChatMessage = {
        id: generateId(),
        text: `Error: ${errorMessage}`,
        sender: 'bot',
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, botMessage]);
      onError?.(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const chatButton = (
    <button
      onClick={() => setIsOpen(!isOpen)}
      style={{
        position: 'fixed',
        bottom: '20px',
        right: '20px',
        width: '60px',
        height: '60px',
        borderRadius: '50%',
        backgroundColor: '#007bff',
        color: 'white',
        border: 'none',
        cursor: 'pointer',
        boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
        fontSize: '24px',
        zIndex: 1000,
        transition: 'all 0.3s ease',
      }}
      disabled={disabled}
    >
      {isOpen ? '×' : '💬'}
    </button>
  );

  const chatWindow = (
    <div
      className={`go-chatbot ${className}`}
      style={{
        position: 'fixed',
        bottom: '90px',
        right: '20px',
        width: '350px',
        height: maxHeight,
        backgroundColor: 'white',
        borderRadius: '12px',
        boxShadow: '0 8px 32px rgba(0,0,0,0.12)',
        display: 'flex',
        flexDirection: 'column',
        zIndex: 999,
        overflow: 'hidden',
        ...style,
      }}
    >
      {/* Header */}
      <div
        style={{
          padding: '16px',
          backgroundColor: '#007bff',
          color: 'white',
          borderRadius: '12px 12px 0 0',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        {header || <span style={{ fontWeight: 'bold' }}>Chat Support</span>}
        <button
          onClick={() => setIsOpen(false)}
          style={{
            background: 'none',
            border: 'none',
            color: 'white',
            cursor: 'pointer',
            fontSize: '18px',
            padding: '0',
            width: '24px',
            height: '24px',
          }}
        >
          ×
        </button>
      </div>

      {/* Messages */}
      <div
        style={{
          flex: 1,
          overflowY: 'auto',
          padding: '16px',
          display: 'flex',
          flexDirection: 'column',
          gap: '12px',
        }}
      >
        {messages.map((message) => (
          <div
            key={message.id}
            style={{
              display: 'flex',
              justifyContent: message.sender === 'user' ? 'flex-end' : 'flex-start',
            }}
          >
            <div
              style={{
                maxWidth: '80%',
                padding: '8px 12px',
                borderRadius: '18px',
                backgroundColor: message.sender === 'user' ? '#007bff' : '#f1f3f5',
                color: message.sender === 'user' ? 'white' : '#333',
                fontSize: '14px',
                lineHeight: '1.4',
                wordWrap: 'break-word',
              }}
            >
              <div>{message.text}</div>
              <div
                style={{
                  fontSize: '10px',
                  opacity: 0.7,
                  marginTop: '4px',
                  textAlign: 'right',
                }}
              >
                {formatTime(message.timestamp)}
              </div>
            </div>
          </div>
        ))}

        {isLoading && showTypingIndicator && (
          <div style={{ display: 'flex', justifyContent: 'flex-start' }}>
            <div
              style={{
                padding: '8px 12px',
                borderRadius: '18px',
                backgroundColor: '#f1f3f5',
                color: '#333',
                fontSize: '14px',
              }}
            >
              <div style={{ display: 'flex', gap: '4px', alignItems: 'center' }}>
                <span>Bot is typing</span>
                <div style={{ display: 'flex', gap: '2px' }}>
                  <div
                    style={{
                      width: '4px',
                      height: '4px',
                      borderRadius: '50%',
                      backgroundColor: '#999',
                      animation: 'pulse 1.5s infinite',
                    }}
                  />
                  <div
                    style={{
                      width: '4px',
                      height: '4px',
                      borderRadius: '50%',
                      backgroundColor: '#999',
                      animation: 'pulse 1.5s infinite 0.5s',
                    }}
                  />
                  <div
                    style={{
                      width: '4px',
                      height: '4px',
                      borderRadius: '50%',
                      backgroundColor: '#999',
                      animation: 'pulse 1.5s infinite 1s',
                    }}
                  />
                </div>
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <div
        style={{
          padding: '16px',
          borderTop: '1px solid #e9ecef',
          display: 'flex',
          gap: '8px',
        }}
      >
        <input
          ref={inputRef}
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder={placeholder}
          disabled={isLoading || disabled}
          style={{
            flex: 1,
            padding: '8px 12px',
            border: '1px solid #dee2e6',
            borderRadius: '20px',
            outline: 'none',
            fontSize: '14px',
          }}
        />
        <button
          onClick={sendMessage}
          disabled={!inputValue.trim() || isLoading || disabled}
          style={{
            padding: '8px 16px',
            backgroundColor: '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '20px',
            cursor: isLoading || disabled ? 'not-allowed' : 'pointer',
            fontSize: '14px',
            opacity: !inputValue.trim() || isLoading || disabled ? 0.6 : 1,
          }}
        >
          {isLoading ? '...' : 'Send'}
        </button>
      </div>
    </div>
  );

  return (
    <>
      <style>
        {`
          @keyframes pulse {
            0%, 100% { opacity: 0.4; }
            50% { opacity: 1; }
          }
        `}
      </style>
      {chatButton}
      {isOpen && chatWindow}
    </>
  );
};

export default GoChatbot;
