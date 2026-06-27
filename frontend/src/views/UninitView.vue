<template>
  <div class="uninit">
    <TitleBar @minimize="onMinimize" @maximize="onMaximize" @close="onClose" />
    <div class="uninit-body">
      <div class="card">
        <div class="logo">E</div>
        <h1>Ease 未初始化</h1>
        <p class="hint">请在终端中运行以下命令设置账户密码：</p>
        <pre class="cmd">ease-ui init</pre>
        <p class="hint">设置完成后重新启动本应用。</p>
        <button class="quit" @click="onClose">退出</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import TitleBar from '../components/TitleBar.vue'
import { WindowMinimise, WindowToggleMaximise, WindowQuit } from '../composables/useWails'
function onMinimize()  { WindowMinimise() }
function onMaximize()  { WindowToggleMaximise() }
function onClose()     { WindowQuit() }
</script>

<style scoped>
.uninit { display: flex; flex-direction: column; height: 100vh; }
.uninit-body {
  flex: 1; display: flex; align-items: center; justify-content: center;
  background: var(--bg-primary);
}
.card {
  text-align: center; padding: 32px;
  border: 1px solid var(--border); border-radius: var(--radius-lg);
  background: var(--bg-panel); max-width: 360px;
}
.logo {
  width: 48px; height: 48px; border-radius: 12px;
  background: linear-gradient(135deg, var(--accent), var(--accent-deep));
  display: flex; align-items: center; justify-content: center;
  color: white; font-size: 22px; font-weight: 700;
  margin: 0 auto 16px;
  box-shadow: var(--shadow-accent);
}
h1 { font-size: 16px; margin-bottom: 12px; }
.hint { color: var(--text-secondary); font-size: 12px; margin: 8px 0; }
.cmd {
  font-family: var(--font-mono); background: var(--bg-input);
  padding: 8px 12px; border-radius: var(--radius-md);
  display: inline-block; margin: 8px 0;
  border: 1px solid var(--border);
}
.quit {
  margin-top: 16px; padding: 8px 20px;
  background: var(--bg-input); border: 1px solid var(--border);
  border-radius: var(--radius-md); color: var(--text-primary);
}
.quit:hover { background: var(--border); }
</style>
