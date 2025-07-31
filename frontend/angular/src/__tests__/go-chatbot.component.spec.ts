import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { FormsModule } from '@angular/forms';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import { GoChatbotComponent, ChatMessage, ChatResponse } from '../go-chatbot.component';

describe('GoChatbotComponent', () => {
  let component: GoChatbotComponent;
  let fixture: ComponentFixture<GoChatbotComponent>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [GoChatbotComponent],
      imports: [HttpClientTestingModule, FormsModule],
    }).compileComponents();

    fixture = TestBed.createComponent(GoChatbotComponent);
    component = fixture.componentInstance;
    httpMock = TestBed.inject(HttpTestingController);
    fixture.detectChanges();
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should render chat button', () => {
    const button = fixture.debugElement.query(By.css('.go-chatbot-button'));
    expect(button).toBeTruthy();
    expect(button.nativeElement.textContent.trim()).toBe('💬');
  });

  it('should toggle chat window when button is clicked', () => {
    const button = fixture.debugElement.query(By.css('.go-chatbot-button'));

    // Initially closed
    expect(component.isOpen).toBe(false);
    expect(fixture.debugElement.query(By.css('.go-chatbot-window'))).toBeFalsy();

    // Click to open
    button.nativeElement.click();
    fixture.detectChanges();

    expect(component.isOpen).toBe(true);
    expect(fixture.debugElement.query(By.css('.go-chatbot-window'))).toBeTruthy();
    expect(button.nativeElement.textContent.trim()).toBe('×');
  });

  it('should close chat window when close button is clicked', () => {
    component.isOpen = true;
    fixture.detectChanges();

    const closeButton = fixture.debugElement.query(By.css('.go-chatbot-close'));
    closeButton.nativeElement.click();
    fixture.detectChanges();

    expect(component.isOpen).toBe(false);
  });

  it('should display initial messages', () => {
    const initialMessages: ChatMessage[] = [
      {
        id: '1',
        text: 'Hello!',
        sender: 'user',
        timestamp: new Date()
      },
      {
        id: '2',
        text: 'Hi there!',
        sender: 'bot',
        timestamp: new Date()
      }
    ];

    component.initialMessages = initialMessages;
    component.ngOnInit();
    component.isOpen = true;
    fixture.detectChanges();

    const messageElements = fixture.debugElement.queryAll(By.css('.go-chatbot-message'));
    expect(messageElements.length).toBe(2);
    expect(messageElements[0].nativeElement.textContent).toContain('Hello!');
    expect(messageElements[1].nativeElement.textContent).toContain('Hi there!');
  });

  it('should send message when send button is clicked', () => {
    component.isOpen = true;
    fixture.detectChanges();

    spyOn(component.messageSent, 'emit');

    const input = fixture.debugElement.query(By.css('.go-chatbot-input'));
    const sendButton = fixture.debugElement.query(By.css('.go-chatbot-send'));

    input.nativeElement.value = 'Test message';
    input.nativeElement.dispatchEvent(new Event('input'));
    component.inputValue = 'Test message';
    fixture.detectChanges();

    sendButton.nativeElement.click();

    expect(component.messageSent.emit).toHaveBeenCalledWith('Test message');
    expect(component.messages.length).toBe(1);
    expect(component.messages[0].text).toBe('Test message');
    expect(component.messages[0].sender).toBe('user');
  });

  it('should send message when Enter key is pressed', () => {
    component.isOpen = true;
    fixture.detectChanges();

    spyOn(component.messageSent, 'emit');

    const input = fixture.debugElement.query(By.css('.go-chatbot-input'));

    input.nativeElement.value = 'Test message';
    component.inputValue = 'Test message';

    const enterEvent = new KeyboardEvent('keydown', { key: 'Enter' });
    input.nativeElement.dispatchEvent(enterEvent);
    component.sendMessage();

    expect(component.messageSent.emit).toHaveBeenCalledWith('Test message');
    expect(component.messages.length).toBe(1);
  });

  it('should handle successful API response', () => {
    spyOn(component.responseReceived, 'emit');

    component.isOpen = true;
    component.inputValue = 'Hello';
    component.sendMessage();

    const req = httpMock.expectOne('/chat/');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual({ message: 'Hello' });

    const mockResponse: ChatResponse = {
      success: true,
      response: 'Hi there!'
    };

    req.flush(mockResponse);
    fixture.detectChanges();

    expect(component.messages.length).toBe(2);
    expect(component.messages[1].text).toBe('Hi there!');
    expect(component.messages[1].sender).toBe('bot');
    expect(component.responseReceived.emit).toHaveBeenCalledWith('Hi there!');
    expect(component.isLoading).toBe(false);
  });

  it('should handle API error response', () => {
    spyOn(component.error, 'emit');

    component.isOpen = true;
    component.inputValue = 'Hello';
    component.sendMessage();

    const req = httpMock.expectOne('/chat/');
    const mockResponse: ChatResponse = {
      success: false,
      error: 'Server error'
    };

    req.flush(mockResponse);
    fixture.detectChanges();

    expect(component.messages.length).toBe(2);
    expect(component.messages[1].text).toBe('Error: Server error');
    expect(component.messages[1].sender).toBe('bot');
    expect(component.error.emit).toHaveBeenCalledWith('Server error');
    expect(component.isLoading).toBe(false);
  });

  it('should handle network error', () => {
    spyOn(component.error, 'emit');

    component.isOpen = true;
    component.inputValue = 'Hello';
    component.sendMessage();

    const req = httpMock.expectOne('/chat/');
    req.error(new ErrorEvent('Network error'));
    fixture.detectChanges();

    expect(component.messages.length).toBe(2);
    expect(component.messages[1].text).toContain('Error:');
    expect(component.messages[1].sender).toBe('bot');
    expect(component.error.emit).toHaveBeenCalled();
    expect(component.isLoading).toBe(false);
  });

  it('should disable input when loading', () => {
    component.isOpen = true;
    component.isLoading = true;
    fixture.detectChanges();

    const input = fixture.debugElement.query(By.css('.go-chatbot-input'));
    const sendButton = fixture.debugElement.query(By.css('.go-chatbot-send'));

    expect(input.nativeElement.disabled).toBe(true);
    expect(sendButton.nativeElement.disabled).toBe(true);
  });

  it('should disable component when disabled prop is true', () => {
    component.disabled = true;
    component.isOpen = true;
    fixture.detectChanges();

    const button = fixture.debugElement.query(By.css('.go-chatbot-button'));
    const input = fixture.debugElement.query(By.css('.go-chatbot-input'));
    const sendButton = fixture.debugElement.query(By.css('.go-chatbot-send'));

    expect(button.nativeElement.disabled).toBe(true);
    expect(input.nativeElement.disabled).toBe(true);
    expect(sendButton.nativeElement.disabled).toBe(true);
  });

  it('should not send empty messages', () => {
    component.isOpen = true;
    component.inputValue = '   ';
    fixture.detectChanges();

    const sendButton = fixture.debugElement.query(By.css('.go-chatbot-send'));
    expect(sendButton.nativeElement.disabled).toBe(true);

    component.sendMessage();
    expect(component.messages.length).toBe(0);
  });

  it('should show typing indicator when loading', () => {
    component.isOpen = true;
    component.isLoading = true;
    component.showTypingIndicator = true;
    fixture.detectChanges();

    const typingIndicator = fixture.debugElement.query(By.css('.go-chatbot-typing'));
    expect(typingIndicator).toBeTruthy();
    expect(typingIndicator.nativeElement.textContent).toContain('Bot is typing');
  });

  it('should use custom placeholder', () => {
    component.placeholder = 'Custom placeholder';
    component.isOpen = true;
    fixture.detectChanges();

    const input = fixture.debugElement.query(By.css('.go-chatbot-input'));
    expect(input.nativeElement.placeholder).toBe('Custom placeholder');
  });

  it('should use custom API endpoint', () => {
    component.apiEndpoint = '/custom/chat';
    component.isOpen = true;
    component.inputValue = 'Hello';
    component.sendMessage();

    const req = httpMock.expectOne('/custom/chat');
    expect(req.request.method).toBe('POST');
    req.flush({ success: true, response: 'Response' });
  });

  it('should apply custom styles', () => {
    component.style = { width: '400px' };
    component.isOpen = true;
    fixture.detectChanges();

    const chatWindow = fixture.debugElement.query(By.css('.go-chatbot-window'));
    expect(chatWindow.nativeElement.style.width).toBe('400px');
  });

  it('should apply custom className', () => {
    component.className = 'custom-class';
    component.isOpen = true;
    fixture.detectChanges();

    const chatWindow = fixture.debugElement.query(By.css('.go-chatbot-window'));
    expect(chatWindow.nativeElement.classList.contains('custom-class')).toBe(true);
  });

  it('should track messages by ID', () => {
    const message: ChatMessage = {
      id: 'test-id',
      text: 'Test',
      sender: 'user',
      timestamp: new Date()
    };

    const result = component.trackMessage(0, message);
    expect(result).toBe('test-id');
  });

  it('should format time correctly', () => {
    const date = new Date('2023-01-01T12:30:00');
    const formatted = component.formatTime(date);
    expect(formatted).toMatch(/12:30/);
  });

  it('should generate unique IDs', () => {
    const id1 = component.generateId();
    const id2 = component.generateId();

    expect(id1).toContain('msg-');
    expect(id2).toContain('msg-');
    expect(id1).not.toBe(id2);
  });
});
