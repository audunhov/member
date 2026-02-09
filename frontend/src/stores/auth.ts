import { defineStore } from "pinia";

type User = {
  id: string,
}

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: null as User | null,
    isInitialized: false,
  }),
  getters: {
    isAuthenticated: (state) => !!state.user,
  },
  actions: {
    async checkAuth() {
      try {
        const res = await fetch('/api/v1/auth/me')
        if (res.ok) {
          this.user = await res.json()
        } else {
          this.user = null
        }
      } catch (err) {
        this.user = null
      } finally {
        this.isInitialized = true
      }

    }
  }
})
