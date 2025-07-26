class ChatApp {
    constructor() {
        this.currentUserId = '';
        this.isConnected = false;
        this.messageHistory = [];
        this.initializeElements();
        this.attachEventListeners();
        this.loadChatHistory();
    }

    initializeElements() {
        this.userIdInput = document.getElementById('userId');
        this.connectBtn = document.getElementById('connectBtn');
        this.chatContainer = document.getElementById('chatContainer');
        this.chatMessages = document.getElementById('chatMessages');
        this.messageInput = document.getElementById('messageInput');
        this.sendBtn = document.getElementById('sendBtn');
        this.loading = document.getElementById('loading');
        this.charCount = document.getElementById('charCount');
        this.status = document.getElementById('status');
    }

    attachEventListeners() {
        // 连接按钮事件
        this.connectBtn.addEventListener('click', () => this.handleConnect());
        
        // 发送按钮事件
        this.sendBtn.addEventListener('click', () => this.handleSendMessage());
        
        // 输入框回车事件
        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.handleSendMessage();
            }
        });

        // 用户ID输入框回车事件
        this.userIdInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                this.handleConnect();
            }
        });

        // 字符计数
        this.messageInput.addEventListener('input', () => {
            const length = this.messageInput.value.length;
            this.charCount.textContent = `${length}/500`;
            
            if (length > 450) {
                this.charCount.style.color = '#f44336';
            } else {
                this.charCount.style.color = '#666';
            }
        });
    }

    handleConnect() {
        const userId = this.userIdInput.value.trim();
        
        if (!userId) {
            this.showError('请输入用户ID');
            return;
        }

        this.currentUserId = userId;
        this.isConnected = true;
        
        // 显示聊天界面
        this.chatContainer.style.display = 'flex';
        this.connectBtn.textContent = '重新连接';
        this.status.textContent = `已连接 - 用户: ${userId}`;
        this.status.className = 'status-connected';
        
        // 加载历史消息
        this.loadUserHistory();
        
        // 聚焦到消息输入框
        this.messageInput.focus();
    }

    async handleSendMessage() {
        const message = this.messageInput.value.trim();
        
        if (!message) {
            return;
        }

        if (!this.isConnected) {
            this.showError('请先连接用户ID');
            return;
        }

        // 清空输入框
        this.messageInput.value = '';
        this.charCount.textContent = '0/500';
        
        // 显示用户消息
        this.displayMessage('user', message);
        
        // 禁用发送按钮
        this.sendBtn.disabled = true;
        this.sendBtn.textContent = '发送中...';
        this.showLoading();

        try {
            const response = await this.sendMessageToAPI(message);
            
            if (response.status === 'success') {
                // 显示AI响应
                this.displayMessage('assistant', response.data.reply);
                this.saveToHistory('user', message);
                this.saveToHistory('assistant', response.data.reply);
            } else {
                this.showError(response.message || '发送消息失败');
            }
        } catch (error) {
            console.error('发送消息错误:', error);
            this.showError('网络错误，请检查连接后重试');
        } finally {
            // 恢复发送按钮
            this.sendBtn.disabled = false;
            this.sendBtn.textContent = '发送';
            this.hideLoading();
            this.messageInput.focus();
        }
    }

    async sendMessageToAPI(message) {
        const response = await fetch('/chat', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                message: message,
                userId: this.currentUserId
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.json();
    }

    displayMessage(role, content) {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${role}`;
        
        const bubbleDiv = document.createElement('div');
        bubbleDiv.className = 'message-bubble';
        bubbleDiv.textContent = content;
        
        const timeDiv = document.createElement('div');
        timeDiv.className = 'message-time';
        timeDiv.textContent = this.formatTime(new Date());
        
        messageDiv.appendChild(bubbleDiv);
        messageDiv.appendChild(timeDiv);
        this.chatMessages.appendChild(messageDiv);
        
        // 滚动到底部
        this.scrollToBottom();
    }

    formatTime(date) {
        return date.toLocaleTimeString('zh-CN', {
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    scrollToBottom() {
        setTimeout(() => {
            this.chatMessages.scrollTop = this.chatMessages.scrollHeight;
        }, 100);
    }

    showLoading() {
        this.loading.style.display = 'flex';
    }

    hideLoading() {
        this.loading.style.display = 'none';
    }

    showError(message) {
        // 创建错误提示
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-message';
        errorDiv.textContent = message;
        
        // 插入到聊天消息区域
        this.chatMessages.appendChild(errorDiv);
        this.scrollToBottom();
        
        // 3秒后自动删除
        setTimeout(() => {
            if (errorDiv.parentNode) {
                errorDiv.parentNode.removeChild(errorDiv);
            }
        }, 3000);
    }

    // 保存消息到本地存储
    saveToHistory(role, content) {
        const message = {
            role: role,
            content: content,
            timestamp: new Date().toISOString(),
            userId: this.currentUserId
        };
        
        this.messageHistory.push(message);
        this.saveChatHistory();
    }

    // 保存聊天历史到本地存储
    saveChatHistory() {
        const historyKey = `chat_history_${this.currentUserId}`;
        localStorage.setItem(historyKey, JSON.stringify(this.messageHistory));
    }

    // 加载聊天历史
    loadChatHistory() {
        if (!this.currentUserId) return;
        
        const historyKey = `chat_history_${this.currentUserId}`;
        const history = localStorage.getItem(historyKey);
        
        if (history) {
            try {
                this.messageHistory = JSON.parse(history);
                // 显示历史消息
                this.messageHistory.forEach(msg => {
                    this.displayMessage(msg.role, msg.content);
                });
            } catch (error) {
                console.error('加载聊天历史失败:', error);
                this.messageHistory = [];
            }
        }
    }

    // 从服务器加载用户历史（可选功能）
    async loadUserHistory() {
        try {
            const response = await fetch(`/api/history?userId=${this.currentUserId}`);
            if (response.ok) {
                const result = await response.json();
                if (result.status === 'success' && Array.isArray(result.data)) {
                    // 清空当前显示的消息
                    this.chatMessages.innerHTML = '';
                    this.messageHistory = [];
                    
                    // 显示服务器历史消息
                    result.data.forEach(msg => {
                        this.displayMessage(msg.role, msg.content);
                        this.messageHistory.push(msg);
                    });
                }
            }
        } catch (error) {
            console.error('加载服务器历史失败:', error);
            // 使用本地历史作为后备
            this.loadChatHistory();
        }
    }
}

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new ChatApp();
});