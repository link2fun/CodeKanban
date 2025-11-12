import { defineStore } from 'pinia';
import { computed, ref, watch } from 'vue';

export interface ThemeSettings {
  primaryColor: string;
  surfaceColor: string;
  bodyColor: string;
}

interface GeneralSettings {
  theme: ThemeSettings;
}

const STORAGE_KEY = 'general_settings';

const defaultTheme: ThemeSettings = {
  primaryColor: '#18a058',
  surfaceColor: '#ffffff',
  bodyColor: '#f7f8fa',
};

const defaultSettings: GeneralSettings = {
  theme: { ...defaultTheme },
};

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<GeneralSettings>(loadSettings());

  const theme = computed(() => settings.value.theme);

  watch(
    settings,
    newSettings => {
      saveSettings(newSettings);
    },
    { deep: true },
  );

  function updateTheme(partial: Partial<ThemeSettings>) {
    settings.value.theme = {
      ...settings.value.theme,
      ...partial,
    };
  }

  function resetTheme() {
    settings.value.theme = { ...defaultTheme };
  }

  return {
    theme,
    updateTheme,
    resetTheme,
  };
});

function loadSettings(): GeneralSettings {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
      const parsed = JSON.parse(stored) as Partial<GeneralSettings>;
      return {
        theme: {
          ...defaultTheme,
          ...(parsed.theme ?? {}),
        },
      };
    }
  } catch (error) {
    console.warn('Failed to load settings, falling back to defaults.', error);
  }
  return cloneDefaultSettings();
}

function saveSettings(settings: GeneralSettings) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
  } catch (error) {
    console.error('Failed to persist settings:', error);
  }
}

function cloneDefaultSettings(): GeneralSettings {
  return {
    theme: { ...defaultTheme },
  };
}
