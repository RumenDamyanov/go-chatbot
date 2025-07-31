/**
 * @jest-environment jsdom
 */

const GoChatbot = require('../src/go-chatbot.js');

// Mock fetch
global.fetch = jest.fn();

describe('GoChatbot', () => {
  let chatbot;

  beforeEach(() => {
    // Clear DOM
    document.body.innerHTML = '';

    // Reset fetch mock
    fetch.mockClear();

    // Create new chatbot instance
    chatbot = new GoChatbot({
      apiEndpoint: '/test/chat'
    });
  });

  afterEach(() => {
    if (chatbot) {
      chatbot.destroy();
    }
  });

  describe('Initialization', () => {
    it('should create chatbot with default config', () => {
      expect(chatbot.config.apiEndpoint).toBe('/test/chat');
      expect(chatbot.config.placeholder).toBe('Type your message...');
      expect(chatbot.config.maxHeight).toBe('400px');
      expect(chatbot.isOpen).toBe(false);
    });

    it('should create DOM elements', () => {
      expect(document.querySelector('.go-chatbot-container')).toBeTruthy();
      expect(document.querySelector('.go-chatbot-button')).toBeTruthy();
      expect(document.querySelector('.go-chatbot-window')).toBeTruthy();
    });

    it('should apply styles to document', () => {
      expect(document.getElementById('go-chatbot-styles')).toBeTruthy();
    });
  });

  describe('Chat Window Toggle', () => {
    it('should open chat window', () => {
      chatbot.open();

      expect(chatbot.isOpen).toBe(true);
      expect(chatbot.elements.window.style.display).toBe('block');
      expect(chatbot.elements.button.innerHTML).toBe('×');
    });

    it('should close chat window', () => {
      chatbot.open();
      chatbot.close();

      expect(chatbot.isOpen).toBe(false);
      expect(chatbot.elements.button.innerHTML).toBe('💬');
    });

    it('should toggle chat window', () => {
      expect(chatbot.isOpen).toBe(false);

      chatbot.toggle();
      expect(chatbot.isOpen).toBe(true);

      chatbot.toggle();
      expect(chatbot.isOpen).toBe(false);
    });

    it('should not open when disabled', () => {
      chatbot.config.disabled = true;
      chatbot.open();

      expect(chatbot.isOpen).toBe(false);
    });
  });

  describe('Message Handling', () => {
    it('should add message to chat', () => {
      const message = {
        text: 'Hello world',
        sender: 'user'
      };

      chatbot.addMessage(message);

      expect(chatbot.messages).toHaveLength(1);
      expect(chatbot.messages[0].text).toBe('Hello world');
      expect(chatbot.messages[0].sender).toBe('user');
      expect(chatbot.messages[0].id).toBeTruthy();
      expect(chatbot.messages[0].timestamp).toBeInstanceOf(Date);
    });

    it('should clear all messages', () => {
      chatbot.addMessage({ text: 'Message 1', sender: 'user' });
      chatbot.addMessage({ text: 'Message 2', sender: 'bot' });

      expect(chatbot.messages).toHaveLength(2);

      chatbot.clearMessages();

      expect(chatbot.messages).toHaveLength(0);
    });

    it('should generate unique message IDs', () => {
      const id1 = chatbot.generateId();
      const id2 = chatbot.generateId();

      expect(id1).toContain('msg-');
      expect(id2).toContain('msg-');
      expect(id1).not.toBe(id2);
    });
  });

  describe('Message Sending', () => {
    beforeEach(() => {
      chatbot.open();
      chatbot.elements.input.value = 'Test message';
    });

    it('should send message successfully', async () => {
      const mockResponse = {
        success: true,
        response: 'Bot response'
      };

      fetch.mockResolvedValueOnce({
        json: jest.fn().resolvedValueOnce(mockResponse)
      });

      const messageSentSpy = jest.fn();
      const responseReceivedSpy = jest.fn();

      chatbot.on('messageSent', messageSentSpy);
      chatbot.on('responseReceived', responseReceivedSpy);

      await chatbot.sendMessage();

      expect(fetch).toHaveBeenCalledWith('/test/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ message: 'Test message' })
      });

      expect(messageSentSpy).toHaveBeenCalledWith('Test message');
      expect(responseReceivedSpy).toHaveBeenCalledWith('Bot response');
      expect(chatbot.messages).toHaveLength(2);
      expect(chatbot.messages[0].sender).toBe('user');
      expect(chatbot.messages[1].sender).toBe('bot');
    });

    it('should handle API error response', async () => {
      const mockResponse = {
        success: false,
        error: 'Server error'
      };

      fetch.mockResolvedValueOnce({
        json: jest.fn().resolvedValueOnce(mockResponse)
      });

      const errorSpy = jest.fn();
      chatbot.on('error', errorSpy);

      await chatbot.sendMessage();

      expect(errorSpy).toHaveBeenCalledWith('Server error');
      expect(chatbot.messages).toHaveLength(2);
      expect(chatbot.messages[1].text).toContain('Error: Server error');
    });

    it('should handle network error', async () => {
      fetch.mockRejectedValueOnce(new Error('Network error'));

      const errorSpy = jest.fn();
      chatbot.on('error', errorSpy);

      await chatbot.sendMessage();

      expect(errorSpy).toHaveBeenCalledWith('Network error');
      expect(chatbot.messages).toHaveLength(2);
      expect(chatbot.messages[1].text).toContain('Error: Network error');
    });

    it('should not send empty message', async () => {
      chatbot.elements.input.value = '   ';

      await chatbot.sendMessage();

      expect(fetch).not.toHaveBeenCalled();
      expect(chatbot.messages).toHaveLength(0);
    });

    it('should not send when loading', async () => {
      chatbot.isLoading = true;

      await chatbot.sendMessage();

      expect(fetch).not.toHaveBeenCalled();
    });

    it('should not send when disabled', async () => {
      chatbot.config.disabled = true;

      await chatbot.sendMessage();

      expect(fetch).not.toHaveBeenCalled();
    });
  });

  describe('Event System', () => {
    it('should add and trigger event handlers', () => {
      const handler = jest.fn();

      chatbot.on('messageSent', handler);
      chatbot.emit('messageSent', 'test data');

      expect(handler).toHaveBeenCalledWith('test data');
    });

    it('should remove event handlers', () => {
      const handler = jest.fn();

      chatbot.on('messageSent', handler);
      chatbot.off('messageSent', handler);
      chatbot.emit('messageSent', 'test data');

      expect(handler).not.toHaveBeenCalled();
    });

    it('should handle multiple event handlers', () => {
      const handler1 = jest.fn();
      const handler2 = jest.fn();

      chatbot.on('messageSent', handler1);
      chatbot.on('messageSent', handler2);
      chatbot.emit('messageSent', 'test data');

      expect(handler1).toHaveBeenCalledWith('test data');
      expect(handler2).toHaveBeenCalledWith('test data');
    });
  });

  describe('Configuration', () => {
    it('should update configuration', () => {
      chatbot.updateConfig({
        placeholder: 'New placeholder',
        maxHeight: '500px'
      });

      expect(chatbot.config.placeholder).toBe('New placeholder');
      expect(chatbot.config.maxHeight).toBe('500px');
      expect(chatbot.elements.input.placeholder).toBe('New placeholder');
    });

    it('should handle disabled state', () => {
      chatbot.setDisabled(true);

      expect(chatbot.config.disabled).toBe(true);
      expect(chatbot.elements.button.disabled).toBe(true);
      expect(chatbot.elements.input.disabled).toBe(true);
    });

    it('should use initial messages', () => {
      const initialMessages = [
        { id: '1', text: 'Welcome', sender: 'bot', timestamp: new Date() }
      ];

      const newChatbot = new GoChatbot({
        initialMessages
      });

      expect(newChatbot.messages).toHaveLength(1);
      expect(newChatbot.messages[0].text).toBe('Welcome');

      newChatbot.destroy();
    });
  });

  describe('UI Interactions', () => {
    it('should handle button click', () => {
      const button = chatbot.elements.button;

      button.click();
      expect(chatbot.isOpen).toBe(true);

      button.click();
      expect(chatbot.isOpen).toBe(false);
    });

    it('should handle close button click', () => {
      chatbot.open();

      const closeButton = chatbot.elements.header.querySelector('.go-chatbot-close');
      closeButton.click();

      expect(chatbot.isOpen).toBe(false);
    });

    it('should handle send button click', () => {
      chatbot.open();
      chatbot.elements.input.value = 'Test';

      const sendButton = chatbot.elements.sendButton;
      sendButton.click();

      expect(chatbot.messages).toHaveLength(1);
    });

    it('should handle Enter key in input', () => {
      chatbot.open();
      chatbot.elements.input.value = 'Test';

      const enterEvent = new KeyboardEvent('keydown', { key: 'Enter' });
      chatbot.elements.input.dispatchEvent(enterEvent);

      expect(chatbot.messages).toHaveLength(1);
    });

    it('should not send on Shift+Enter', () => {
      chatbot.open();
      chatbot.elements.input.value = 'Test';

      const enterEvent = new KeyboardEvent('keydown', { key: 'Enter', shiftKey: true });
      chatbot.elements.input.dispatchEvent(enterEvent);

      expect(chatbot.messages).toHaveLength(0);
    });
  });

  describe('Utility Functions', () => {
    it('should format time correctly', () => {
      const date = new Date('2023-01-01T12:30:00');
      const formatted = chatbot.formatTime(date);

      expect(formatted).toMatch(/12:30/);
    });

    it('should get correct position styles', () => {
      const styles = chatbot.getPositionStyles();
      expect(styles).toContain('bottom: 20px');
      expect(styles).toContain('right: 20px');
    });

    it('should create message elements', () => {
      const message = {
        id: '1',
        text: 'Test message',
        sender: 'user',
        timestamp: new Date()
      };

      const element = chatbot.createMessageElement(message);

      expect(element.className).toContain('go-chatbot-message--user');
      expect(element.textContent).toContain('Test message');
    });

    it('should create typing indicator', () => {
      const indicator = chatbot.createTypingIndicator();

      expect(indicator.className).toContain('go-chatbot-message--bot');
      expect(indicator.textContent).toContain('Bot is typing');
    });
  });

  describe('Cleanup', () => {
    it('should destroy chatbot instance', () => {
      const container = chatbot.elements.container;
      const styles = document.getElementById('go-chatbot-styles');

      expect(container.parentNode).toBeTruthy();
      expect(styles).toBeTruthy();

      chatbot.destroy();

      expect(container.parentNode).toBeFalsy();
    });
  });

  describe('Auto-initialization', () => {
    beforeEach(() => {
      // Clean up existing chatbot
      if (chatbot) {
        chatbot.destroy();
      }
      document.body.innerHTML = '';
    });

    it('should auto-initialize from data attributes', () => {
      // Create element with data attributes
      const element = document.createElement('div');
      element.setAttribute('data-go-chatbot', '');
      element.setAttribute('data-go-chatbot-api-endpoint', '/custom/chat');
      element.setAttribute('data-go-chatbot-placeholder', 'Custom placeholder');
      document.body.appendChild(element);

      // Trigger auto-initialization
      const event = new Event('DOMContentLoaded');
      document.dispatchEvent(event);

      // Check if chatbot was created
      const chatbotContainer = document.querySelector('.go-chatbot-container');
      expect(chatbotContainer).toBeTruthy();
    });
  });
});
