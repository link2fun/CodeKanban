<template>
  <div class="project-workspace">
    <n-layout has-sider>
      <n-layout-sider bordered :width="320" :collapsed-width="0" show-trigger="arrow-circle">
        <WorktreeList />
      </n-layout-sider>

      <n-layout-content>
        <div class="workspace-content">
          <KanbanBoard :project-id="currentProjectId" />
        </div>
      </n-layout-content>
    </n-layout>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useProjectStore } from '@/stores/project';
import WorktreeList from '@/components/worktree/WorktreeList.vue';
import KanbanBoard from '@/components/kanban/KanbanBoard.vue';

const route = useRoute();
const projectStore = useProjectStore();

const currentProjectId = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''));

const loadProject = (id: string) => {
  if (!id) {
    return;
  }
  projectStore.fetchProject(id);
};

onMounted(() => {
  if (currentProjectId.value) {
    loadProject(currentProjectId.value);
  }
});

watch(
  () => route.params.id,
  newId => {
    if (typeof newId === 'string') {
      loadProject(newId);
    }
  },
);
</script>

<style scoped>
.project-workspace {
  height: 100vh;
}

.workspace-content {
  padding: 24px;
}
</style>
