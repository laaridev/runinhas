export interface RuneConfig {
  key: string;
  name: string;
  icon: string;
  description: string;
  min: number;
  max: number;
  step: number;
  defaultValue: number;
}

export interface InstallResult {
  success: boolean;
  message: string;
  installedAt: string;
}
