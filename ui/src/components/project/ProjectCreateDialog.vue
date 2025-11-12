<template>
  <n-modal
    v-model:show="visible"
    preset="dialog"
    title="创建项目"
    positive-text="创建"
    negative-text="取消"
    :loading="loading"
    @positive-click="handleCreate"
  >
    <n-form ref="formRef" :model="formData" :rules="rules" label-placement="top">
      <n-form-item label="项目名称" path="name">
        <n-input v-model:value="formData.name" placeholder="输入项目名称" />
      </n-form-item>
      <n-form-item label="项目目录" path="path">
        <n-input
          v-model:value="formData.path"
          placeholder="选择或输入本地项目目录，例如 C:\\Projects\\demo"
        />
        <template #feedback>
          <n-text depth="3">
            可以直接选择任意本地文件夹；若目录中没有 .git，将无法使用分支与 Worktree 功能。
          </n-text>
        </template>
      </n-form-item>
      <n-form-item label="项目描述" path="description">
        <n-input
          v-model:value="formData.description"
          type="textarea"
          :rows="3"
          placeholder="输入项目描述（可选）"
        />
      </n-form-item>
      <n-form-item label="隐藏路径" path="hidePath">
        <n-space align="center">
          <n-switch v-model:value="formData.hidePath" />
          <n-text depth="3">开启后，项目列表与侧边栏中将不再展示绝对路径。</n-text>
        </n-space>
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useMessage, type FormInst, type FormRules } from 'naive-ui';
import { useProjectStore } from '@/stores/project';
import type { Project } from '@/types/models';

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
  success: [project: Project];
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
  name: '',
  path: '',
  description: '',
  hidePath: false,
});

const rules: FormRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: ['blur', 'input'] }],
  path: [{ required: true, message: '请输入项目目录', trigger: ['blur', 'input'] }],
};

watch(visible, newVal => {
  if (!newVal) {
    formData.value = { name: '', path: '', description: '', hidePath: false };
  }
});

async function handleCreate() {
  try {
    await formRef.value?.validate();
    loading.value = true;
    const project = await projectStore.createProject(formData.value);
    message.success('项目创建成功');
    visible.value = false;
    emit('success', project);
  } catch (error: any) {
    if (error?.message) {
      message.error(error.message);
    }
    return false;
  } finally {
    loading.value = false;
  }

  return true;
}
</script>
