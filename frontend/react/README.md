# React Component for go-chatbot

A modern React TypeScript component for integrating go-chatbot into React applications.

## Features

- 🚀 **TypeScript Support** - Full type safety and IntelliSense
- 💬 **Floating Chat Widget** - Non-intrusive chat button with popup window
- 🎨 **Customizable** - Extensive styling and behavior options
- 📱 **Responsive** - Works on desktop and mobile devices
- ⚡ **Performance** - Optimized with React hooks and efficient rendering
- 🧪 **Well Tested** - Comprehensive test suite with high coverage
- ♿ **Accessible** - ARIA-compliant and keyboard navigation support

## Installation

```bash
npm install @go-chatbot/react
# or
yarn add @go-chatbot/react
```

## Basic Usage

```tsx
import React from 'react';
import { GoChatbot } from '@go-chatbot/react';

function App() {
  return (
    <div>
      <h1>My App</h1>
      <GoChatbot />
    </div>
  );
}

export default App;
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `apiEndpoint` | `string` | `'/chat/'` | API endpoint for chatbot requests |
| `className` | `string` | `''` | Custom CSS classes |
| `placeholder` | `string` | `'Type your message...'` | Input placeholder text |
| `maxHeight` | `string` | `'400px'` | Maximum height of chat window |
| `style` | `React.CSSProperties` | `{}` | Custom inline styles |
| `onMessageSent` | `(message: string) => void` | `undefined` | Callback when message is sent |
| `onResponseReceived` | `(response: string) => void` | `undefined` | Callback when response is received |
| `onError` | `(error: string) => void` | `undefined` | Callback when error occurs |
| `initialMessages` | `ChatMessage[]` | `[]` | Initial messages to display |
| `showTypingIndicator` | `boolean` | `true` | Show typing indicator during loading |
| `header` | `React.ReactNode` | `'Chat Support'` | Custom header content |
| `disabled` | `boolean` | `false` | Disable the chat component |

## Advanced Usage

### Custom Styling

```tsx
import { GoChatbot } from '@go-chatbot/react';

function App() {
  return (
    <GoChatbot
      className="my-custom-chat"
      style={{
        fontFamily: 'Arial, sans-serif',
        borderRadius: '8px',
      }}
      maxHeight="500px"
    />
  );
}
```

### Event Handlers

```tsx
import { GoChatbot } from '@go-chatbot/react';

function App() {
  const handleMessageSent = (message: string) => {
    console.log('User sent:', message);
    // Track analytics, etc.
  };

  const handleResponseReceived = (response: string) => {
    console.log('Bot responded:', response);
    // Log response, etc.
  };

  const handleError = (error: string) => {
    console.error('Chat error:', error);
    // Handle error, show notification, etc.
  };

  return (
    <GoChatbot
      onMessageSent={handleMessageSent}
      onResponseReceived={handleResponseReceived}
      onError={handleError}
    />
  );
}
```

### Initial Messages

```tsx
import { GoChatbot, ChatMessage } from '@go-chatbot/react';

function App() {
  const initialMessages: ChatMessage[] = [
    {
      id: '1',
      text: 'Welcome! How can I help you today?',
      sender: 'bot',
      timestamp: new Date(),
    },
  ];

  return (
    <GoChatbot initialMessages={initialMessages} />
  );
}
```

### Custom Header

```tsx
import { GoChatbot } from '@go-chatbot/react';

function App() {
  const customHeader = (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
      <img src="/logo.png" alt="Logo" width="20" height="20" />
      <span>AI Assistant</span>
    </div>
  );

  return (
    <GoChatbot header={customHeader} />
  );
}
```

## Types

### ChatMessage

```typescript
interface ChatMessage {
  id: string;
  text: string;
  sender: 'user' | 'bot';
  timestamp: Date;
}
```

### GoChatbotProps

```typescript
interface GoChatbotProps {
  apiEndpoint?: string;
  className?: string;
  placeholder?: string;
  maxHeight?: string;
  style?: React.CSSProperties;
  onMessageSent?: (message: string) => void;
  onResponseReceived?: (response: string) => void;
  onError?: (error: string) => void;
  initialMessages?: ChatMessage[];
  showTypingIndicator?: boolean;
  header?: React.ReactNode;
  disabled?: boolean;
}
```

## Backend Integration

The component expects a REST API endpoint that accepts POST requests with the following format:

**Request:**
```json
{
  "message": "User's message"
}
```

**Response (Success):**
```json
{
  "success": true,
  "response": "Bot's response"
}
```

**Response (Error):**
```json
{
  "success": false,
  "error": "Error message"
}
```

### go-chatbot Backend

This component is designed to work with the [go-chatbot](https://go.rumenx.com/chatbot) backend. Example server setup:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go.rumenx.com/chatbot"
    "go.rumenx.com/chatbot/adapters"
    "go.rumenx.com/chatbot/config"
)

func main() {
    cfg := config.Default()
    chatbot, _ := gochatbot.New(cfg)

    r := gin.Default()
    adapter := adapters.NewGinAdapter(chatbot)
    adapter.SetupRoutes(r)

    r.Run(":8080")
}
```

## Development

### Building

```bash
npm run build
```

### Testing

```bash
# Run tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm test -- --coverage
```

### Linting

```bash
# Check for linting errors
npm run lint

# Fix linting errors
npm run lint:fix
```

## Browser Support

- Chrome >= 60
- Firefox >= 60
- Safari >= 12
- Edge >= 79

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT © [Rumen Damyanov](https://github.com/RumenDamyanov)
