<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from "vue";

export interface GlassNavItem {
  /** Wert, der über v-model mit der Auswahl verbunden wird. */
  value: string | number;
  /** Sichtbare Beschriftung und zugänglicher Name bei Icon-Navigationen. */
  label: string;
}

const props = withDefaults(
  defineProps<{
    items: GlassNavItem[];
    modelValue?: string | number;
    /** Zugängliche Beschriftung der gesamten Navigation. */
    ariaLabel?: string;
  }>(),
  { ariaLabel: "Navigation" },
);

const emit = defineEmits<{ "update:modelValue": [value: string | number] }>();

const rootEl = ref<HTMLElement | null>(null);
const optionEls = ref<HTMLLabelElement[]>([]);

const activeIndex = computed(() => {
  const i = props.items.findIndex((it) => it.value === props.modelValue);
  return i === -1 ? 0 : i;
});

// Legt die Richtung fest, aus der sich der aktive Indikator verformt.
const originSide = ref<"left" | "right">("left");
const animKey = ref(0);

const indicator = ref({ left: 0, width: 0 });

function measure() {
  const el = optionEls.value[activeIndex.value];
  const root = rootEl.value;
  if (!el || !root) return;
  indicator.value = { left: el.offsetLeft, width: el.offsetWidth };
}

function select(value: string | number) {
  emit("update:modelValue", value);
}

watch(activeIndex, (next, prev) => {
  originSide.value = next >= prev ? "left" : "right";
  animKey.value++;
  nextTick(measure);
});

onMounted(async () => {
  await nextTick();
  measure();
  if (rootEl.value && "ResizeObserver" in window) {
    new ResizeObserver(measure).observe(rootEl.value);
  }
});
</script>

<template>
  <div
    ref="rootEl"
    class="liquid-glass-nav"
    role="tablist"
    :aria-label="ariaLabel"
  >
    <!-- Der Indikator wird anhand der gemessenen Breite des aktiven Eintrags verschoben. -->
    <span
      :key="animKey"
      class="liquid-glass-nav-indicator"
      :style="{
        left: indicator.left + 'px',
        width: indicator.width + 'px',
        transformOrigin: originSide,
      }"
      aria-hidden="true"
    />

    <label
      v-for="item in items"
      :key="item.value"
      ref="optionEls"
      class="liquid-glass-nav-option"
      :class="item.value === modelValue && 'liquid-glass-nav-option-active'"
    >
      <input
        class="sr-only"
        type="radio"
        :name="'glass-nav-' + ariaLabel"
        :value="item.value"
        :checked="item.value === modelValue"
        role="tab"
        :aria-selected="item.value === modelValue"
        @change="select(item.value)"
      />
      <span class="liquid-glass-nav-content">
        <slot name="icon" :item="item" :active="item.value === modelValue">{{
          item.label
        }}</slot>
      </span>
    </label>

    <!-- Der Filter erzeugt die Brechung des Liquid-Glass-Effekts. -->
    <div class="liquid-glass-nav-filter" aria-hidden="true">
      <svg xmlns="glass-filter">
        <filter id="liquid-glass-nav" primitiveUnits="objectBoundingBox">
          <feImage
            result="map"
            width="100%"
            height="100%"
            x="0"
            y="0"
            href="/liquid-map.webp"
            />
          <feGaussianBlur
            in="SourceGraphic"
            stdDeviation="0.04"
            result="blur"
          />
          <feDisplacementMap
            in="blur"
            in2="map"
            scale="0.5"
            xChannelSelector="R"
            yChannelSelector="G"
          />
        </filter>
      </svg>
    </div>
  </div>
</template>
