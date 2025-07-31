# Go Chatbot - Vanilla JavaScript Component

A framework-agnostic, lightweight chat widget that integrates seamlessly with your Go backend. Works with any website, application, or framework without dependencies.

## Features

- **Zero Dependencies**: Pure vanilla JavaScript, works everywhere
- **Lightweight**: ~15KB minified, ~5KB gzipped
- **Framework Agnostic**: Works with React, Vue, Angular, or plain HTML
- **Customizable**: Full control over styling, behavior, and positioning
- **TypeScript Support**: Complete type definitions included
- **Responsive**: Mobile-friendly and accessible
- **Modern**: ES6+ with fallbacks, async/await support
- **Flexible**: Event-driven architecture with comprehensive API

## Installation

### CDN (Recommended for quick start)

```html
<script src="https://unpkg.com/@go-chatbot/vanilla@latest/dist/go-chatbot.min.js"></script>
```

### NPM

```bash
npm install @go-chatbot/vanilla
```

### Download

Download `go-chatbot.min.js` from the [releases page](https://github.com/your-org/go-chatbot/releases).

## Quick Start

### HTML + Script Tag

```html
<!DOCTYPE html>
<html>
<head>
    <title>My Website</title>
</head>
<body>
    <h1>Welcome to my website</h1>

    <!-- Your existing content -->

    <!-- Load the chatbot -->
    <script src="https://unpkg.com/@go-chatbot/vanilla@latest/dist/go-chatbot.min.js"></script>
    <script>
        // Initialize the chatbot
        const chatbot = new GoChatbot({
            apiEndpoint: '/api/chat'
        });
    </script>
</body>
</html>
```

### Auto-Initialization with Data Attributes

```html
<!-- Add this element anywhere in your HTML -->
<div data-go-chatbot
     data-go-chatbot-api-endpoint="/api/chat"
     data-go-chatbot-placeholder="How can I help you?">
</div>

<script src="https://unpkg.com/@go-chatbot/vanilla@latest/dist/go-chatbot.min.js"></script>
<!-- Chatbot will auto-initialize on DOMContentLoaded -->
```

### ES6 Modules

```javascript
import GoChatbot from '@go-chatbot/vanilla';

const chatbot = new GoChatbot({
    apiEndpoint: '/api/chat',
    placeholder: 'Ask me anything...'
});
```

### CommonJS

```javascript
const GoChatbot = require('@go-chatbot/vanilla');

const chatbot = new GoChatbot({
    apiEndpoint: '/api/chat'
});
```

## Configuration Options

```javascript
const chatbot = new GoChatbot({
    // API Configuration
    apiEndpoint: '/chat/',              // Backend chat endpoint

    // UI Configuration
    placeholder: 'Type your message...', // Input placeholder text
    maxHeight: '400px',                 // Max height of messages area
    className: '',                      // Custom CSS class for styling
    style: {},                          // Custom inline styles

    // Behavior Configuration
    showTypingIndicator: true,          // Show "bot is typing" indicator
    disabled: false,                    // Disable the entire widget

    // Positioning
    position: 'bottom-right',           // 'bottom-right', 'bottom-left', 'top-right', 'top-left'

    // Theming
    theme: 'light',                     // 'light' or 'dark'

    // Initial State
    initialMessages: []                 // Pre-populate with messages
});
```

## API Reference

### Methods

#### `open()`
Open the chat window.

```javascript
chatbot.open();
```

#### `close()`
Close the chat window.

```javascript
chatbot.close();
```

#### `toggle()`
Toggle the chat window open/closed state.

```javascript
chatbot.toggle();
```

#### `sendMessage()`
Programmatically send the current input value.

```javascript
chatbot.sendMessage();
```

#### `addMessage(message)`
Add a message to the chat programmatically.

```javascript
chatbot.addMessage({
    text: 'Hello from the system!',
    sender: 'bot'  // 'user' or 'bot'
});
```

#### `clearMessages()`
Clear all messages from the chat.

```javascript
chatbot.clearMessages();
```

#### `setDisabled(disabled)`
Enable or disable the chatbot.

```javascript
chatbot.setDisabled(true);  // Disable
chatbot.setDisabled(false); // Enable
```

#### `updateConfig(newConfig)`
Update configuration options.

```javascript
chatbot.updateConfig({
    placeholder: 'New placeholder text',
    theme: 'dark'
});
```

#### `destroy()`
Remove the chatbot from the DOM and clean up.

```javascript
chatbot.destroy();
```

### Events

#### `on(event, handler)`
Add an event listener.

```javascript
chatbot.on('messageSent', (message) => {
    console.log('User sent:', message);
});

chatbot.on('responseReceived', (response) => {
    console.log('Bot responded:', response);
});

chatbot.on('error', (error) => {
    console.error('Chat error:', error);
});

chatbot.on('opened', () => {
    console.log('Chat opened');
});

chatbot.on('closed', () => {
    console.log('Chat closed');
});
```

#### `off(event, handler)`
Remove an event listener.

```javascript
const handler = (message) => console.log(message);

chatbot.on('messageSent', handler);
chatbot.off('messageSent', handler);
```

### Properties

Access internal state (read-only):

```javascript
chatbot.isOpen      // Boolean: is chat window open
chatbot.isLoading   // Boolean: is API request in progress
chatbot.messages    // Array: all chat messages
chatbot.config      // Object: current configuration
```

## Examples

### Basic Usage

```javascript
const chatbot = new GoChatbot({
    apiEndpoint: '/api/chat'
});
```

### Custom Styling

```javascript
const chatbot = new GoChatbot({
    className: 'my-custom-chatbot',
    style: {
        width: '400px',
        height: '600px'
    }
});
```

Add custom CSS:

```css
.my-custom-chatbot .go-chatbot-header {
    background: linear-gradient(135deg, #667eea, #764ba2);
}

.my-custom-chatbot .go-chatbot-button {
    background: #667eea;
}
```

### Dark Theme

```javascript
const chatbot = new GoChatbot({
    theme: 'dark',
    className: 'dark-chatbot'
});
```

### Different Positioning

```javascript
// Top left corner
const chatbot = new GoChatbot({
    position: 'top-left'
});

// Bottom left corner
const chatbot = new GoChatbot({
    position: 'bottom-left'
});
```

### With Welcome Message

```javascript
const chatbot = new GoChatbot({
    initialMessages: [
        {
            text: 'Hello! How can I help you today?',
            sender: 'bot',
            timestamp: new Date()
        }
    ]
});

// Auto-open to show welcome message
chatbot.open();
```

### Event Handling

```javascript
const chatbot = new GoChatbot({
    apiEndpoint: '/api/chat'
});

// Analytics tracking
chatbot.on('messageSent', (message) => {
    analytics.track('chat_message_sent', {
        message_length: message.length
    });
});

// Error handling
chatbot.on('error', (error) => {
    console.error('Chat error:', error);
    // Show user-friendly error message
    showNotification('Sorry, the chat is temporarily unavailable');
});

// Custom responses
chatbot.on('responseReceived', (response) => {
    if (response.includes('contact')) {
        // Show contact form
        showContactForm();
    }
});
```

### Integration with Forms

```javascript
// Initialize chatbot
const chatbot = new GoChatbot({
    apiEndpoint: '/api/support-chat'
});

// Pre-populate with form data
document.getElementById('helpButton').addEventListener('click', () => {
    const userEmail = document.getElementById('email').value;
    const issue = document.getElementById('issue').value;

    chatbot.addMessage({
        text: `Hi! I'm having trouble with: ${issue}. My email is ${userEmail}`,
        sender: 'user'
    });

    chatbot.open();
});
```

### Multiple Chatbots

```javascript
// Support chatbot
const supportBot = new GoChatbot({
    apiEndpoint: '/api/support',
    position: 'bottom-right',
    className: 'support-bot'
});

