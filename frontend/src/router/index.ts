import { createRouter, createWebHashHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import SettingsView from '../views/SettingsView.vue'
import UninitView from '../views/UninitView.vue'
import { useAuthStore } from '../stores/auth'
import { IsInitialized } from '../composables/useWails'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/login', component: LoginView },
    { path: '/uninit', component: UninitView },
    { path: '/home', component: HomeView, meta: { requiresAuth: true } },
    { path: '/settings', component: SettingsView },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (!auth.initialized) {
    const init = await IsInitialized()
    auth.initialized = init
    if (!init) return '/uninit'
  }
  if (to.meta.requiresAuth && !auth.loggedIn) {
    return '/login'
  }
  return true
})

export default router
