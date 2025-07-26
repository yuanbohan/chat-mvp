class ChatApp {
    constructor() {
        this.currentUserId = this.generateUserId();
        this.isLoading = false;

        this.initializeElements();
        this.bindEvents();
        this.initializeTheme();
        this.setWelcomeTime();
        this.loadChatHistory();
    }

    // 初始化DOM元素
    initializeElements() {
        this.messagesContainer = document.getElementById('chat-messages');
        this.messageInput = document.getElementById('message-input');
        this.sendBtn = document.getElementById('send-btn');
        this.clearBtn = document.getElementById('clear-btn');
        this.themeToggle = document.getElementById('theme-toggle');
        this.loading = document.getElementById('loading');
        this.charCount = document.getElementById('char-count');
        this.statusToast = document.getElementById('status-toast');
        this.confirmModal = document.getElementById('confirm-modal');
        this.confirmClear = document.getElementById('confirm-clear');
        this.confirmCancel = document.getElementById('confirm-cancel');
    }

    // 绑定事件监听器
    bindEvents() {
        // 发送消息
        this.sendBtn.addEventListener('click', () => this.sendMessage());

        // 输入框事件
        this.messageInput.addEventListener('input', () => this.updateCharCount());
        this.messageInput.addEventListener('keydown', (e) => this.handleInputKeydown(e));

        // 清除历史
        this.clearBtn.addEventListener('click', () => this.showConfirmModal());
        this.confirmClear.addEventListener('click', () => this.clearHistory());
        this.confirmCancel.addEventListener('click', () => this.hideConfirmModal());

        // 主题切换
        this.themeToggle.addEventListener('click', () => this.toggleTheme());

        // 模态框点击外部关闭
        this.confirmModal.addEventListener('click', (e) => {
            if (e.target === this.confirmModal) {
                this.hideConfirmModal();
            }
        });

        // 自动调整输入框高度
        this.messageInput.addEventListener('input', () => this.autoResizeTextarea());
    }

    // 生成用户ID
    generateUserId() {
        return 'user_' + Math.random().toString(36).substr(2, 9) + '_' + Date.now();
    }

    // 设置欢迎消息时间
    setWelcomeTime() {
        const welcomeTime = document.getElementById('welcome-time');
        if (welcomeTime) {
            welcomeTime.textContent = this.formatTime(new Date());
        }
    }

    // 格式化时间
    formatTime(date) {
        return date.toLocaleTimeString('zh-CN', {
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    // 处理输入框按键事件
    handleInputKeydown(e) {
        if (e.key === 'Enter') {
            if (e.shiftKey) {
                // Shift+Enter 换行
                return;
            } else {
                // Enter 发送消息
                e.preventDefault();
                this.sendMessage();
            }
        }
    }

    // 自动调整文本域高度
    autoResizeTextarea() {
        const textarea = this.messageInput;
        textarea.style.height = 'auto';
        const newHeight = Math.min(textarea.scrollHeight, 120);
        textarea.style.height = newHeight + 'px';
    }

    // 更新字符计数
    updateCharCount() {
        const count = this.messageInput.value.length;
        this.charCount.textContent = count;

        // 更新发送按钮状态
        this.sendBtn.disabled = count === 0 || count > 2000 || this.isLoading;

        // 字符数超限时的样式
        if (count > 2000) {
            this.charCount.style.color = 'var(--danger-color)';
        } else {
            this.charCount.style.color = 'var(--text-secondary)';
        }
    }

    // 发送消息
    async sendMessage() {
        const message = this.messageInput.value.trim();
        if (!message || this.isLoading) return;

        // 添加用户消息到界面
        this.addMessage(message, 'user');

        // 清空输入框
        this.messageInput.value = '';
        this.updateCharCount();
        this.autoResizeTextarea();

        // 显示加载状态
        this.setLoading(true);

        try {
            const response = await fetch('/api/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    message: message,
                    user_id: this.currentUserId
                })
            });

            const data = await response.json();

            if (data.success) {
                this.addMessage(data.message, 'assistant');
            } else {
                this.showToast('发送失败: ' + (data.error || '未知错误'), 'error');
                console.error('Chat error:', data.error);
            }
        } catch (error) {
            console.error('Network error:', error);
            this.showToast('网络错误，请检查连接后重试', 'error');
        } finally {
            this.setLoading(false);
        }
    }

    // 添加消息到界面
    addMessage(content, role, timestamp = new Date()) {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${role}`;

        const avatar = role === 'user' ? '👤' : '🤖';
        const timeStr = this.formatTime(timestamp);

        messageDiv.innerHTML = `
            <div class="message-avatar">${avatar}</div>
            <div class="message-content">
                <div class="message-text">${this.formatMessageContent(content)}</div>
                <div class="message-time">${timeStr}</div>
            </div>
        `;

        this.messagesContainer.appendChild(messageDiv);
        this.scrollToBottom();
    }

    // 格式化消息内容
    formatMessageContent(content) {
        // 简单的HTML转义
        const div = document.createElement('div');
        div.textContent = content;
        let escaped = div.innerHTML;

        // 简单的Markdown支持
        escaped = escaped
            .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
            .replace(/\*(.*?)\*/g, '<em>$1</em>')
            .replace(/`(.*?)`/g, '<code>$1</code>')
            .replace(/\n/g, '<br>');

        return escaped;
    }

    // 设置加载状态
    setLoading(loading) {
        this.isLoading = loading;

        if (loading) {
            this.loading.classList.remove('hidden');
            this.sendBtn.disabled = true;
        } else {
            this.loading.classList.add('hidden');
            this.updateCharCount(); // 重新计算发送按钮状态
        }

        this.scrollToBottom();
    }

    // 滚动到底部
    scrollToBottom() {
        setTimeout(() => {
            this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
        }, 100);
    }

    // 加载聊天历史
    async loadChatHistory() {
        try {
            const response = await fetch(`/api/history?user_id=${this.currentUserId}`);
            const data = await response.json();

            if (data.success && data.history && data.history.length > 0) {
                // 清除欢迎消息
                const welcomeMsg = this.messagesContainer.querySelector('.message.assistant');
                if (welcomeMsg) {
                    welcomeMsg.remove();
                }

                // 添加历史消息
                data.history.forEach(msg => {
                    this.addMessage(msg.content, msg.role, new Date(msg.timestamp));
                });
            }
        } catch (error) {
            console.error('Failed to load chat history:', error);
        }
    }

    // 显示确认模态框
    showConfirmModal() {
        this.confirmModal.classList.remove('hidden');
        setTimeout(() => {
            this.confirmModal.classList.add('show');
        }, 10);
    }

    // 隐藏确认模态框
    hideConfirmModal() {
        this.confirmModal.classList.remove('show');
        setTimeout(() => {
            this.confirmModal.classList.add('hidden');
        }, 300);
    }

    // 清除聊天历史
    async clearHistory() {
        this.hideConfirmModal();

        try {
            const response = await fetch('/api/clear-history', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user_id: this.currentUserId
                })
            });

            const data = await response.json();

            if (data.success) {
                // 清空消息容器
                this.messagesContainer.innerHTML = '';

                // 重新添加欢迎消息
                this.addMessage(
                    '你好！我是基于 GPT-4o 的智能助手。我可以帮助你解答问题、提供建议或进行有趣的对话。请随意向我提问！',
                    'assistant'
                );

                this.showToast('聊天历史已清除', 'success');
            } else {
                this.showToast('清除失败: ' + (data.error || '未知错误'), 'error');
            }
        } catch (error) {
            console.error('Clear history error:', error);
            this.showToast('网络错误，请稍后重试', 'error');
        }
    }

    // 显示状态提示
    showToast(message, type = 'info') {
        this.statusToast.textContent = message;
        this.statusToast.className = `toast ${type}`;
        this.statusToast.classList.add('show');

        setTimeout(() => {
            this.statusToast.classList.remove('show');
        }, 3000);
    }

    // 初始化主题
    initializeTheme() {
        const savedTheme = localStorage.getItem('chat-theme') || 'light';
        this.setTheme(savedTheme);
    }

    // 切换主题
    toggleTheme() {
        const currentTheme = document.body.classList.contains('dark-theme') ? 'dark' : 'light';
        const newTheme = currentTheme === 'light' ? 'dark' : 'light';
        this.setTheme(newTheme);
    }

    // 设置主题
    setTheme(theme) {
        if (theme === 'dark') {
            document.body.classList.add('dark-theme');
            this.themeToggle.querySelector('.icon').textContent = '☀️';
            this.themeToggle.title = '切换到亮色主题';
        } else {
            document.body.classList.remove('dark-theme');
            this.themeToggle.querySelector('.icon').textContent = '🌙';
            this.themeToggle.title = '切换到深色主题';
        }
        localStorage.setItem('chat-theme', theme);
    }
}

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new ChatApp();
});

// 页面可见性变化时的处理
document.addEventListener('visibilitychange', () => {
    if (!document.hidden) {
        // 页面重新可见时，可以在这里添加刷新逻辑
        console.log('Page is now visible');
    }
});

// 处理在线/离线状态
window.addEventListener('online', () => {
    console.log('网络已连接');
});

window.addEventListener('offline', () => {
    console.log('网络已断开');
});
