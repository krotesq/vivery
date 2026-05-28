import { defineStore } from "pinia"
import { watch } from "vue"
import { useLocalStorage } from "@vueuse/core"

export const useThemeStore = defineStore("theme", () => {

  // define data
  const theme = useLocalStorage("theme", "cupcake")

  // create setter
  function setTheme(t: string) {
    theme.value = t
  }

  // ensure that theme is set on first page load
  watch(theme, (t) => {
    document.documentElement.setAttribute("data-theme", t)
  }, { immediate: true })

  return { theme, setTheme }
})