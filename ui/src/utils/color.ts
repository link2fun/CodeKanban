const HEX_PATTERN = /^#?[0-9a-f]{3}([0-9a-f]{3})?$/i;

function normalizeHex(hex: string): string | null {
  if (!HEX_PATTERN.test(hex)) {
    return null;
  }
  let normalized = hex.startsWith('#') ? hex.slice(1) : hex;
  if (normalized.length === 3) {
    normalized = normalized
      .split('')
      .map(char => char + char)
      .join('');
  }
  return normalized.toLowerCase();
}

function clamp(value: number) {
  return Math.max(0, Math.min(255, value));
}

function toHex(value: number) {
  return value.toString(16).padStart(2, '0');
}

function adjustColor(hex: string, ratio: number): string {
  const normalized = normalizeHex(hex);
  if (!normalized) {
    return hex;
  }
  const r = parseInt(normalized.slice(0, 2), 16);
  const g = parseInt(normalized.slice(2, 4), 16);
  const b = parseInt(normalized.slice(4, 6), 16);

  const amount = Math.round(255 * ratio);

  const nextR = clamp(r + amount);
  const nextG = clamp(g + amount);
  const nextB = clamp(b + amount);

  return `#${toHex(nextR)}${toHex(nextG)}${toHex(nextB)}`;
}

export function lightenColor(hex: string, ratio = 0.1): string {
  return adjustColor(hex, Math.abs(ratio));
}

export function darkenColor(hex: string, ratio = 0.1): string {
  return adjustColor(hex, -Math.abs(ratio));
}
