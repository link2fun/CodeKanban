<template>
  <n-card class="task-card" size="small" @click="emit('click')">
    <div class="task-card__header">
      <div class="task-card__title">
        <n-ellipsis :line-clamp="2">
          {{ task.title }}
        </n-ellipsis>
        <n-tag v-if="priorityLabel" size="tiny" :type="priorityType" :bordered="false">
          {{ priorityLabel }}
        </n-tag>
      </div>
      <div class="task-card__actions">
        <n-tooltip trigger="hover" placement="bottom">
          <template #trigger>
            <n-button quaternary circle size="tiny" @click.stop="handleCopy">
              <n-icon size="14">
                <CopyOutline />
              </n-icon>
            </n-button>
          </template>
          复制任务名
        </n-tooltip>
        <n-tooltip trigger="hover" placement="bottom">
          <template #trigger>
            <n-button quaternary circle size="tiny" @click.stop="handleEdit">
              <n-icon size="14">
                <CreateOutline />
              </n-icon>
            </n-button>
          </template>
          快捷编辑
        </n-tooltip>
        <n-tooltip trigger="hover" placement="bottom">
          <template #trigger>
            <n-button quaternary circle size="tiny" @click.stop="handleDelete">
              <n-icon size="14">
                <TrashOutline />
              </n-icon>
            </n-button>
          </template>
          快捷删除
        </n-tooltip>
      </div>
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
import { CopyOutline, CreateOutline, TrashOutline } from '@vicons/ionicons5';
import type { Task } from '@/types/models';

const props = defineProps<{
  task: Task;
}>();

const emit = defineEmits<{
  click: [];
  edit: [];
  delete: [];
  copy: [];
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

const handleEdit = () => emit('edit');
const handleDelete = () => emit('delete');
const handleCopy = () => emit('copy');
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

.task-card__title {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.task-card__actions {
  display: flex;
  align-items: center;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s ease;
  pointer-events: none;
}

.task-card:hover .task-card__actions {
  opacity: 1;
  pointer-events: auto;
}

.task-card__meta {
  margin-bottom: 8px;
}
</style>
