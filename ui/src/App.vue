<script setup lang="ts">
import { computed, useCssVars, watchEffect } from 'vue';
import { RouterView } from 'vue-router';
import { storeToRefs } from 'pinia';
import { zhCN, dateZhCN, type GlobalThemeOverrides } from 'naive-ui';
import AppInitializer from '@/components/common/AppInitializer.vue';
import NotePad from '@/components/notepad/NotePad.vue';
import { useSettingsStore } from '@/stores/settings';
import { darkenColor, lightenColor } from '@/utils/color';

const settingsStore = useSettingsStore();
const { theme } = storeToRefs(settingsStore);

const themeOverrides = computed<GlobalThemeOverrides>(() => {
  const { primaryColor, bodyColor, surfaceColor } = theme.value;
  const primaryHover = lightenColor(primaryColor, 0.08);
  const primaryPressed = darkenColor(primaryColor, 0.12);

  return {
    common: {
      bodyColor,
      cardColor: surfaceColor,
      modalColor: surfaceColor,
      popoverColor: surfaceColor,
      primaryColor,
      primaryColorHover: primaryHover,
      primaryColorPressed: primaryPressed,
      primaryColorSuppl: primaryHover,
    },
    Layout: {
      color: surfaceColor,
      siderColor: surfaceColor,
      headerColor: surfaceColor,
      footerColor: surfaceColor,
    },
    Scrollbar: {
      width: '8px',
      height: '8px',
    },
  };
});

useCssVars(() => ({
  'app-body-color': theme.value.bodyColor,
  'app-surface-color': theme.value.surfaceColor,
}));

watchEffect(() => {
  if (typeof document === 'undefined') {
    return;
  }
  document.body.style.backgroundColor = theme.value.bodyColor;
});
</script>

<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme-overrides="themeOverrides">
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <n-modal-provider>
              <AppInitializer />
              <RouterView />
              <NotePad />
            </n-modal-provider>
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>
