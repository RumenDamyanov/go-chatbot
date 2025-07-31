/**
 * Go Chatbot - Vanilla JavaScript Component
 * A framework-agnostic chat widget for Go backend integration
 */

class GoChatbot {
  constructor(options = {}) {
    // Configuration
    this.config = {
      apiEndpoint: '/chat/',
      placeholder: 'Type your message...',
      maxHeight: '400px',
      className: '',
      style: {},
      initialMessages: [],
      showTypingIndicator: true,
      disabled: false,
      position: 'bottom-right',
      theme: 'light',
      ...options
    };

    // State
    this.messages = [...this.config.initialMessages];
    this.isLoading = false;
    this.isOpen = false;
    this.messageId = 0;

    // Event handlers
    this.eventHandlers = {
      messageSent: [],
      responseReceived: [],
      error: [],
      opened: [],
      closed: []
    };

    // DOM elements
    this.elements = {};

    // Initialize
    this.init();
  }

  /**
   * Initialize the chatbot
   */
  init() {
    this.createElements();
    this.attachEventListeners();
    this.applyStyles();
    this.render();
  }

  /**
   * Create DOM elements
   */
  createElements() {
    // Main container
    this.elements.container = document.createElement('div');
    this.elements.container.className = 'go-chatbot-container';

    // Chat button
    this.elements.button = document.createElement('button');
    this.elements.button.className = 'go-chatbot-button';
    this.elements.button.innerHTML = '💬';
    this.elements.button.setAttribute('aria-label', 'Open chat');

    // Chat window
    this.elements.window = document.createElement('div');
    this.elements.window.className = `go-chatbot-window ${this.config.className}`;
    this.elements.window.style.display = 'none';

    // Header
    this.elements.header = document.createElement('div');
    this.elements.header.className = 'go-chatbot-header';
    this.elements.header.innerHTML = `
      <span class="go-chatbot-title">Chat Support</span>
      <button class="go-chatbot-close" aria-label="Close chat">×</button>
    `;

    // Messages container
    this.elements.messages = document.createElement('div');
    this.elements.messages.className = 'go-chatbot-messages';
    this.elements.messages.style.maxHeight = this.config.maxHeight;

    // Input container
    this.elements.inputContainer = document.createElement('div');
    this.elements.inputContainer.className = 'go-chatbot-input-container';

    this.elements.input = document.createElement('input');
    this.elements.input.type = 'text';
    this.elements.input.className = 'go-chatbot-input';
    this.elements.input.placeholder = this.config.placeholder;

    this.elements.sendButton = document.createElement('button');
    this.elements.sendButton.className = 'go-chatbot-send';
    this.elements.sendButton.textContent = 'Send';

    // Assemble elements
    this.elements.inputContainer.appendChild(this.elements.input);
    this.elements.inputContainer.appendChild(this.elements.sendButton);

    this.elements.window.appendChild(this.elements.header);
    this.elements.window.appendChild(this.elements.messages);
    this.elements.window.appendChild(this.elements.inputContainer);

    this.elements.container.appendChild(this.elements.button);
    this.elements.container.appendChild(this.elements.window);

    // Add to DOM
    document.body.appendChild(this.elements.container);
  }

