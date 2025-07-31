# Go Chatbot - Frontend Components

This directory contains frontend integration components for the Go Chatbot system, providing seamless chat functionality for popular JavaScript frameworks and vanilla JavaScript.

## Available Components

### 🔵 React Component (`/react/`)

Modern React TypeScript component with hooks, comprehensive props interface, and full TypeScript support.

- **Technology**: React 18+, TypeScript, modern hooks
- **Testing**: Jest + React Testing Library
- **Build**: Vite, ESLint, TypeScript compiler
- **Features**: Floating chat widget, customizable styling, event handling

[📖 React Documentation](./react/README.md)

### 🟢 Vue Component (`/vue/`)

Vue 3 component using Composition API with full TypeScript integration and reactive state management.

- **Technology**: Vue 3, Composition API, TypeScript
- **Testing**: Vitest + Vue Test Utils
- **Build**: Vite, ESLint, TypeScript support
- **Features**: Scoped styling, slot-based customization, reactive props

[📖 Vue Documentation](./vue/README.md)

### 🔴 Angular Component (`/angular/`)

Angular 17+ component with standalone support, RxJS integration, and comprehensive module system.

- **Technology**: Angular 17+, RxJS, TypeScript
- **Testing**: Jest + Angular Testing utilities
- **Build**: ng-packagr, Angular CLI
- **Features**: Standalone component, reactive forms, HTTP client integration

[📖 Angular Documentation](./angular/README.md)

### ⚡ Vanilla JavaScript (`/vanilla/`)

Framework-agnostic vanilla JavaScript component that works anywhere without dependencies.

- **Technology**: Pure JavaScript ES6+, TypeScript definitions
- **Testing**: Jest with jsdom
- **Build**: Terser for minification
- **Features**: Zero dependencies, 15KB minified, auto-initialization

[📖 Vanilla JS Documentation](./vanilla/README.md)

## Quick Start

### React

```bash
npm install @go-chatbot/react
```

```jsx
import { GoChatbot } from '@go-chatbot/react';

function App() {
  return (
    <div>
      <h1>My App</h1>
      <GoChatbot
        apiEndpoint="/api/chat"
        onMessageSent={(msg) => console.log('Sent:', msg)}
      />
    </div>
  );
}
```

### Vue

```bash
npm install @go-chatbot/vue
```

```vue
<template>
  <div>
    <h1>My App</h1>
    <GoChatbot
      api-endpoint="/api/chat"
      @message-sent="handleMessage"
    />
  </div>
</template>

<script setup>
import { GoChatbot } from '@go-chatbot/vue';

const handleMessage = (msg) => console.log('Sent:', msg);
</script>
```

### Angular

```bash
npm install @go-chatbot/angular
```

```typescript
// app.module.ts
import { GoChatbotModule } from '@go-chatbot/angular';

@NgModule({
  imports: [GoChatbotModule]
})
export class AppModule { }
```

```html
<!-- app.component.html -->
<h1>My App</h1>
<go-chatbot
  apiEndpoint="/api/chat"
  (messageSent)="handleMessage($event)">
</go-chatbot>
```

### Vanilla JavaScript

```html
<script src="https://unpkg.com/@go-chatbot/vanilla@latest/dist/go-chatbot.min.js"></script>
<script>
  const chatbot = new GoChatbot({
    apiEndpoint: '/api/chat'
  });
</script>
```

## Common Features

All components share these core features:

### 🎨 **Customizable Styling**
- CSS custom properties
- Theme support (light/dark)
- Custom CSS classes
- Inline styles
- Responsive design

### 📱 **Responsive Design**
- Mobile-friendly interface
- Touch-optimized interactions
- Adaptive layouts
- Screen reader support

### ⚡ **Performance Optimized**
- Lazy loading
- Efficient DOM updates
- Minimal bundle size
- Tree-shakeable exports

### 🔧 **Developer Experience**
- Full TypeScript support
- Comprehensive documentation
- Example projects
- Testing utilities

### 🎯 **Accessibility**
- ARIA labels and roles
- Keyboard navigation
- Screen reader support
- High contrast compatibility

## API Integration

All components expect your Go backend to implement this simple HTTP API:

### Endpoint: `POST /api/chat`

**Request:**
```json
{
  "message": "User's message text"
}
```

**Success Response:**
```json
{
  "success": true,
  "response": "Bot's response text"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Error message"
}
```

### Go Backend Example

```go
package main

import (
    "encoding/json"
    "net/http"
)

type ChatRequest struct {
    Message string `json:"message"`
}

type ChatResponse struct {
    Success  bool   `json:"success"`
    Response string `json:"response,omitempty"`
    Error    string `json:"error,omitempty"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ChatRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid JSON format")
        return
    }

    // Process the message with your AI/chatbot logic
    response := processMessage(req.Message)

    sendSuccess(w, response)
}

