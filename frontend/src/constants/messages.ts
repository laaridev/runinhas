/**
 * User-facing messages for the application
 */

export const TOAST_MESSAGES = {
  // Success messages
  AUDIO_GENERATED: 'Áudio gerado com sucesso!',
  MESSAGE_SAVED: 'Mensagem salva e áudio gerado com sucesso!',
  MESSAGE_RESTORED: 'Mensagem padrão restaurada com sucesso!',
  
  // Error messages
  AUDIO_ERROR: 'Erro ao gerar áudio. Tente novamente.',
  AUDIO_PREVIEW_ERROR: 'Erro ao testar voz',
  AUDIO_PLAY_ERROR: 'Erro ao reproduzir áudio',
  SAVE_ERROR: 'Erro ao salvar',
  RESTORE_ERROR: 'Erro ao restaurar',
  
  // Warning messages
  RETRY_GENERATING: (attempt: number, max: number) => 
    `Tentando gerar áudio novamente (${attempt}/${max})...`,
  RETRY_SAVING: (attempt: number, max: number) => 
    `Salvando mensagem (tentativa ${attempt}/${max})...`,
  MESSAGE_MISSING_PLACEHOLDER: 'A mensagem deve conter {seconds} para incluir o tempo!',
  MESSAGE_EMPTY: 'A mensagem não pode estar vazia!',
  
  // Info messages
  GENERATING_FIRST_TIME: 'Gerando áudio pela primeira vez...',
  RESTORING_DEFAULT: (attempt: number, max: number) => 
    `Restaurando mensagem padrão (${attempt}/${max})...`,
  GENERATING_DEFAULT: (attempt: number, max: number) => 
    `Gerando áudio padrão (${attempt}/${max})...`,
} as const;

export const VALIDATION_MESSAGES = {
  PLACEHOLDER_REQUIRED: (placeholder: string) => 
    `A mensagem deve conter ${placeholder} para incluir o tempo!`,
  EMPTY_MESSAGE: 'A mensagem não pode estar vazia!',
} as const;
