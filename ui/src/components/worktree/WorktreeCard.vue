<template>
  <n-card class="worktree-card" :class="{ 'is-main': worktree.isMain }" size="small">
    <template #header>
      <n-space justify="space-between" align="center">
        <n-ellipsis style="max-width: 180px">
          {{ worktree.branchName }}
        </n-ellipsis>
        <n-dropdown :options="actions" @select="handleAction">
          <n-button text size="small">
            <n-icon><EllipsisHorizontalOutline /></n-icon>
          </n-button>
        </n-dropdown>
      </n-space>
    </template>

    <n-space vertical size="small">
      <GitStatusBadge :worktree="worktree" />

      <n-text depth="3" class="meta-text">
        {{ worktree.headCommit || '无提交信息' }}
      </n-text>

      <n-text depth="3" class="meta-text">
        {{ formatTime(worktree.statusUpdatedAt) }}
      </n-text>
    </n-space>
  </n-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { EllipsisHorizontalOutline } from '@vicons/ionicons5';
import GitStatusBadge from '@/components/common/GitStatusBadge.vue';
import type { Worktree } from '@/types/models';

dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

const props = defineProps<{
  worktree: Worktree;
}>();

const emit = defineEmits<{
  refresh: [id: string];
  delete: [worktree: Worktree];
  'open-explorer': [path: string];
  'open-terminal': [path: string];
}>();

const actions = computed(() => [
  { label: '打开文件管理器', key: 'explorer' },
  { label: '打开终端', key: 'terminal' },
  { label: '刷新状态', key: 'refresh' },
  { label: '删除', key: 'delete', disabled: props.worktree.isMain },
]);

function handleAction(key: string | number) {
  switch (key) {
    case 'explorer':
      emit('open-explorer', props.worktree.path);
      break;
    case 'terminal':
      emit('open-terminal', props.worktree.path);
      break;
    case 'refresh':
      emit('refresh', props.worktree.id);
      break;
    case 'delete':
      emit('delete', props.worktree);
      break;
    default:
      break;
  }
}

function formatTime(time: string | null) {
  if (!time) {
    return '未更新';
  }
  return dayjs(time).fromNow();
}
</script>

<style scoped>
.worktree-card {
  margin-bottom: 8px;
  cursor: pointer;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
}

.worktree-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.worktree-card.is-main {
  border-color: var(--n-color-primary);
}

.meta-text {
  font-size: 12px;
}
</style>
