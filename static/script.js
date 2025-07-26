class ChatApp {
    constructor() {
        this.userId = this.generateUserId();
        this.isTyping = false;
        this.init();
    }

    init() {
        this.bindElements();
        this.bindEvents();
        this.loadTheme();
        this.loadHistory();
    }

    bindElements() {
        this.elements = {
            messageInput: document.getElementById('message-input'),
            sendButton: document.getElementById('send-button'),
            chatMessages: document.getElementById('chat-messages'),
            typingIndicator: document.getElementById('typing-indicator'),
            themeToggle: document.getElementById('theme-toggle'),
            clearHistory: document.getElementById('clear-history'),
            charCount: document.getElementById('char-count'),
            errorToast: document.getElementById('error-toast'),
            welcomeMessage: document.querySelector('.welcome-message')
        };
    }

    bindEvents() {
        // Send button click
        this.elements.sendButton.addEventListener('click', () => this.sendMessage());

        // Input events
        this.elements.messageInput.addEventListener('input', () => this.handleInput());
        this.elements.messageInput.addEventListener('keydown', (e) => this.handleKeydown(e));

        // Theme toggle
        this.elements.themeToggle.addEventListener('click', () => this.toggleTheme());

        // Clear history
        this.elements.clearHistory.addEventListener('click', () => this.clearHistory());

        // Quick action buttons
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('quick-btn')) {
                const message = e.target.dataset.message;
                this.elements.messageInput.value = message;
                this.handleInput();
                this.sendMessage();
            }
        });

        // Error toast close
        document.querySelector('.error-close').addEventListener('click', () => {
            this.hideError();
        });

        // Auto-hide error toast
        setTimeout(() => {
            this.hideError();
        }, 5000);
    }

    generateUserId() {
        let userId = localStorage.getItem('chat-user-id');
        if (!userId) {
            userId = 'user_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
            localStorage.setItem('chat-user-id', userId);
        }
        return userId;
    }

    handleInput() {
        const input = this.elements.messageInput;
        const text = input.value.trim();
        const charCount = input.value.length;
        
        // Update character count
        this.elements.charCount.textContent = charCount;
        
        // Update send button state
        this.elements.sendButton.disabled = !text || this.isTyping;
        
        // Auto-resize textarea
        input.style.height = 'auto';
        input.style.height = Math.min(input.scrollHeight, 120) + 'px';
    }

    handleKeydown(event) {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            this.sendMessage();
        }
    }

    async sendMessage() {
        const message = this.elements.messageInput.value.trim();
        if (!message || this.isTyping) return;

        // Hide welcome message if it exists
        if (this.elements.welcomeMessage) {
            this.elements.welcomeMessage.style.display = 'none';
        }

        // Add user message to chat
        this.addMessage(message, 'user');
        
        // Clear input
        this.elements.messageInput.value = '';
        this.handleInput();
        
        // Show typing indicator
        this.showTyping();
        
        try {
            const response = await this.callChatAPI(message);
            this.hideTyping();
            
            if (response.success) {
                this.addMessage(response.message, 'assistant');
            } else {
                this.showError('发送消息失败：' + (response.error || '未知错误'));
            }
        } catch (error) {
            this.hideTyping();
            this.showError('网络错误：' + error.message);
        }
    }

    async callChatAPI(message) {
        const response = await fetch('/api/chat', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                message: message,
                user_id: this.userId
            })
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    }

    addMessage(content, role) {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${role}`;
        
        const timestamp = new Date().toLocaleTimeString('zh-CN', {
            hour: '2-digit',
            minute: '2-digit'
        });
        
        messageDiv.innerHTML = `
            <div class="message-content">
                ${this.formatMessage(content)}
                <div class="message-time">${timestamp}</div>
            </div>
        `;
        
        this.elements.chatMessages.appendChild(messageDiv);
        this.scrollToBottom();
    }

    formatMessage(content) {
        // Simple markdown-style formatting
        return content
            .replace(/\n/g, '<br>')
            .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
            .replace(/\*(.*?)\*/g, '<em>$1</em>')
            .replace(/`(.*?)`/g, '<code>$1</code>');
    }

    showTyping() {
        this.isTyping = true;
        this.elements.sendButton.disabled = true;
        this.elements.typingIndicator.style.display = 'flex';
        this.scrollToBottom();
    }

    hideTyping() {
        this.isTyping = false;
        this.elements.typingIndicator.style.display = 'none';
        this.handleInput(); // Re-enable send button if input has text
    }

    scrollToBottom() {
        setTimeout(() => {
            this.elements.chatMessages.scrollTop = this.elements.chatMessages.scrollHeight;
        }, 100);
    }

    showError(message) {
        const errorToast = this.elements.errorToast;
        const errorMessage = errorToast.querySelector('.error-message');
        
        errorMessage.textContent = message;
        errorToast.style.display = 'flex';
        
        // Auto-hide after 5 seconds
        setTimeout(() => {
            this.hideError();
        }, 5000);
    }

    hideError() {
        this.elements.errorToast.style.display = 'none';
    }

    toggleTheme() {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('chat-theme', newTheme);
        
        // Update theme toggle icon
        const themeIcon = this.elements.themeToggle.querySelector('.theme-icon');
        themeIcon.textContent = newTheme === 'dark' ? '☀️' : '🌙';
    }

    loadTheme() {
        const savedTheme = localStorage.getItem('chat-theme') || 'light';
        document.documentElement.setAttribute('data-theme', savedTheme);
        
        // Update theme toggle icon
        const themeIcon = this.elements.themeToggle.querySelector('.theme-icon');
        themeIcon.textContent = savedTheme === 'dark' ? '☀️' : '🌙';
    }

    async loadHistory() {
        try {
            const response = await fetch(`/api/history?user_id=${this.userId}`);
            if (response.ok) {
                const data = await response.json();
                if (data.success && data.messages && data.messages.length > 0) {
                    // Hide welcome message
                    if (this.elements.welcomeMessage) {
                        this.elements.welcomeMessage.style.display = 'none';
                    }
                    
                    // Add historical messages
                    data.messages.forEach(msg => {
                        this.addMessage(msg.content, msg.role);
                    });
                }
            }
        } catch (error) {
            console.warn('Failed to load chat history:', error);
        }
    }

    async clearHistory() {
        if (!confirm('确定要清空所有对话记录吗？此操作无法撤销。')) {
            return;
        }
        
        try {
            const response = await fetch(`/api/history?user_id=${this.userId}`, {
                method: 'DELETE'
            });
            
            if (response.ok) {
                // Clear chat messages
                this.elements.chatMessages.innerHTML = `
                    <div class="welcome-message">
                        <div class="welcome-icon">👋</div>
                        <h2>欢迎使用 Chat MVP</h2>
                        <p>我是基于 OpenAI GPT-4o 的智能助手，很高兴为您服务！</p>
                        <div class="quick-actions">
                            <button class="quick-btn" data-message="你好，请介绍一下你自己">自我介绍</button>
                            <button class="quick-btn" data-message="请帮我解释一下人工智能">AI 知识</button>
                            <button class="quick-btn" data-message="给我讲一个有趣的故事">讲个故事</button>
                        </div>
                    </div>
                `;
                
                // Re-bind welcome message element
                this.elements.welcomeMessage = document.querySelector('.welcome-message');
                
                this.showError('对话记录已清空');
            } else {
                this.showError('清空对话记录失败');
            }
        } catch (error) {
            this.showError('网络错误：' + error.message);
        }
    }
}

// Initialize the chat app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new ChatApp();
});