import { Component, Input, Output, EventEmitter, OnInit, OnDestroy, ElementRef, ViewChild, AfterViewChecked } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subscription } from 'rxjs';

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

@Component({
  selector: 'go-chatbot',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <!-- Chat Button -->
    <button
      (click)="toggleChat()"
      [disabled]="disabled"
      class="go-chatbot-button"
      [style]="chatButtonStyle"
    >
      {{ isOpen ? '×' : '💬' }}
    </button>

    <!-- Chat Window -->
    <div
      *ngIf="isOpen"
      [class]="'go-chatbot-window ' + className"
      [ngStyle]="{ ...chatWindowStyle, ...style }"
    >
      <!-- Header -->
      <div class="go-chatbot-header">
        <ng-content select="[slot=header]">
          <span class="go-chatbot-title">Chat Support</span>
        </ng-content>
        <button (click)="closeChat()" class="go-chatbot-close">×</button>
      </div>

      <!-- Messages -->
      <div class="go-chatbot-messages" [style.max-height]="maxHeight" #messagesContainer>
        <div
          *ngFor="let message of messages; trackBy: trackMessage"
          [class]="'go-chatbot-message go-chatbot-message--' + message.sender"
        >
          <div class="go-chatbot-message-content">
            <div class="go-chatbot-message-text">{{ message.text }}</div>
            <div class="go-chatbot-message-time">
              {{ formatTime(message.timestamp) }}
            </div>
          </div>
        </div>

        <!-- Typing Indicator -->
        <div
          *ngIf="isLoading && showTypingIndicator"
          class="go-chatbot-message go-chatbot-message--bot"
        >
          <div class="go-chatbot-message-content">
            <div class="go-chatbot-typing">
              <span>Bot is typing</span>
              <div class="go-chatbot-dots">
                <div class="go-chatbot-dot"></div>
                <div class="go-chatbot-dot"></div>
                <div class="go-chatbot-dot"></div>
              </div>
            </div>
          </div>
        </div>

        <div #messagesEnd></div>
      </div>

      <!-- Input -->
      <div class="go-chatbot-input-container">
        <input
          #messageInput
          [(ngModel)]="inputValue"
          (keydown.enter)="sendMessage()"
          [disabled]="isLoading || disabled"
          [placeholder]="placeholder"
          class="go-chatbot-input"
          type="text"
        />
        <button
          (click)="sendMessage()"
          [disabled]="!inputValue?.trim() || isLoading || disabled"
          class="go-chatbot-send"
        >
          {{ isLoading ? '...' : 'Send' }}
        </button>
      </div>
    </div>
  `,
  styles: [`
    .go-chatbot-button {
      position: fixed;
      bottom: 20px;
      right: 20px;
      width: 60px;
      height: 60px;
      border-radius: 50%;
      background-color: #007bff;
      color: white;
      border: none;
      cursor: pointer;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
      font-size: 24px;
      z-index: 1000;
      transition: all 0.3s ease;
    }

    .go-chatbot-button:disabled {
      cursor: not-allowed;
      opacity: 0.6;
    }

    .go-chatbot-window {
      position: fixed;
      bottom: 90px;
      right: 20px;
      width: 350px;
      background-color: white;
      border-radius: 12px;
      box-shadow: 0 8px 32px rgba(0,0,0,0.12);
      display: flex;
      flex-direction: column;
      z-index: 999;
      overflow: hidden;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
    }

    .go-chatbot-header {
      padding: 16px;
      background-color: #007bff;
      color: white;
      border-radius: 12px 12px 0 0;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .go-chatbot-title {
      font-weight: bold;
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
    }

    .go-chatbot-messages {
      flex: 1;
      overflow-y: auto;
      padding: 16px;
      display: flex;
      flex-direction: column;
      gap: 12px;
    }

    .go-chatbot-message {
      display: flex;
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
      transition: opacity 0.2s ease;
    }

    .go-chatbot-send:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }

    .go-chatbot-send:hover:not(:disabled) {
      background-color: #0056b3;
    }
  `]
})
export class GoChatbotComponent implements OnInit, OnDestroy, AfterViewChecked {
  @Input() apiEndpoint: string = '/chat/';
  @Input() className: string = '';
  @Input() placeholder: string = 'Type your message...';
  @Input() maxHeight: string = '400px';
  @Input() style: { [key: string]: any } = {};
  @Input() initialMessages: ChatMessage[] = [];
  @Input() showTypingIndicator: boolean = true;
  @Input() disabled: boolean = false;

  @Output() messageSent = new EventEmitter<string>();
  @Output() responseReceived = new EventEmitter<string>();
  @Output() error = new EventEmitter<string>();

  @ViewChild('messagesEnd') messagesEnd!: ElementRef;
  @ViewChild('messageInput') messageInput!: ElementRef;
  @ViewChild('messagesContainer') messagesContainer!: ElementRef;

  messages: ChatMessage[] = [];
  inputValue: string = '';
  isLoading: boolean = false;
  isOpen: boolean = false;
  private subscription: Subscription = new Subscription();
  private shouldScrollToBottom: boolean = false;

  chatButtonStyle = {
    position: 'fixed',
    bottom: '20px',
    right: '20px',
    width: '60px',
    height: '60px',
    borderRadius: '50%',
    backgroundColor: '#007bff',
    color: 'white',
    border: 'none',
    cursor: 'pointer',
    boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
    fontSize: '24px',
    zIndex: 1000,
    transition: 'all 0.3s ease'
  };

  chatWindowStyle = {
    position: 'fixed',
    bottom: '90px',
    right: '20px',
    width: '350px',
    backgroundColor: 'white',
    borderRadius: '12px',
    boxShadow: '0 8px 32px rgba(0,0,0,0.12)',
    display: 'flex',
    flexDirection: 'column',
    zIndex: 999,
    overflow: 'hidden'
  };

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.messages = [...this.initialMessages];
    this.shouldScrollToBottom = true;
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }

  ngAfterViewChecked(): void {
    if (this.shouldScrollToBottom) {
      this.scrollToBottom();
      this.shouldScrollToBottom = false;
    }
  }

  generateId(): string {
    return `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  formatTime(date: Date): string {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  trackMessage(index: number, message: ChatMessage): string {
    return message.id;
  }

  scrollToBottom(): void {
    if (this.messagesEnd) {
      this.messagesEnd.nativeElement.scrollIntoView({ behavior: 'smooth' });
    }
  }

  toggleChat(): void {
    this.isOpen = !this.isOpen;
    if (this.isOpen) {
      setTimeout(() => {
        this.messageInput?.nativeElement?.focus();
      }, 0);
    }
  }

  closeChat(): void {
    this.isOpen = false;
  }

  sendMessage(): void {
    if (!this.inputValue?.trim() || this.isLoading || this.disabled) return;

    const userMessage: ChatMessage = {
      id: this.generateId(),
      text: this.inputValue.trim(),
      sender: 'user',
      timestamp: new Date(),
    };

    this.messages.push(userMessage);
    const messageText = this.inputValue.trim();
    this.inputValue = '';
    this.isLoading = true;
    this.shouldScrollToBottom = true;

    this.messageSent.emit(messageText);

    const subscription = this.http.post<ChatResponse>(this.apiEndpoint, { message: messageText })
      .subscribe({
        next: (data) => {
          if (!data.success) {
            throw new Error(data.error || 'Unknown error occurred');
          }

          const botMessage: ChatMessage = {
            id: this.generateId(),
            text: data.response!,
            sender: 'bot',
            timestamp: new Date(),
          };

          this.messages.push(botMessage);
          this.responseReceived.emit(data.response!);
          this.shouldScrollToBottom = true;
        },
        error: (error) => {
          const errorMessage = error.message || 'Failed to send message';
          const botMessage: ChatMessage = {
            id: this.generateId(),
            text: `Error: ${errorMessage}`,
            sender: 'bot',
            timestamp: new Date(),
          };

          this.messages.push(botMessage);
          this.error.emit(errorMessage);
          this.shouldScrollToBottom = true;
        },
        complete: () => {
          this.isLoading = false;
        }
      });

    this.subscription.add(subscription);
  }
}
