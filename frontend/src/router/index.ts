import { createRouter, createWebHistory } from 'vue-router'
import Login from "@/views/Login.vue"
import Dashboard from "@/views/Dashboard.vue"
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/login",
      component: Login,
      meta: {
        public: true
      }
    },
    {
      path: "/dashboard",
      component: Dashboard,
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  if (!authStore.isInitialized) {
    await authStore.checkAuth()
  }

  const isPublic = to.matched.some(record => record.meta.public)

  if (isPublic) {
    if (authStore.isAuthenticated && to.path === "/login") {
      return next("/dashboard")
    }
    return next()
  }

  if (!authStore.isAuthenticated) {
    return next({
      path: "/login",
      query: { redirect: to.fullPath }
    })
  }

  next()

})

export default router
