<template>
  <n-space size="small">
    <n-tag
      v-if="worktree.statusAhead > 0 || worktree.statusBehind > 0"
      size="small"
      type="info"
      :bordered="false"
    >
      ↑{{ worktree.statusAhead }} ↓{{ worktree.statusBehind }}
    </n-tag>

    <n-tag v-if="worktree.statusModified > 0" size="small" type="warning" :bordered="false">
      ~{{ worktree.statusModified }}
    </n-tag>

    <n-tag v-if="worktree.statusStaged > 0" size="small" type="success" :bordered="false">
      +{{ worktree.statusStaged }}
    </n-tag>

    <n-tag v-if="worktree.statusUntracked > 0" size="small" :bordered="false">
      ?{{ worktree.statusUntracked }}
    </n-tag>
  </n-space>
</template>

<script setup lang="ts">
import type { Worktree } from '@/types/models';

defineProps<{
  worktree: Worktree;
}>();
</script>
