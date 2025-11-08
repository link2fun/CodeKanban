<script setup lang="ts">
import { computed } from 'vue';
import { RouterView } from 'vue-router';
import { zhCN, dateZhCN, type GlobalThemeOverrides } from 'naive-ui';

const themeOverrides = computed<GlobalThemeOverrides>(() => ({
  common: {},
  Scrollbar: {
    width: '8px',
    height: '8px',
  },
}));
</script>

<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme-overrides="themeOverrides">
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <n-modal-provider>
              <RouterView />
            </n-modal-provider>
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>
