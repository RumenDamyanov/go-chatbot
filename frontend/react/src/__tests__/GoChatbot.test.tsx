import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';
import GoChatbot, { ChatMessage } from '../GoChatbot';

// Mock fetch globally
global.fetch = jest.fn();

describe('GoChatbot', () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockClear();
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it('renders chat button by default', () => {
    render(<GoChatbot />);

    const chatButton = screen.getByRole('button');
    expect(chatButton).toBeInTheDocument();
    expect(chatButton).toHaveTextContent('💬');
  });

  it('opens chat window when button is clicked', async () => {
    render(<GoChatbot />);

    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    expect(screen.getByText('Chat Support')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Type your message...')).toBeInTheDocument();
  });

  it('closes chat window when close button is clicked', async () => {
    render(<GoChatbot />);

    // Open chat
    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    // Close chat
    const closeButton = screen.getByText('×');
    await userEvent.click(closeButton);

    expect(screen.queryByText('Chat Support')).not.toBeInTheDocument();
  });

  it('sends message and displays response', async () => {
    const mockResponse = {
      success: true,
      response: 'Hello! How can I help you?'
    };

    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse,
    });

    const onMessageSent = jest.fn();
    const onResponseReceived = jest.fn();

    render(
      <GoChatbot
        onMessageSent={onMessageSent}
        onResponseReceived={onResponseReceived}
      />
    );

    // Open chat
    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    // Type and send message
    const input = screen.getByPlaceholderText('Type your message...');
    const sendButton = screen.getByText('Send');

    await userEvent.type(input, 'Hello');
    await userEvent.click(sendButton);

    // Wait for response
    await waitFor(() => {
      expect(screen.getByText('Hello')).toBeInTheDocument();
      expect(screen.getByText('Hello! How can I help you?')).toBeInTheDocument();
    });

    expect(onMessageSent).toHaveBeenCalledWith('Hello');
    expect(onResponseReceived).toHaveBeenCalledWith('Hello! How can I help you?');
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

    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockErrorResponse,
    });

    const onError = jest.fn();

    render(<GoChatbot onError={onError} />);

    // Open chat and send message
    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    const input = screen.getByPlaceholderText('Type your message...');
    const sendButton = screen.getByText('Send');

    await userEvent.type(input, 'Hello');
    await userEvent.click(sendButton);

    // Wait for error message
    await waitFor(() => {
      expect(screen.getByText('Error: Server error')).toBeInTheDocument();
    });

    expect(onError).toHaveBeenCalledWith('Server error');
  });

  it('handles network errors', async () => {
    (fetch as jest.Mock).mockRejectedValueOnce(new Error('Network error'));

    const onError = jest.fn();

    render(<GoChatbot onError={onError} />);

    // Open chat and send message
    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    const input = screen.getByPlaceholderText('Type your message...');
    const sendButton = screen.getByText('Send');

    await userEvent.type(input, 'Hello');
    await userEvent.click(sendButton);

    // Wait for error message
    await waitFor(() => {
      expect(screen.getByText('Error: Network error')).toBeInTheDocument();
    });

    expect(onError).toHaveBeenCalledWith('Network error');
  });

  it('sends message on Enter key press', async () => {
    const mockResponse = {
      success: true,
      response: 'Response'
    };

    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse,
    });

    render(<GoChatbot />);

    // Open chat
    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    // Type message and press Enter
    const input = screen.getByPlaceholderText('Type your message...');
    await userEvent.type(input, 'Hello{enter}');

    await waitFor(() => {
      expect(fetch).toHaveBeenCalled();
    });
  });

  it('does not send empty messages', async () => {
    render(<GoChatbot />);

    // Open chat
    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    // Try to send empty message
    const sendButton = screen.getByText('Send');
    await userEvent.click(sendButton);

    expect(fetch).not.toHaveBeenCalled();
  });

  it('displays initial messages', () => {
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

    render(<GoChatbot initialMessages={initialMessages} />);

    // Open chat
    const chatButton = screen.getByRole('button');
    fireEvent.click(chatButton);

    expect(screen.getByText('Welcome!')).toBeInTheDocument();
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });

  it('shows typing indicator when loading', async () => {
    (fetch as jest.Mock).mockImplementation(() =>
      new Promise(resolve => setTimeout(() => resolve({
        ok: true,
        json: async () => ({ success: true, response: 'Response' }),
      }), 100))
    );

    render(<GoChatbot showTypingIndicator={true} />);

    // Open chat and send message
    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    const input = screen.getByPlaceholderText('Type your message...');
    const sendButton = screen.getByText('Send');

    await userEvent.type(input, 'Hello');
    await userEvent.click(sendButton);

    expect(screen.getByText('Bot is typing')).toBeInTheDocument();
  });

  it('can be disabled', () => {
    render(<GoChatbot disabled={true} />);

    const chatButton = screen.getByRole('button');
    expect(chatButton).toBeDisabled();
  });

  it('uses custom API endpoint', async () => {
    const mockResponse = {
      success: true,
      response: 'Response'
    };

    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse,
    });

    render(<GoChatbot apiEndpoint="/custom/chat" />);

    // Open chat and send message
    const chatButton = screen.getByRole('button');
    await userEvent.click(chatButton);

    const input = screen.getByPlaceholderText('Type your message...');
    const sendButton = screen.getByText('Send');

    await userEvent.type(input, 'Hello');
    await userEvent.click(sendButton);

    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith('/custom/chat', expect.any(Object));
    });
  });

  it('renders custom header', () => {
    const customHeader = <div>Custom Support</div>;

    render(<GoChatbot header={customHeader} />);

    // Open chat
    const chatButton = screen.getByRole('button');
    fireEvent.click(chatButton);

    expect(screen.getByText('Custom Support')).toBeInTheDocument();
  });

  it('applies custom className and styles', () => {
    render(
      <GoChatbot
        className="custom-chat"
        style={{ backgroundColor: 'red' }}
      />
    );

    // Open chat
    const chatButton = screen.getByRole('button');
    fireEvent.click(chatButton);

    const chatWindow = document.querySelector('.go-chatbot');
    expect(chatWindow).toHaveClass('custom-chat');
    expect(chatWindow).toHaveStyle('background-color: red');
  });
});
