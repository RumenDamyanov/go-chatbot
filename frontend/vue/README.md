# Vue Component for go-chatbot

A modern Vue 3 TypeScript component for integrating go-chatbot into Vue.js applications.

## Features

- 🚀 **Vue 3 + TypeScript** - Built with the latest Vue 3 Composition API and full TypeScript support
- 💬 **Floating Chat Widget** - Non-intrusive chat button with popup window
- 🎨 **Customizable** - Extensive styling options via props and slots
- 📱 **Responsive** - Works seamlessly on desktop and mobile devices
- ⚡ **Performance** - Optimized with Vue 3 reactivity and efficient DOM updates
- 🧪 **Well Tested** - Comprehensive test suite with Vitest
- ♿ **Accessible** - ARIA-compliant and keyboard navigation support
- 🎭 **Flexible** - Supports custom headers via slots

## Installation

```bash
npm install @go-chatbot/vue
# or
yarn add @go-chatbot/vue
# or
pnpm add @go-chatbot/vue
```

## Basic Usage

```vue
<template>
  <div>
    <h1>My Vue App</h1>
    <GoChatbot />
  </div>
</template>

<script setup>
import { GoChatbot } from '@go-chatbot/vue'
</script>
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `apiEndpoint` | `string` | `'/chat/'` | API endpoint for chatbot requests |
| `className` | `string` | `''` | Custom CSS classes |
| `placeholder` | `string` | `'Type your message...'` | Input placeholder text |
| `maxHeight` | `string` | `'400px'` | Maximum height of chat window |
| `style` | `Record<string, any>` | `{}` | Custom inline styles |
| `initialMessages` | `ChatMessage[]` | `[]` | Initial messages to display |
| `showTypingIndicator` | `boolean` | `true` | Show typing indicator during loading |
| `disabled` | `boolean` | `false` | Disable the chat component |

## Events

| Event | Payload | Description |
|-------|---------|-------------|
| `message-sent` | `string` | Emitted when user sends a message |
| `response-received` | `string` | Emitted when bot response is received |
| `error` | `string` | Emitted when an error occurs |

## Slots

| Slot | Description |
|------|-------------|
| `header` | Custom header content for the chat window |

## Advanced Usage

### Custom Styling

```vue
<template>
  <GoChatbot
    class-name="my-custom-chat"
    :style="{ fontFamily: 'Arial, sans-serif' }"
    max-height="500px"
  />
</template>
```

### Event Handlers

```vue
<template>
  <GoChatbot
    @message-sent="handleMessageSent"
    @response-received="handleResponseReceived"
    @error="handleError"
  />
</template>

<script setup>
const handleMessageSent = (message: string) => {
  console.log('User sent:', message)
  // Track analytics, etc.
}

const handleResponseReceived = (response: string) => {
  console.log('Bot responded:', response)
  // Log response, etc.
}

const handleError = (error: string) => {
  console.error('Chat error:', error)
  // Handle error, show notification, etc.
}
</script>
```

### Initial Messages

```vue
<template>
  <GoChatbot :initial-messages="initialMessages" />
</template>

<script setup>
import type { ChatMessage } from '@go-chatbot/vue'

const initialMessages: ChatMessage[] = [
  {
    id: '1',
    text: 'Welcome! How can I help you today?',
    sender: 'bot',
    timestamp: new Date(),
  },
]
</script>
```

### Custom Header

```vue
<template>
  <GoChatbot>
    <template #header>
      <div style="display: flex; align-items: center; gap: 8px;">
        <img src="/logo.png" alt="Logo" width="20" height="20" />
        <span>AI Assistant</span>
      </div>
    </template>
  </GoChatbot>
</template>
```

### Options API (Vue 2 style)

```vue
<template>
  <GoChatbot
    :api-endpoint="apiEndpoint"
    :initial-messages="initialMessages"
    @message-sent="handleMessageSent"
    @response-received="handleResponseReceived"
    @error="handleError"
  />
</template>

<script>
import { GoChatbot } from '@go-chatbot/vue'

export default {
  components: {
    GoChatbot
  },
  data() {
    return {
      apiEndpoint: '/api/chat',
      initialMessages: [
        {
          id: '1',
          text: 'Hello! How can I help?',
          sender: 'bot',
          timestamp: new Date()
        }
      ]
    }
  },
  methods: {
    handleMessageSent(message) {
      console.log('Message sent:', message)
    },
    handleResponseReceived(response) {
      console.log('Response received:', response)
    },
    handleError(error) {
      console.error('Error:', error)
    }
  }
}
</script>
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

# Run tests with UI
npm run test:ui

# Run tests with coverage
npm run test:coverage
```

### Linting

```bash
# Check for linting errors
npm run lint

# Fix linting errors
npm run lint:fix
```

### Type Checking

```bash
npm run type-check
```

## Browser Support

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## Vue Compatibility

- Vue 3.3+ (Composition API)
- Vue 2.7+ (with @vue/composition-api)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT © [Rumen Damyanov](https://github.com/RumenDamyanov)
