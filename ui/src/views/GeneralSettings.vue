<template>
  <div class="general-settings-page">
    <n-page-header @back="handleBack">
      <template #title>
        <n-space align="center">
          <n-icon size="24">
            <SettingsOutline />
          </n-icon>
          <span>总设置</span>
        </n-space>
      </template>
      <template #extra>
        <n-space>
          <n-button tertiary @click="handleResetTheme">
            <template #icon>
              <n-icon>
                <RefreshOutline />
              </n-icon>
            </template>
            恢复默认主题
          </n-button>
        </n-space>
      </template>
    </n-page-header>

    <n-space vertical size="large">
      <n-card title="主题设置" size="huge">
        <n-form label-placement="left" label-width="120">
          <n-form-item label="主色调">
            <n-color-picker v-model:value="primaryColor" :modes="['hex']" :actions="['confirm']" />
          </n-form-item>
          <n-form-item label="界面背景色">
            <n-color-picker v-model:value="bodyColor" :modes="['hex']" :actions="['confirm']" />
          </n-form-item>
          <n-form-item label="卡片背景色">
            <n-color-picker v-model:value="surfaceColor" :modes="['hex']" :actions="['confirm']" />
          </n-form-item>
        </n-form>
      </n-card>

      <n-card title="实时预览" size="huge">
        <div class="preview-panel" :style="{ backgroundColor: surfaceColor }">
          <div class="preview-banner" :style="{ backgroundColor: primaryColor }">
            <n-space align="center" size="small">
              <n-icon size="24" color="#fff">
                <ColorPaletteOutline />
              </n-icon>
              <span>看看新主题的样子</span>
            </n-space>
          </div>
          <div class="preview-content">
            <n-space vertical size="medium">
              <n-button type="primary">主要按钮</n-button>
              <n-tag type="primary" :bordered="false">主色标签</n-tag>
              <n-alert type="info" title="提示">
                主题色会实时保存，所有页面会立即生效。
              </n-alert>
            </n-space>
          </div>
        </div>
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useTitle } from '@vueuse/core';
import { ColorPaletteOutline, SettingsOutline, RefreshOutline } from '@vicons/ionicons5';
import { useSettingsStore } from '@/stores/settings';
import { APP_NAME } from '@/constants/app';

useTitle(`总设置 - ${APP_NAME}`);

const router = useRouter();
const settingsStore = useSettingsStore();
const { theme } = storeToRefs(settingsStore);

const primaryColor = computed({
  get: () => theme.value.primaryColor,
  set: value => settingsStore.updateTheme({ primaryColor: value || '#18a058' }),
});

const bodyColor = computed({
  get: () => theme.value.bodyColor,
  set: value => settingsStore.updateTheme({ bodyColor: value || '#f7f8fa' }),
});

const surfaceColor = computed({
  get: () => theme.value.surfaceColor,
  set: value => settingsStore.updateTheme({ surfaceColor: value || '#ffffff' }),
});

function handleBack() {
  router.back();
}

function handleResetTheme() {
  settingsStore.resetTheme();
}
</script>

<style scoped>
.general-settings-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px;
}

.preview-panel {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--n-border-color);
}

.preview-banner {
  color: #fff;
  padding: 16px;
  font-size: 16px;
  font-weight: 600;
}

.preview-content {
  padding: 24px;
  background-color: rgba(255, 255, 255, 0.8);
}
</style>
