import { defineStore } from "pinia"
import { watch } from "vue"
import { useLocalStorage } from "@vueuse/core"

const DEFAULT_THEME = "dark"
const THEMES = new Set(["light", "dark"])

function applyTheme(theme: string) {
  document.documentElement.setAttribute("data-theme", theme)
}

function normalizeTheme(theme: string | null) {
  return theme && THEMES.has(theme) ? theme : DEFAULT_THEME
}

export function initializeTheme() {
  const theme = normalizeTheme(localStorage.getItem("theme"))
  localStorage.setItem("theme", theme)
  applyTheme(theme)
}

export const useThemeStore = defineStore("theme", () => {

  // define data
  const theme = useLocalStorage("theme", DEFAULT_THEME)

  // create setter
  function setTheme(t: string) {
    theme.value = normalizeTheme(t)
  }

  // ensure that theme is set on first page load
  watch(theme, (t) => {
    applyTheme(normalizeTheme(t))
  }, { immediate: true })

  return { theme, setTheme }
})
