# Go Chatbot - Angular Component

A modern, customizable Angular component for integrating chat functionality with your Go backend.

## Features

- **Modern Angular**: Built with Angular 17+ and TypeScript
- **Standalone Component**: Can be used independently or as part of Angular modules
- **Reactive**: Uses RxJS for HTTP communication and state management
- **Customizable**: Full control over styling and behavior
- **TypeScript**: Complete type safety with interfaces
- **Accessible**: Keyboard navigation and screen reader support
- **Responsive**: Works on all screen sizes
- **Template Slots**: Custom header content support

## Installation

```bash
npm install @go-chatbot/angular
```

## Basic Usage

### Module Import (Traditional)

```typescript
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { GoChatbotModule } from '@go-chatbot/angular';

import { AppComponent } from './app.component';

@NgModule({
  declarations: [AppComponent],
  imports: [
    BrowserModule,
    HttpClientModule,
    GoChatbotModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
```

### Standalone Component (Modern)

```typescript
import { Component } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { GoChatbotComponent } from '@go-chatbot/angular';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [HttpClientModule, GoChatbotComponent],
  template: `
    <div class="app">
      <h1>My App</h1>
      <go-chatbot
        apiEndpoint="/api/chat"
        (messageSent)="onMessageSent($event)"
        (responseReceived)="onResponseReceived($event)">
      </go-chatbot>
    </div>
  `
})
export class AppComponent {
  onMessageSent(message: string) {
    console.log('User sent:', message);
  }

  onResponseReceived(response: string) {
    console.log('Bot responded:', response);
  }
}
```

### Template Usage

```html
<go-chatbot
  apiEndpoint="/api/chat"
  placeholder="Ask me anything..."
  maxHeight="500px"
  [showTypingIndicator]="true"
  [disabled]="false"
  className="custom-chatbot"
  [style]="{ right: '30px', bottom: '30px' }"
  [initialMessages]="initialMessages"
  (messageSent)="handleMessageSent($event)"
  (responseReceived)="handleResponse($event)"
  (error)="handleError($event)">

  <!-- Custom header content -->
  <div slot="header">
    <img src="/logo.png" alt="Logo" style="height: 20px;">
    <span>Support Chat</span>
  </div>
</go-chatbot>
```

## API Reference

### Component Props

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `apiEndpoint` | `string` | `'/chat/'` | Backend chat API endpoint |
| `className` | `string` | `''` | Additional CSS class for the chat window |
| `placeholder` | `string` | `'Type your message...'` | Input placeholder text |
| `maxHeight` | `string` | `'400px'` | Maximum height of the messages area |
| `style` | `object` | `{}` | Custom styles for the chat window |
| `initialMessages` | `ChatMessage[]` | `[]` | Pre-populated messages |
| `showTypingIndicator` | `boolean` | `true` | Show typing indicator when loading |
| `disabled` | `boolean` | `false` | Disable the entire component |

### Events

| Event | Type | Description |
|-------|------|-------------|
| `messageSent` | `string` | Emitted when user sends a message |
| `responseReceived` | `string` | Emitted when bot response is received |
| `error` | `string` | Emitted when an error occurs |

### Types

```typescript
export interface ChatMessage {
  id: string;
  text: string;
  sender: 'user' | 'bot';
  timestamp: Date;
}

export interface ChatResponse {
  success: boolean;
  response?: string;
  error?: string;
}
```

## Examples

### Custom Styling

```typescript
@Component({
  template: `
    <go-chatbot
      [style]="{
        width: '400px',
        height: '600px',
        borderRadius: '20px'
      }"
      className="my-custom-chat">
    </go-chatbot>
  `,
  styles: [`
    ::ng-deep .my-custom-chat {
      --primary-color: #6366f1;
      --text-color: #1f2937;
      --background-color: #ffffff;
    }

    ::ng-deep .my-custom-chat .go-chatbot-header {
      background: linear-gradient(135deg, #6366f1, #8b5cf6);
    }
  `]
})
export class CustomChatComponent { }
```

### With Initial Messages

```typescript
@Component({
  template: `
    <go-chatbot [initialMessages]="welcomeMessages">
    </go-chatbot>
  `
})
export class WelcomeChatComponent {
  welcomeMessages: ChatMessage[] = [
    {
      id: '1',
      text: 'Hello! How can I help you today?',
      sender: 'bot',
      timestamp: new Date()
    }
  ];
}
```

### Error Handling

```typescript
@Component({
  template: `
    <go-chatbot
      (error)="handleChatError($event)"
      (messageSent)="logMessage($event)">
    </go-chatbot>
  `
})
export class ErrorHandlingComponent {
  handleChatError(error: string) {
    console.error('Chat error:', error);
    // Show user-friendly error message
    this.showNotification('Chat is temporarily unavailable');
  }

  logMessage(message: string) {
    // Analytics or logging
    console.log('User message:', message);
  }

  showNotification(message: string) {
    // Your notification system
  }
}
```

### Custom API Integration

```typescript
@Component({
  template: `
    <go-chatbot
      apiEndpoint="/api/v2/chat/sessions/{{sessionId}}"
      (messageSent)="onMessage($event)">
    </go-chatbot>
  `
})
export class CustomAPIComponent {
  sessionId = 'user-123';

  onMessage(message: string) {
    // Custom analytics or processing
    this.analytics.track('chat_message_sent', {
      message_length: message.length,
      session_id: this.sessionId
    });
  }
}
```

## Backend Integration

The component expects your Go backend to handle POST requests to the specified endpoint with this format:

### Request Format
```json
{
  "message": "User's message text"
}
```

### Response Format
```json
{
  "success": true,
  "response": "Bot's response text"
}
```

### Error Response Format
```json
{
  "success": false,
  "error": "Error message"
}
```

## Development

### Setup

```bash
# Install dependencies
npm install

# Run tests
npm test

# Run tests in watch mode
npm run test:watch

# Build the library
npm run build

# Lint code
npm run lint
```

### Testing

The component includes comprehensive tests covering:

- Component rendering and initialization
- User interactions (typing, clicking, keyboard navigation)
- API communication and error handling
- Event emission
- Prop validation
- Custom styling and configuration

Run tests with:

```bash
npm test
```

## Browser Support

- Chrome 88+
- Firefox 78+
- Safari 14+
- Edge 88+

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
