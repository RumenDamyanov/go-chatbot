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

export interface GoChatbotConfig {
  apiEndpoint?: string;
  placeholder?: string;
  maxHeight?: string;
  className?: string;
  style?: Record<string, any>;
  initialMessages?: ChatMessage[];
  showTypingIndicator?: boolean;
  disabled?: boolean;
  position?: 'bottom-right' | 'bottom-left' | 'top-right' | 'top-left';
  theme?: 'light' | 'dark';
}

export type EventHandler<T = any> = (data: T) => void;

export declare class GoChatbot {
  constructor(options?: GoChatbotConfig);

  /**
   * Open the chat window
   */
  open(): void;

  /**
   * Close the chat window
   */
  close(): void;

  /**
   * Toggle the chat window
   */
  toggle(): void;

  /**
   * Send a message programmatically
   */
  sendMessage(): Promise<void>;

  /**
   * Add a message to the chat
   */
  addMessage(message: Partial<ChatMessage>): void;

  /**
   * Clear all messages
   */
  clearMessages(): void;

  /**
   * Enable or disable the chatbot
   */
  setDisabled(disabled: boolean): void;

  /**
   * Update the chatbot configuration
   */
  updateConfig(config: Partial<GoChatbotConfig>): void;

  /**
   * Add an event listener
   */
  on(event: 'messageSent', handler: EventHandler<string>): void;
  on(event: 'responseReceived', handler: EventHandler<string>): void;
  on(event: 'error', handler: EventHandler<string>): void;
  on(event: 'opened', handler: EventHandler): void;
  on(event: 'closed', handler: EventHandler): void;

  /**
   * Remove an event listener
   */
  off(event: string, handler: EventHandler): void;

  /**
   * Destroy the chatbot instance
   */
  destroy(): void;
}

declare global {
  interface Window {
    GoChatbot: typeof GoChatbot;
  }
}

export default GoChatbot;
