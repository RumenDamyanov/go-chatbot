<template>
  <div>
    <!-- Chat Button -->
    <button
      @click="toggleChat"
      :disabled="disabled"
      class="go-chatbot-button"
      :style="chatButtonStyle"
    >
      {{ isOpen ? '×' : '💬' }}
    </button>

    <!-- Chat Window -->
    <div
      v-if="isOpen"
      :class="['go-chatbot-window', className]"
      :style="{ ...chatWindowStyle, ...style }"
    >
      <!-- Header -->
      <div class="go-chatbot-header">
        <slot name="header">
          <span class="go-chatbot-title">Chat Support</span>
        </slot>
        <button @click="closeChat" class="go-chatbot-close">×</button>
      </div>

      <!-- Messages -->
      <div class="go-chatbot-messages" :style="{ maxHeight }">
        <div
          v-for="message in messages"
          :key="message.id"
          :class="['go-chatbot-message', `go-chatbot-message--${message.sender}`]"
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
          v-if="isLoading && showTypingIndicator"
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

        <div ref="messagesEnd"></div>
      </div>

      <!-- Input -->
      <div class="go-chatbot-input-container">
        <input
          ref="messageInput"
          v-model="inputValue"
          @keydown.enter.prevent="sendMessage"
          :disabled="isLoading || disabled"
          :placeholder="placeholder"
          class="go-chatbot-input"
        />
        <button
          @click="sendMessage"
          :disabled="!inputValue.trim() || isLoading || disabled"
          class="go-chatbot-send"
        >
          {{ isLoading ? '...' : 'Send' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue';

export interface ChatMessage {
  id: string;
  text: string;
  sender: 'user' | 'bot';
  timestamp: Date;
}

export interface Props {
  apiEndpoint?: string;
  className?: string;
  placeholder?: string;
  maxHeight?: string;
  style?: Record<string, any>;
  initialMessages?: ChatMessage[];
  showTypingIndicator?: boolean;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  apiEndpoint: '/chat/',
  className: '',
  placeholder: 'Type your message...',
  maxHeight: '400px',
  style: () => ({}),
  initialMessages: () => [],
  showTypingIndicator: true,
  disabled: false,
});

const emit = defineEmits<{
  messageSent: [message: string];
  responseReceived: [response: string];
  error: [error: string];
}>();

// Reactive state
const messages = ref<ChatMessage[]>([...props.initialMessages]);
const inputValue = ref('');
const isLoading = ref(false);
const isOpen = ref(false);
const messagesEnd = ref<HTMLElement>();
const messageInput = ref<HTMLInputElement>();

// Computed styles
const chatButtonStyle = computed(() => ({
  position: 'fixed',
  bottom: '20px',
  right: '20px',
  width: '60px',
  height: '60px',
  borderRadius: '50%',
  backgroundColor: '#007bff',
  color: 'white',
  border: 'none',
  cursor: props.disabled ? 'not-allowed' : 'pointer',
  boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
  fontSize: '24px',
  zIndex: 1000,
  transition: 'all 0.3s ease',
}));

const chatWindowStyle = computed(() => ({
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
  overflow: 'hidden',
}));

// Methods
const generateId = (): string => `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

const formatTime = (date: Date): string => {
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

const scrollToBottom = async (): Promise<void> => {
  await nextTick();
  messagesEnd.value?.scrollIntoView({ behavior: 'smooth' });
};

const toggleChat = (): void => {
  isOpen.value = !isOpen.value;
};

const closeChat = (): void => {
  isOpen.value = false;
};

const sendMessage = async (): Promise<void> => {
  if (!inputValue.value.trim() || isLoading.value || props.disabled) return;

  const userMessage: ChatMessage = {
    id: generateId(),
    text: inputValue.value.trim(),
    sender: 'user',
    timestamp: new Date(),
  };

  messages.value.push(userMessage);
  const messageText = inputValue.value.trim();
  inputValue.value = '';
  isLoading.value = true;

  emit('messageSent', messageText);

  try {
    const response = await fetch(props.apiEndpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message: messageText }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();

    if (!data.success) {
      throw new Error(data.error || 'Unknown error occurred');
    }

    const botMessage: ChatMessage = {
      id: generateId(),
      text: data.response,
      sender: 'bot',
      timestamp: new Date(),
    };

    messages.value.push(botMessage);
    emit('responseReceived', data.response);
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Failed to send message';
    const botMessage: ChatMessage = {
      id: generateId(),
      text: `Error: ${errorMessage}`,
      sender: 'bot',
      timestamp: new Date(),
    };

    messages.value.push(botMessage);
    emit('error', errorMessage);
  } finally {
    isLoading.value = false;
  }
};

// Watchers
watch([messages, isLoading], scrollToBottom);

watch(isOpen, async (newValue) => {
  if (newValue) {
    await nextTick();
    messageInput.value?.focus();
  }
});

// Lifecycle
onMounted(() => {
  scrollToBottom();
});
</script>

<style scoped>
.go-chatbot-window {
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
</style>