// Sales chatbot
const salesBot = new GoChatbot({
    apiEndpoint: '/api/sales',
    position: 'bottom-left',
    className: 'sales-bot'
});

// Show different bots on different pages
if (window.location.pathname.includes('/support')) {
    salesBot.destroy();
} else if (window.location.pathname.includes('/pricing')) {
    supportBot.destroy();
}
```

## Backend Integration

The chatbot sends POST requests to your API endpoint with this format:

### Request Format

```json
{
  "message": "User's message text"
}
```

### Response Format

**Success:**
```json
{
  "success": true,
  "response": "Bot's response text"
}
```

**Error:**
```json
{
  "success": false,
  "error": "Error message"
}
```

### Go Backend Example

```go
func chatHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Message string `json:"message"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Process the message (call your AI service, database, etc.)
    response := processChatMessage(req.Message)

    resp := struct {
        Success  bool   `json:"success"`
        Response string `json:"response"`
    }{
        Success:  true,
        Response: response,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
```

## TypeScript Support

Full TypeScript definitions are included:

```typescript
import GoChatbot, { ChatMessage, GoChatbotConfig } from '@go-chatbot/vanilla';

const config: GoChatbotConfig = {
    apiEndpoint: '/api/chat',
    placeholder: 'Type here...',
    showTypingIndicator: true
};

const chatbot = new GoChatbot(config);

chatbot.on('messageSent', (message: string) => {
    console.log('Message sent:', message);
});
```

## Styling Customization

### CSS Custom Properties

```css
.go-chatbot-container {
    --primary-color: #007bff;
    --text-color: #333;
    --background-color: #ffffff;
    --border-radius: 12px;
    --shadow: 0 8px 32px rgba(0,0,0,0.12);
}
```

### Custom Themes

```css
/* Dark theme */
.go-chatbot-container.dark-theme {
    --primary-color: #4a5568;
    --text-color: #ffffff;
    --background-color: #2d3748;
    --message-bg-bot: #4a5568;
    --message-bg-user: #3182ce;
}

/* Brand theme */
.go-chatbot-container.brand-theme {
    --primary-color: #your-brand-color;
    --border-radius: 8px;
}
```

### Responsive Design

The chatbot is responsive by default, but you can customize breakpoints:

```css
@media (max-width: 768px) {
    .go-chatbot-window {
        width: 100vw !important;
        height: 100vh !important;
        border-radius: 0 !important;
        bottom: 0 !important;
        right: 0 !important;
    }
}
```

## Browser Support

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+
- iOS Safari 12+
- Android Chrome 60+

## Performance

- **Size**: ~15KB minified, ~5KB gzipped
- **Runtime**: Minimal memory footprint
- **Network**: Efficient API calls with error handling
- **Rendering**: Optimized DOM updates

## Development

### Setup

```bash
git clone https://github.com/your-org/go-chatbot.git
cd go-chatbot/frontend/vanilla
npm install
```

### Development Server

```bash
npm run dev
# Opens http://localhost:8080 with examples
```

### Testing

```bash
npm test              # Run tests
npm run test:watch    # Watch mode
npm run test:coverage # Coverage report
```

### Building

```bash
npm run build         # Create minified version
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.

## Support

- 📚 [Documentation](https://github.com/your-org/go-chatbot)
- 🐛 [Issue Tracker](https://github.com/your-org/go-chatbot/issues)
- 💬 [Discussions](https://github.com/your-org/go-chatbot/discussions)
- 📧 [Email Support](mailto:support@go-chatbot.dev)