func sendSuccess(w http.ResponseWriter, message string) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ChatResponse{
        Success:  true,
        Response: message,
    })
}

func sendError(w http.ResponseWriter, error string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(ChatResponse{
        Success: false,
        Error:   error,
    })
}

func processMessage(message string) string {
    // Your chatbot logic here
    return "I received your message: " + message
}

func main() {
    http.HandleFunc("/api/chat", chatHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Configuration Options

All components support these common configuration options:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `apiEndpoint` | `string` | `'/chat/'` | Backend API endpoint |
| `placeholder` | `string` | `'Type your message...'` | Input placeholder |
| `maxHeight` | `string` | `'400px'` | Max height of messages |
| `className` | `string` | `''` | Custom CSS class |
| `style` | `object` | `{}` | Custom inline styles |
| `initialMessages` | `array` | `[]` | Pre-populated messages |
| `showTypingIndicator` | `boolean` | `true` | Show typing animation |
| `disabled` | `boolean` | `false` | Disable component |

## Events

All components emit these events:

| Event | Payload | Description |
|-------|---------|-------------|
| `messageSent` | `string` | User sent a message |
| `responseReceived` | `string` | Bot response received |
| `error` | `string` | Error occurred |
| `opened` | `void` | Chat window opened |
| `closed` | `void` | Chat window closed |

## Styling Examples

### CSS Custom Properties

```css
.go-chatbot-container {
  --primary-color: #your-brand-color;
  --text-color: #333333;
  --background-color: #ffffff;
  --border-radius: 12px;
  --shadow: 0 8px 32px rgba(0,0,0,0.12);
}
```

### Dark Theme

```css
.go-chatbot-container.dark-theme {
  --primary-color: #4a5568;
  --text-color: #ffffff;
  --background-color: #2d3748;
  --message-bg-user: #3182ce;
  --message-bg-bot: #4a5568;
}
```

### Custom Positioning

```css
.go-chatbot-container {
  /* Bottom left */
  bottom: 20px;
  left: 20px;

  /* Top right */
  top: 20px;
  right: 20px;
}
```

## Browser Support

- **Modern Browsers**: Chrome 88+, Firefox 78+, Safari 14+, Edge 88+
- **Mobile**: iOS Safari 14+, Android Chrome 88+
- **Framework Compatibility**:
  - React 16.8+ (hooks support)
  - Vue 3.0+
  - Angular 15+
  - Vanilla: Any modern browser

## Testing

Each component includes comprehensive test suites:

```bash
# Run all frontend tests
npm run test:all

# Run specific component tests
cd react && npm test
cd vue && npm test
cd angular && npm test
cd vanilla && npm test
```

Test coverage includes:
- Component rendering
- User interactions
- API communication
- Error handling
- Accessibility
- Event emission

## Development

### Prerequisites

- Node.js 18+
- npm 9+

### Setup

```bash
# Install dependencies for all components
npm run install:all

# Or install individually
cd react && npm install
cd vue && npm install
cd angular && npm install
cd vanilla && npm install
```

### Development Servers

```bash
# React development server
cd react && npm run dev

# Vue development server
cd vue && npm run dev

# Angular development server
cd angular && npm run dev

# Vanilla JavaScript examples
cd vanilla && npm run dev
```

### Building

```bash
# Build all components
npm run build:all

# Build individually
cd react && npm run build
cd vue && npm run build
cd angular && npm run build
cd vanilla && npm run build
```

## Examples and Demos

Each component directory includes:

- **`/examples/`** - Live demo applications
- **`/src/__tests__/`** - Comprehensive test suites
- **`/README.md`** - Detailed documentation
- **`/package.json`** - Build and development scripts

### Live Demos

- [React Demo](./react/examples/)
- [Vue Demo](./vue/examples/)
- [Angular Demo](./angular/examples/)
- [Vanilla JS Demo](./vanilla/examples/)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`npm run test:all`)
6. Update documentation as needed
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Component Guidelines

- Follow the established patterns for consistency
- Maintain TypeScript support across all components
- Include comprehensive tests
- Update documentation
- Follow accessibility best practices
- Ensure responsive design

## Troubleshooting

### Common Issues

**CORS Errors**
```javascript
// Ensure your Go backend includes CORS headers
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
```

**TypeScript Errors**
```bash
# Install type definitions
npm install --save-dev @types/node

# Check tsconfig.json for proper module resolution
```

**Build Issues**
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Getting Help

- 📚 [Documentation](https://github.com/your-org/go-chatbot)
- 🐛 [Issue Tracker](https://github.com/your-org/go-chatbot/issues)
- 💬 [Discussions](https://github.com/your-org/go-chatbot/discussions)
- 📧 [Email Support](mailto:support@go-chatbot.dev)

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

---

**Made with ❤️ by the Go Chatbot Team**

Choose the component that best fits your stack and start building amazing chat experiences!
