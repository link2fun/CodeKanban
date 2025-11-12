import { computed, onBeforeUnmount, ref } from 'vue';

const panelOrder = ref<string[]>([]);

type UsePanelStackOptions = {
  baseZIndex?: number;
};

function registerPanel(panelKey: string) {
  if (!panelOrder.value.includes(panelKey)) {
    panelOrder.value = [...panelOrder.value, panelKey];
  }
}

export function usePanelStack(panelKey: string, options?: UsePanelStackOptions) {
  const baseZIndex = options?.baseZIndex ?? 1000;

  registerPanel(panelKey);

  const bringToFront = () => {
    registerPanel(panelKey);
    panelOrder.value = [...panelOrder.value.filter(key => key !== panelKey), panelKey];
  };

  const zIndex = computed(() => {
    const index = panelOrder.value.indexOf(panelKey);
    return baseZIndex + (index === -1 ? 0 : index);
  });

  onBeforeUnmount(() => {
    panelOrder.value = panelOrder.value.filter(key => key !== panelKey);
  });

  return {
    zIndex,
    bringToFront,
  };
}
