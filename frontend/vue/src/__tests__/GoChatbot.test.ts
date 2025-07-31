import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import GoChatbot from '../GoChatbot.vue';
import type { ChatMessage } from '../GoChatbot.vue';

// Mock fetch globally
global.fetch = vi.fn();

describe('GoChatbot', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders chat button by default', () => {
    const wrapper = mount(GoChatbot);

    const chatButton = wrapper.find('.go-chatbot-button');
    expect(chatButton.exists()).toBe(true);
    expect(chatButton.text()).toBe('💬');
  });

  it('opens chat window when button is clicked', async () => {
    const wrapper = mount(GoChatbot);

    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    expect(wrapper.find('.go-chatbot-window').exists()).toBe(true);
    expect(wrapper.find('.go-chatbot-title').text()).toBe('Chat Support');
    expect(wrapper.find('.go-chatbot-input').exists()).toBe(true);
  });

  it('closes chat window when close button is clicked', async () => {
    const wrapper = mount(GoChatbot);

    // Open chat
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    // Close chat
    const closeButton = wrapper.find('.go-chatbot-close');
    await closeButton.trigger('click');

    expect(wrapper.find('.go-chatbot-window').exists()).toBe(false);
  });

  it('sends message and displays response', async () => {
    const mockResponse = {
      success: true,
      response: 'Hello! How can I help you?'
    };

    (fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse,
    });

    const wrapper = mount(GoChatbot);

    // Open chat
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    // Type and send message
    const input = wrapper.find('.go-chatbot-input');
    const sendButton = wrapper.find('.go-chatbot-send');

    await input.setValue('Hello');
    await sendButton.trigger('click');

    // Wait for response
    await wrapper.vm.$nextTick();
    await new Promise(resolve => setTimeout(resolve, 0));

    const messages = wrapper.findAll('.go-chatbot-message');
    expect(messages.length).toBe(2);
    expect(messages[0].text()).toContain('Hello');
    expect(messages[1].text()).toContain('Hello! How can I help you?');

    expect(fetch).toHaveBeenCalledWith('/chat/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message: 'Hello' }),
    });
  });

  it('handles API errors gracefully', async () => {
    const mockErrorResponse = {
      success: false,
      error: 'Server error'
    };

    (fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => mockErrorResponse,
    });

    const wrapper = mount(GoChatbot);

    // Open chat and send message
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    const input = wrapper.find('.go-chatbot-input');
    const sendButton = wrapper.find('.go-chatbot-send');

    await input.setValue('Hello');
    await sendButton.trigger('click');

    // Wait for error message
    await wrapper.vm.$nextTick();
    await new Promise(resolve => setTimeout(resolve, 0));

    const messages = wrapper.findAll('.go-chatbot-message');
    expect(messages[1].text()).toContain('Error: Server error');
  });

  it('handles network errors', async () => {
    (fetch as any).mockRejectedValueOnce(new Error('Network error'));

    const wrapper = mount(GoChatbot);

    // Open chat and send message
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    const input = wrapper.find('.go-chatbot-input');
    const sendButton = wrapper.find('.go-chatbot-send');

    await input.setValue('Hello');
    await sendButton.trigger('click');

    // Wait for error message
    await wrapper.vm.$nextTick();
    await new Promise(resolve => setTimeout(resolve, 0));

    const messages = wrapper.findAll('.go-chatbot-message');
    expect(messages[1].text()).toContain('Error: Network error');
  });

  it('sends message on Enter key press', async () => {
    const mockResponse = {
      success: true,
      response: 'Response'
    };

    (fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse,
    });

    const wrapper = mount(GoChatbot);

    // Open chat
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    // Type message and press Enter
    const input = wrapper.find('.go-chatbot-input');
    await input.setValue('Hello');
    await input.trigger('keydown.enter');

    await wrapper.vm.$nextTick();

    expect(fetch).toHaveBeenCalled();
  });

  it('does not send empty messages', async () => {
    const wrapper = mount(GoChatbot);

    // Open chat
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    // Try to send empty message
    const sendButton = wrapper.find('.go-chatbot-send');
    await sendButton.trigger('click');

    expect(fetch).not.toHaveBeenCalled();
  });

  it('displays initial messages', async () => {
    const initialMessages: ChatMessage[] = [
      {
        id: '1',
        text: 'Welcome!',
        sender: 'bot',
        timestamp: new Date(),
      },
      {
        id: '2',
        text: 'Hello',
        sender: 'user',
        timestamp: new Date(),
      },
    ];

    const wrapper = mount(GoChatbot, {
      props: { initialMessages }
    });

    // Open chat
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    const messages = wrapper.findAll('.go-chatbot-message');
    expect(messages[0].text()).toContain('Welcome!');
    expect(messages[1].text()).toContain('Hello');
  });

  it('shows typing indicator when loading', async () => {
    (fetch as any).mockImplementation(() =>
      new Promise(resolve => setTimeout(() => resolve({
        ok: true,
        json: async () => ({ success: true, response: 'Response' }),
      }), 100))
    );

    const wrapper = mount(GoChatbot, {
      props: { showTypingIndicator: true }
    });

    // Open chat and send message
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    const input = wrapper.find('.go-chatbot-input');
    const sendButton = wrapper.find('.go-chatbot-send');

    await input.setValue('Hello');
    await sendButton.trigger('click');

    await wrapper.vm.$nextTick();

    expect(wrapper.find('.go-chatbot-typing').exists()).toBe(true);
  });

  it('can be disabled', () => {
    const wrapper = mount(GoChatbot, {
      props: { disabled: true }
    });

    const chatButton = wrapper.find('.go-chatbot-button');
    expect(chatButton.attributes('disabled')).toBeDefined();
  });

  it('uses custom API endpoint', async () => {
    const mockResponse = {
      success: true,
      response: 'Response'
    };

    (fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse,
    });

    const wrapper = mount(GoChatbot, {
      props: { apiEndpoint: '/custom/chat' }
    });

    // Open chat and send message
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    const input = wrapper.find('.go-chatbot-input');
    const sendButton = wrapper.find('.go-chatbot-send');

    await input.setValue('Hello');
    await sendButton.trigger('click');

    await wrapper.vm.$nextTick();

    expect(fetch).toHaveBeenCalledWith('/custom/chat', expect.any(Object));
  });

  it('renders custom header slot', async () => {
    const wrapper = mount(GoChatbot, {
      slots: {
        header: '<div>Custom Support</div>'
      }
    });

    // Open chat
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    expect(wrapper.text()).toContain('Custom Support');
  });

  it('applies custom className and styles', async () => {
    const wrapper = mount(GoChatbot, {
      props: {
        className: 'custom-chat',
        style: { backgroundColor: 'red' }
      }
    });

    // Open chat
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    const chatWindow = wrapper.find('.go-chatbot-window');
    expect(chatWindow.classes()).toContain('custom-chat');
    expect(chatWindow.attributes('style')).toContain('background-color: red');
  });

  it('emits events correctly', async () => {
    const mockResponse = {
      success: true,
      response: 'Bot response'
    };

    (fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse,
    });

    const wrapper = mount(GoChatbot);

    // Open chat and send message
    const chatButton = wrapper.find('.go-chatbot-button');
    await chatButton.trigger('click');

    const input = wrapper.find('.go-chatbot-input');
    const sendButton = wrapper.find('.go-chatbot-send');

    await input.setValue('Hello');
    await sendButton.trigger('click');

    await wrapper.vm.$nextTick();
    await new Promise(resolve => setTimeout(resolve, 0));

    const emitted = wrapper.emitted();
    expect(emitted.messageSent).toBeTruthy();
    expect(emitted.messageSent[0]).toEqual(['Hello']);
    expect(emitted.responseReceived).toBeTruthy();
    expect(emitted.responseReceived[0]).toEqual(['Bot response']);
  });
});
