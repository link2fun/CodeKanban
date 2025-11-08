import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import type { Task, TaskComment } from '@/types/models';

type StatusBuckets = Record<string, Task[]>;

const defaultStatuses = ['todo', 'in_progress', 'done', 'archived'] as const;

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([]);
  const selectedTaskId = ref<string | null>(null);
  const commentsMap = ref<Record<string, TaskComment[]>>({});

  const selectedTask = computed<Task | null>(() => {
    return tasks.value.find(task => task.id === selectedTaskId.value) ?? null;
  });

  const tasksByStatus = computed<StatusBuckets>(() => {
    const buckets: StatusBuckets = {};
    for (const status of defaultStatuses) {
      buckets[status] = [];
    }
    for (const task of tasks.value) {
      if (!buckets[task.status]) {
        buckets[task.status] = [];
      }
      buckets[task.status].push(task);
    }
    for (const bucket of Object.values(buckets)) {
      bucket.sort((a, b) => a.orderIndex - b.orderIndex);
    }
    return buckets;
  });

  function setTasks(list: Task[]) {
    tasks.value = [...list];
  }

  function upsertTask(task: Task) {
    const index = tasks.value.findIndex(item => item.id === task.id);
    if (index === -1) {
      tasks.value.push(task);
    } else {
      tasks.value.splice(index, 1, task);
    }
  }

  function removeTask(id: string) {
    tasks.value = tasks.value.filter(task => task.id !== id);
    delete commentsMap.value[id];
    if (selectedTaskId.value === id) {
      selectedTaskId.value = null;
    }
  }

  function selectTask(id: string | null) {
    selectedTaskId.value = id;
  }

  function setComments(taskId: string, list: TaskComment[]) {
    commentsMap.value = {
      ...commentsMap.value,
      [taskId]: list,
    };
  }

  function appendComment(taskId: string, comment: TaskComment) {
    const existing = commentsMap.value[taskId] ?? [];
    commentsMap.value = {
      ...commentsMap.value,
      [taskId]: [...existing, comment],
    };
  }

  function removeComment(taskId: string, commentId: string) {
    const existing = commentsMap.value[taskId];
    if (!existing) {
      return;
    }
    commentsMap.value = {
      ...commentsMap.value,
      [taskId]: existing.filter(comment => comment.id !== commentId),
    };
  }

  return {
    tasks,
    selectedTaskId,
    selectedTask,
    tasksByStatus,
    commentsMap,
    setTasks,
    upsertTask,
    removeTask,
    selectTask,
    setComments,
    appendComment,
    removeComment,
  };
});
