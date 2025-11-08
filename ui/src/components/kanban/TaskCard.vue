<template>
  <n-card class="task-card" size="small" @click="emit('click')">
    <div class="task-card__header">
      <n-ellipsis :line-clamp="2">
        {{ task.title }}
      </n-ellipsis>
      <n-tag v-if="priorityLabel" size="tiny" :type="priorityType" :bordered="false">
        {{ priorityLabel }}
      </n-tag>
    </div>

    <div class="task-card__meta">
      <n-space size="small" wrap>
        <n-tag v-if="task.worktree?.branchName" size="tiny" :bordered="false">
          {{ task.worktree.branchName }}
        </n-tag>
        <n-tag v-if="task.dueDate" size="tiny" :type="isOverdue ? 'error' : 'default'" :bordered="false">
          截止 {{ formatDate(task.dueDate) }}
        </n-tag>
      </n-space>
    </div>

    <n-space v-if="task.tags?.length" size="small" wrap>
      <n-tag v-for="tag in task.tags" :key="tag" size="tiny" :bordered="false">
        {{ tag }}
      </n-tag>
    </n-space>
  </n-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import dayjs from 'dayjs';
import type { Task } from '@/types/models';

const props = defineProps<{
  task: Task;
}>();

const emit = defineEmits<{
  click: [];
}>();

const priorityMap: Record<number, { label: string; type: 'default' | 'info' | 'warning' | 'error' }> = {
  0: { label: '', type: 'default' },
  1: { label: '低', type: 'info' },
  2: { label: '中', type: 'warning' },
  3: { label: '高', type: 'error' },
};

const priorityType = computed(() => priorityMap[props.task.priority]?.type ?? 'default');
const priorityLabel = computed(() => priorityMap[props.task.priority]?.label ?? '');

const isOverdue = computed(() => {
  if (!props.task.dueDate) {
    return false;
  }
  return dayjs(props.task.dueDate).isBefore(dayjs(), 'day');
});

const formatDate = (value: string) => dayjs(value).format('MM-DD');
</script>

<style scoped>
.task-card {
  cursor: pointer;
  transition: all 0.2s ease;
}

.task-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.12);
}

.task-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

.task-card__meta {
  margin-bottom: 8px;
}
</style>
