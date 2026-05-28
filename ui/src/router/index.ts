import { createRouter, createWebHistory } from "vue-router"
import { useAuthStore } from "@/stores/auth"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/auth",
      component: () => import("@/layouts/MinimalLayout.vue"),
      children: [
        {
          path: "login",
          name: "login",
          component: () => import("@/views/LoginView.vue"),
        }
      ]
    },

    {
      path: "/",
      name: "default",
      component: () => import("@/layouts/DefaultLayout.vue"),
      meta: { requiresAuth: true },
      children: [
        {
          path: "",
          name: "dashboard",
          component: () => import("@/views/HomeView.vue"),
        }
      ]
    },

    // not found needs to be the last route!
    {
      path: "/:path(.*)*",
      name: "not-found",
      component: () => import("@/views/NotFoundView.vue")
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  // called after reload to check token
  if (!auth.checked) {
    await auth.check()
  }

  // check if user is logged in, if not redirect to login
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return {
      name: "login",
      query: { redirect: to.fullPath }
    }
  }
})

export default router