# 自动隐藏页头页尾功能

## 功能概述

Web 终端在 WebSocket 连接成功后自动隐藏页头和页尾，提供全屏终端体验。断开连接后自动显示页头和页尾，方便用户查看连接状态和操作提示。

## 用户体验改进

### 之前的问题

1. 页头和页尾占据屏幕空间
2. 终端显示区域受限
3. 不够专业的远程终端体验
4. 用户需要手动全屏浏览器才能获得较好体验

### 改进后的效果

1. **连接前**：显示完整的页头和页尾
   - 页头：Logo + 连接状态 + 会话信息 + 断开按钮
   - 页尾：快捷键提示 + 版权信息

2. **连接后**：自动隐藏页头和页尾
   - 终端占满整个屏幕
   - 全屏沉浸式体验
   - 类似本地终端的感觉

3. **断开后**：自动恢复页头和页尾
   - 显示连接状态
   - 提供重连提示

## 实现细节

### CSS 样式

**文件**：`web/terminal.html`

```css
/* 隐藏页头和页尾的样式 */
.header.hidden,
.footer.hidden {
    opacity: 0;
    pointer-events: none;
    transform: translateY(-20px);
}

.footer.hidden {
    transform: translateY(20px);
}

/* 连接后完全隐藏 */
body.has-connected .header,
body.has-connected .footer {
    display: none;
}

body.has-connected .terminal-container {
    padding: 0;
}
```

### JavaScript 逻辑

#### WebSocket 连接成功时隐藏

```javascript
ws.onopen = () => {
    updateConnectionStatus('connected');

    // 自动隐藏页头和页尾
    document.body.classList.add('has-connected');
    document.querySelector('.header').classList.add('hidden');
    document.querySelector('.footer').classList.add('hidden');
};
```

#### WebSocket 断开时显示

```javascript
ws.onclose = (event) => {
    console.log('WebSocket关闭:', event);
    updateConnectionStatus('disconnected');

    // 重新显示页头和页尾
    document.body.classList.remove('has-connected');
    document.querySelector('.header').classList.remove('hidden');
    document.querySelector('.footer').classList.remove('hidden');
};
```

#### 移除不必要的终端消息

同时移除了所有连接状态的终端输出，保持终端干净：

**移除的消息**：
- `✅ WebSocket连接已建立`
- `❌ 连接已断开: code=${event.code}`
- `❌ WebSocket错误`
- `❌ 连接失败: ${error.message}`
- `❌ 解析消息失败: ${e.message}`

**保留**：浏览器控制台日志（用于调试）

## 视觉效果对比

### 之前

```
┌─────────────────────────────────────────┐
│ CliGool - 远程终端    已连接   [断开] │  ← 页头始终显示
├─────────────────────────────────────────┤
│                                         │
│  $ ls -la                              │
│  total 16                              │
│  ...                                   │
│                                         │
├─────────────────────────────────────────┤
│ Ctrl+C 复制  Ctrl+V 粘贴  CliGool...  │  ← 页尾始终显示
└─────────────────────────────────────────┘
```

### 现在

**连接前**：显示页头页尾
```
┌─────────────────────────────────────────┐
│ CliGool - 远程终端    连接中...         │  ← 显示状态
├─────────────────────────────────────────┤
│                                         │
│  (等待连接...)                          │
│                                         │
├─────────────────────────────────────────┤
│ Ctrl+C 复制  Ctrl+V 粘贴  CliGool...  │
└─────────────────────────────────────────┘
```

**连接后**：自动隐藏，全屏终端
```
┌─────────────────────────────────────────┐
│ $ ls -la                              │  ← 全屏显示
│ total 16                              │
│ drwxr-xr-x  4 user  staff   128 Mar  8 │
│ ...                                   │
│                                         │
│                                         │
└─────────────────────────────────────────┘
```

## 测试方法

### 快速测试

使用提供的测试脚本：

```bash
./test-auto-hide-ui.sh
```

### 手动测试

1. **启动 CLI 客户端**：
   ```bash
   ./bin/cligool-darwin-arm64 -server http://localhost:8081 -session test-auto-hide
   ```

2. **打开浏览器**：
   ```
   http://localhost:8081/session/test-auto-hide
   ```

3. **观察效果**：
   - 页面加载时显示页头和页尾
   - 连接成功后页头页尾消失
   - 刷新浏览器或断开后页头页尾重新显示

### 开发者工具检查

打开浏览器开发者工具（F12）：

**Elements 标签**：
- 检查 `<body>` 元素是否有 `has-connected` class
- 检查 `.header` 和 `.footer` 元素是否有 `hidden` class

**Console 标签**：
- 查看连接日志：`WebSocket关闭:` 等

## 兼容性

- ✅ 所有现代浏览器（Chrome、Firefox、Safari、Edge）
- ✅ 移动浏览器（iOS Safari、Chrome Mobile）
- ✅ 平板浏览器（iPad、Android Tablet）

## 相关功能

这个功能与之前实现的其他用户体验改进配合使用：

1. **消息缓存**：防止消息丢失
2. **平滑初始化**：避免视觉闪烁
3. **无CLI提示**：友好的错误提示
4. **自动隐藏UI**：全屏沉浸式体验（本文档）

## 未来改进

可能的进一步优化：

1. **可选显示**：添加设置让用户选择是否隐藏UI
2. **鼠标悬停**：鼠标移动到顶部时临时显示页头
3. **快捷键**：使用快捷键切换UI显示/隐藏
4. **主题切换**：支持亮色/暗色主题

## 相关文件

- `web/terminal.html` - Web 终端实现
- `test-auto-hide-ui.sh` - 测试脚本
- `docs/AUTO_HIDE_UI_FEATURE.md` - 本文档

## 总结

自动隐藏页头页尾功能提供了更专业的远程终端体验：

- ✅ 全屏沉浸式体验
- ✅ 更大的终端显示区域
- ✅ 类似本地终端的感觉
- ✅ 自动切换，无需用户操作
- ✅ 保持友好的断开提示

实施简单，效果显著，是非常有价值的用户体验改进！
