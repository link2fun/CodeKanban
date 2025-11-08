import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import { projectApi, systemApi, worktreeApi } from '@/api/project';
import type { Project, Worktree } from '@/types/models';

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([]);
  const currentProject = ref<Project | null>(null);
  const worktrees = ref<Worktree[]>([]);
  const loading = ref(false);

  const hasProjects = computed(() => projects.value.length > 0);

  async function fetchProjects() {
    loading.value = true;
    try {
      const result = await projectApi.list();
      projects.value = result.items;
    } finally {
      loading.value = false;
    }
  }

  async function fetchProject(id: string) {
    loading.value = true;
    try {
      currentProject.value = await projectApi.get(id);
      await fetchWorktrees(id);
    } finally {
      loading.value = false;
    }
  }

  async function createProject(payload: { name: string; path: string; description?: string }) {
    const project = await projectApi.create(payload);
    projects.value.push(project);
    return project;
  }

  async function deleteProject(id: string) {
    await projectApi.delete(id);
    projects.value = projects.value.filter(project => project.id !== id);
    if (currentProject.value?.id === id) {
      currentProject.value = null;
      worktrees.value = [];
    }
  }

  async function fetchWorktrees(projectId: string) {
    worktrees.value = await worktreeApi.list(projectId);
  }

  async function createWorktree(
    projectId: string,
    payload: { branchName: string; baseBranch?: string; createBranch?: boolean },
  ) {
    const worktree = await worktreeApi.create(projectId, payload);
    worktrees.value.push(worktree);
    return worktree;
  }

  async function deleteWorktree(id: string, force = false, deleteBranch = false) {
    await worktreeApi.delete(id, force, deleteBranch);
    worktrees.value = worktrees.value.filter(worktree => worktree.id !== id);
  }

  async function refreshWorktreeStatus(id: string) {
    const updated = await worktreeApi.refreshStatus(id);
    const index = worktrees.value.findIndex(worktree => worktree.id === id);
    if (index !== -1) {
      worktrees.value.splice(index, 1, updated);
    }
  }

  async function refreshAllWorktrees(projectId: string) {
    await worktreeApi.refreshAll(projectId);
    await fetchWorktrees(projectId);
  }

  async function syncWorktrees(projectId: string) {
    await worktreeApi.sync(projectId);
    await fetchWorktrees(projectId);
  }

  async function openInExplorer(path: string) {
    await systemApi.openExplorer(path);
  }

  async function openInTerminal(path: string) {
    await systemApi.openTerminal(path);
  }

  return {
    projects,
    currentProject,
    worktrees,
    loading,
    hasProjects,
    fetchProjects,
    fetchProject,
    createProject,
    deleteProject,
    fetchWorktrees,
    createWorktree,
    deleteWorktree,
    refreshWorktreeStatus,
    refreshAllWorktrees,
    syncWorktrees,
    openInExplorer,
    openInTerminal,
  };
});
