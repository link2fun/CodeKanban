<template>
  <n-modal
    v-model:show="visible"
    preset="dialog"
    :title="t('worktree.create')"
    :positive-text="t('common.create')"
    :negative-text="t('common.cancel')"
    :loading="loading"
    @positive-click="handleCreate"
  >
    <n-form ref="formRef" :model="formData" :rules="rules" label-placement="top">
      <n-form-item :label="t('branch.branchName')" path="branchName">
        <n-input v-model:value="formData.branchName" :placeholder="t('branch.branchNamePlaceholder')" />
      </n-form-item>

      <n-form-item :label="t('branch.baseBranch')" path="baseBranch">
        <n-input v-model:value="formData.baseBranch" :placeholder="t('branch.baseBranchPlaceholder')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useMessage, type FormInst, type FormRules } from 'naive-ui';
import { useProjectStore } from '@/stores/project';
import type { Worktree } from '@/types/models';
import { useLocale } from '@/composables/useLocale';

const { t } = useLocale();

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
  success: [worktree: Worktree];
}>();

const projectStore = useProjectStore();
const message = useMessage();

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value),
});

const formRef = ref<FormInst | null>(null);
const loading = ref(false);
const formData = ref({
  branchName: '',
  baseBranch: '',
  createBranch: true,
});

const rules: FormRules = {
  branchName: [{ required: true, message: t('validation.branchNameRequired'), trigger: ['blur', 'input'] }],
};

watch(visible, newVal => {
  if (newVal) {
    formData.value.baseBranch =
      projectStore.currentProject?.defaultBranch ?? formData.value.baseBranch ?? 'main';
  } else {
    formData.value = { branchName: '', baseBranch: '', createBranch: true };
  }
});

async function handleCreate() {
  if (!projectStore.currentProject) {
    message.error(t('project.selectProjectFirst'));
    return false;
  }

  try {
    await formRef.value?.validate();
    loading.value = true;
    const worktree = await projectStore.createWorktree(projectStore.currentProject.id, formData.value);
    // 先 emit success 事件，确保父组件能接收到
    emit('success', worktree);
    // 返回 true 让 Naive UI 自动关闭对话框
    return true;
  } catch (error: any) {
    message.error(error?.message ?? t('worktree.createFailed'));
    return false;
  } finally {
    loading.value = false;
  }
}
</script>
