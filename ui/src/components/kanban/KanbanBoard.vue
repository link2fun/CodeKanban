<template>
  <div class="kanban-board">
    <div class="board-header">
      <n-space justify="space-between" align="center">
        <div>
          <h2>任务看板</h2>
          <n-text depth="3">拖拽卡片以重新排序或切换状态</n-text>
        </div>
        <n-button type="primary" :disabled="!projectId" @click="showCreateDialog = true">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
          新建任务
        </n-button>
      </n-space>
    </div>

    <div class="board-body">
      <n-spin :show="boardLoading">
        <n-empty v-if="!projectId" description="请选择一个项目查看任务" />
        <div v-else class="board-columns">
          <KanbanColumn
            v-for="column in columns"
            :key="column.key"
            :title="column.title"
            :status="column.key"
            :tasks="taskStore.tasksByStatus[column.key] ?? []"
            @task-moved="handleTaskMoved"
            @task-clicked="handleTaskClicked"
          />
        </div>
      </n-spin>
    </div>

    <TaskCreateDialog
      v-if="projectId"
      v-model:show="showCreateDialog"
      :project-id="projectId"
      @created="handleTaskCreated"
    />

    <TaskDetailDrawer
      v-model:show="showDetailDrawer"
      :project-id="projectId"
      :task-id="taskStore.selectedTaskId"
      @closed="taskStore.selectTask(null)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useMessage } from 'naive-ui';
import { AddOutline } from '@vicons/ionicons5';
import KanbanColumn from './KanbanColumn.vue';
import TaskCreateDialog from './TaskCreateDialog.vue';
import TaskDetailDrawer from './TaskDetailDrawer.vue';
import { useTaskStore } from '@/stores/task';
import { useTaskActions } from '@/composables/useTaskActions';
import type { Task } from '@/types/models';

const props = defineProps<{
  projectId?: string;
}>();

const taskStore = useTaskStore();
const message = useMessage();
const { listTasks, moveTask } = useTaskActions();

const showCreateDialog = ref(false);
const showDetailDrawer = ref(false);
const boardLoading = ref(false);

const columns = [
  { key: 'todo', title: '待办' },
  { key: 'in_progress', title: '进行中' },
  { key: 'done', title: '已完成' },
] as const;

const currentProjectId = computed(() => props.projectId ?? '');

watch(
  () => currentProjectId.value,
  id => {
    if (id) {
      fetchTasks(id);
    } else {
      taskStore.setTasks([]);
    }
  },
  { immediate: true },
);

async function fetchTasks(projectId: string) {
  boardLoading.value = true;
  try {
    const response = await listTasks.send(projectId);
    const items = response?.body?.items ?? [];
    taskStore.setTasks(items as Task[]);
  } catch (error: any) {
    message.error(error?.message ?? '加载任务失败');
  } finally {
    boardLoading.value = false;
  }
}

async function handleTaskMoved(event: { taskId: string; newStatus: Task['status']; newIndex: number; orderedTasks: Task[] }) {
  const { taskId, newStatus, newIndex, orderedTasks } = event;
  const siblings = orderedTasks;
  let orderIndex = 1000;

  if (siblings.length <= 1) {
    orderIndex = 1000;
  } else if (newIndex <= 0) {
    const next = siblings[1];
    orderIndex = next ? next.orderIndex / 2 : 500;
  } else if (newIndex >= siblings.length - 1) {
    const prev = siblings[newIndex - 1] ?? siblings[siblings.length - 2];
    orderIndex = prev.orderIndex + 1000;
  } else {
    const prev = siblings[newIndex - 1];
    const next = siblings[newIndex + 1];
    orderIndex = prev && next ? (prev.orderIndex + next.orderIndex) / 2 : prev?.orderIndex ?? 1000;
  }

  try {
    const response = await moveTask.send(taskId, { status: newStatus, orderIndex });
    const updated = response?.body?.item as Task | undefined;
    if (updated) {
      taskStore.upsertTask(updated);
    }
  } catch (error: any) {
    message.error(error?.message ?? '移动任务失败');
    fetchTasks(currentProjectId.value);
  }
}

function handleTaskClicked(task: Task) {
  taskStore.selectTask(task.id);
  showDetailDrawer.value = true;
}

function handleTaskCreated(task: Task) {
  taskStore.upsertTask(task);
  message.success('任务创建成功');
}
</script>

<style scoped>
.kanban-board {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.board-header {
  padding: 16px 24px;
  border-bottom: 1px solid var(--n-border-color);
}

.board-header h2 {
  margin: 0 0 4px;
}

.board-body {
  flex: 1;
  padding: 16px;
  overflow: hidden;
}

.board-columns {
  display: grid;
  grid-template-columns: repeat(3, minmax(280px, 1fr));
  gap: 16px;
  height: 100%;
}

@media (max-width: 1200px) {
  .board-columns {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  }
}
</style>