  /**
   * Attach event listeners
   */
  attachEventListeners() {
    // Button click
    this.elements.button.addEventListener('click', () => this.toggle());

    // Close button
    this.elements.header.querySelector('.go-chatbot-close')
      .addEventListener('click', () => this.close());

    // Send button
    this.elements.sendButton.addEventListener('click', () => this.sendMessage());

    // Enter key
    this.elements.input.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        this.sendMessage();
      }
    });

    // Input changes
    this.elements.input.addEventListener('input', () => this.updateSendButton());
  }

  /**
   * Apply default styles
   */
  applyStyles() {
    if (document.getElementById('go-chatbot-styles')) return;

    const styles = document.createElement('style');
    styles.id = 'go-chatbot-styles';
    styles.textContent = `
      .go-chatbot-container {
        position: fixed;
        z-index: 9999;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        ${this.getPositionStyles()}
      }

      .go-chatbot-button {
        width: 60px;
        height: 60px;
        border-radius: 50%;
        background-color: #007bff;
        color: white;
        border: none;
        cursor: pointer;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        font-size: 24px;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .go-chatbot-button:hover {
        background-color: #0056b3;
        transform: scale(1.05);
      }

      .go-chatbot-button:disabled {
        cursor: not-allowed;
        opacity: 0.6;
      }

      .go-chatbot-window {
        position: absolute;
        ${this.getWindowPositionStyles()}
        width: 350px;
        max-width: 90vw;
        background-color: white;
        border-radius: 12px;
        box-shadow: 0 8px 32px rgba(0,0,0,0.12);
        overflow: hidden;
        transform: scale(0.8) translateY(20px);
        opacity: 0;
        transition: all 0.3s ease;
        max-height: 80vh;
      }

      .go-chatbot-window.open {
        transform: scale(1) translateY(0);
        opacity: 1;
      }

      .go-chatbot-header {
        padding: 16px;
        background-color: #007bff;
        color: white;
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .go-chatbot-title {
        font-weight: bold;
        font-size: 16px;
      }

      .go-chatbot-close {
        background: none;
        border: none;
        color: white;
        cursor: pointer;
        font-size: 18px;
        padding: 0;
        width: 24px;
        height: 24px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .go-chatbot-close:hover {
        background-color: rgba(255,255,255,0.1);
      }

      .go-chatbot-messages {
        flex: 1;
        overflow-y: auto;
        padding: 16px;
        display: flex;
        flex-direction: column;
        gap: 12px;
        min-height: 200px;
      }

      .go-chatbot-message {
        display: flex;
        animation: slideIn 0.3s ease;
      }

      @keyframes slideIn {
        from {
          opacity: 0;
          transform: translateY(10px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }

      .go-chatbot-message--user {
        justify-content: flex-end;
      }

      .go-chatbot-message--bot {
        justify-content: flex-start;
      }

      .go-chatbot-message-content {
        max-width: 80%;
        padding: 8px 12px;
        border-radius: 18px;
        font-size: 14px;
        line-height: 1.4;
        word-wrap: break-word;
      }

      .go-chatbot-message--user .go-chatbot-message-content {
        background-color: #007bff;
        color: white;
      }

      .go-chatbot-message--bot .go-chatbot-message-content {
        background-color: #f1f3f5;
        color: #333;
      }

      .go-chatbot-message-time {
        font-size: 10px;
        opacity: 0.7;
        margin-top: 4px;
        text-align: right;
      }

      .go-chatbot-typing {
        display: flex;
        gap: 8px;
        align-items: center;
        padding: 8px 12px;
        background-color: #f1f3f5;
        border-radius: 18px;
        max-width: 80%;
      }

      .go-chatbot-dots {
        display: flex;
        gap: 2px;
      }

      .go-chatbot-dot {
        width: 4px;
        height: 4px;
        border-radius: 50%;
        background-color: #999;
        animation: pulse 1.5s infinite;
      }

      .go-chatbot-dot:nth-child(2) {
        animation-delay: 0.5s;
      }

      .go-chatbot-dot:nth-child(3) {
        animation-delay: 1s;
      }

      @keyframes pulse {
        0%, 100% { opacity: 0.4; }
        50% { opacity: 1; }
      }

      .go-chatbot-input-container {
        padding: 16px;
        border-top: 1px solid #e9ecef;
        display: flex;
        gap: 8px;
      }

      .go-chatbot-input {
        flex: 1;
        padding: 8px 12px;
        border: 1px solid #dee2e6;
        border-radius: 20px;
        outline: none;
        font-size: 14px;
        font-family: inherit;
      }

      .go-chatbot-input:focus {
        border-color: #007bff;
      }

      .go-chatbot-send {
        padding: 8px 16px;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 20px;
        cursor: pointer;
        font-size: 14px;
        transition: all 0.2s ease;
      }

      .go-chatbot-send:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      .go-chatbot-send:hover:not(:disabled) {
        background-color: #0056b3;
      }

      /* Dark theme */
      .go-chatbot-container.theme-dark .go-chatbot-window {
        background-color: #1a1a1a;
        color: white;
      }

      .go-chatbot-container.theme-dark .go-chatbot-message--bot .go-chatbot-message-content {
        background-color: #333;
        color: white;
      }

      .go-chatbot-container.theme-dark .go-chatbot-input {
        background-color: #333;
        border-color: #555;
        color: white;
      }

      .go-chatbot-container.theme-dark .go-chatbot-input-container {
        border-top-color: #555;
      }
    `;

    document.head.appendChild(styles);
  }

  /**
   * Get position styles based on config
   */
  getPositionStyles() {
    const positions = {
      'bottom-right': 'bottom: 20px; right: 20px;',
      'bottom-left': 'bottom: 20px; left: 20px;',
      'top-right': 'top: 20px; right: 20px;',
      'top-left': 'top: 20px; left: 20px;'
    };
    return positions[this.config.position] || positions['bottom-right'];
  }

  /**
   * Get window position styles
   */
  getWindowPositionStyles() {
    const isBottom = this.config.position.includes('bottom');
    const isRight = this.config.position.includes('right');

    let styles = '';
    if (isBottom) styles += 'bottom: 80px; ';
    else styles += 'top: 80px; ';

    if (isRight) styles += 'right: 0; ';
    else styles += 'left: 0; ';

    return styles;
  }

  /**
   * Render messages
   */
  render() {
    this.elements.messages.innerHTML = '';

    this.messages.forEach(message => {
      const messageEl = this.createMessageElement(message);
      this.elements.messages.appendChild(messageEl);
    });

    if (this.isLoading && this.config.showTypingIndicator) {
      const typingEl = this.createTypingIndicator();
      this.elements.messages.appendChild(typingEl);
    }

    this.scrollToBottom();
    this.updateSendButton();
  }

  /**
   * Create message element
   */
  createMessageElement(message) {
    const messageEl = document.createElement('div');
    messageEl.className = `go-chatbot-message go-chatbot-message--${message.sender}`;

    const contentEl = document.createElement('div');
    contentEl.className = 'go-chatbot-message-content';

    const textEl = document.createElement('div');
    textEl.className = 'go-chatbot-message-text';
    textEl.textContent = message.text;

    const timeEl = document.createElement('div');
    timeEl.className = 'go-chatbot-message-time';
    timeEl.textContent = this.formatTime(message.timestamp);

    contentEl.appendChild(textEl);
    contentEl.appendChild(timeEl);
    messageEl.appendChild(contentEl);

    return messageEl;
  }

  /**
   * Create typing indicator
   */
  createTypingIndicator() {
    const container = document.createElement('div');
    container.className = 'go-chatbot-message go-chatbot-message--bot';

    const typing = document.createElement('div');
    typing.className = 'go-chatbot-typing';
    typing.innerHTML = `
      <span>Bot is typing</span>
      <div class="go-chatbot-dots">
        <div class="go-chatbot-dot"></div>
        <div class="go-chatbot-dot"></div>
        <div class="go-chatbot-dot"></div>
      </div>
    `;

    container.appendChild(typing);
    return container;
  }

  /**
   * Generate unique message ID
   */
  generateId() {
    return `msg-${Date.now()}-${++this.messageId}`;
  }

  /**
   * Format timestamp
   */
  formatTime(date) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  /**
   * Scroll to bottom of messages
   */
  scrollToBottom() {
    setTimeout(() => {
      this.elements.messages.scrollTop = this.elements.messages.scrollHeight;
    }, 0);
  }

  /**
   * Update send button state
   */
  updateSendButton() {
    const hasText = this.elements.input.value.trim().length > 0;
    this.elements.sendButton.disabled = !hasText || this.isLoading || this.config.disabled;
    this.elements.sendButton.textContent = this.isLoading ? '...' : 'Send';
  }

  /**
   * Toggle chat window
   */
  toggle() {
    if (this.isOpen) {
      this.close();
    } else {
      this.open();
    }
  }

  /**
   * Open chat window
   */
  open() {
    if (this.config.disabled) return;

    this.isOpen = true;
    this.elements.window.style.display = 'block';
    this.elements.button.innerHTML = '×';

    setTimeout(() => {
      this.elements.window.classList.add('open');
      this.elements.input.focus();
    }, 10);

    this.emit('opened');
  }

  /**
   * Close chat window
   */
  close() {
    this.isOpen = false;
    this.elements.window.classList.remove('open');
    this.elements.button.innerHTML = '💬';

    setTimeout(() => {
      this.elements.window.style.display = 'none';
    }, 300);

    this.emit('closed');
  }

  /**
   * Send message
   */
  async sendMessage() {
    const text = this.elements.input.value.trim();
    if (!text || this.isLoading || this.config.disabled) return;

    // Add user message
    const userMessage = {
      id: this.generateId(),
      text,
      sender: 'user',
      timestamp: new Date()
    };

    this.messages.push(userMessage);
    this.elements.input.value = '';
    this.isLoading = true;

    this.render();
    this.emit('messageSent', text);

    try {
      const response = await fetch(this.config.apiEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ message: text })
      });

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.error || 'Unknown error occurred');
      }

      // Add bot message
      const botMessage = {
        id: this.generateId(),
        text: data.response,
        sender: 'bot',
        timestamp: new Date()
      };

      this.messages.push(botMessage);
      this.emit('responseReceived', data.response);

    } catch (error) {
      // Add error message
      const errorMessage = {
        id: this.generateId(),
        text: `Error: ${error.message}`,
        sender: 'bot',
        timestamp: new Date()
      };

      this.messages.push(errorMessage);
      this.emit('error', error.message);

    } finally {
      this.isLoading = false;
      this.render();
    }
  }

  /**
   * Add event listener
   */
  on(event, handler) {
    if (this.eventHandlers[event]) {
      this.eventHandlers[event].push(handler);
    }
  }

  /**
   * Remove event listener
   */
  off(event, handler) {
    if (this.eventHandlers[event]) {
      const index = this.eventHandlers[event].indexOf(handler);
      if (index > -1) {
        this.eventHandlers[event].splice(index, 1);
      }
    }
  }

  /**
   * Emit event
   */
  emit(event, data) {
    if (this.eventHandlers[event]) {
      this.eventHandlers[event].forEach(handler => handler(data));
    }
  }

  /**
   * Update configuration
   */
  updateConfig(newConfig) {
    this.config = { ...this.config, ...newConfig };
    this.applyStyles();
    this.render();
  }

  /**
   * Add message programmatically
   */
  addMessage(message) {
    const fullMessage = {
      id: this.generateId(),
      timestamp: new Date(),
      ...message
    };
    this.messages.push(fullMessage);
    this.render();
  }

  /**
   * Clear all messages
   */
  clearMessages() {
    this.messages = [];
    this.render();
  }

  /**
   * Enable/disable the chatbot
   */
  setDisabled(disabled) {
    this.config.disabled = disabled;
    this.elements.button.disabled = disabled;
    this.elements.input.disabled = disabled;
    this.updateSendButton();
  }

  /**
   * Destroy the chatbot
   */
  destroy() {
    if (this.elements.container && this.elements.container.parentNode) {
      this.elements.container.parentNode.removeChild(this.elements.container);
    }

    const styles = document.getElementById('go-chatbot-styles');
    if (styles && styles.parentNode) {
      styles.parentNode.removeChild(styles);
    }
  }
}

// Auto-initialize if data-go-chatbot attribute is present
document.addEventListener('DOMContentLoaded', () => {
  const autoInit = document.querySelector('[data-go-chatbot]');
  if (autoInit) {
    const config = {};

    // Parse data attributes
    Object.keys(autoInit.dataset).forEach(key => {
      if (key.startsWith('goChatbot')) {
        const configKey = key.replace('goChatbot', '').toLowerCase();
        config[configKey] = autoInit.dataset[key];
      }
    });

    new GoChatbot(config);
  }
});

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
  module.exports = GoChatbot;
}

// Global for script tag usage
if (typeof window !== 'undefined') {
  window.GoChatbot = GoChatbot;
}
