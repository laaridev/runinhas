import { useState, useEffect } from 'react';

interface CustomMessages {
  [key: string]: string;
}

const STORAGE_KEY = 'runinhas_custom_messages';

export function useCustomMessages() {
  const [messages, setMessages] = useState<CustomMessages>({});

  // Carregar mensagens do localStorage
  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
      try {
        setMessages(JSON.parse(stored));
      } catch (e) {
        console.error('Failed to load custom messages:', e);
      }
    }
  }, []);

  // Salvar mensagem customizada
  const setCustomMessage = (key: string, message: string) => {
    const updated = { ...messages, [key]: message };
    setMessages(updated);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
  };

  // Obter mensagem customizada ou padrão
  const getCustomMessage = (key: string, defaultMessage: string): string => {
    return messages[key] || defaultMessage;
  };

  // Resetar mensagem para padrão
  const resetMessage = (key: string) => {
    const updated = { ...messages };
    delete updated[key];
    setMessages(updated);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
  };

  // Resetar todas as mensagens
  const resetAllMessages = () => {
    setMessages({});
    localStorage.removeItem(STORAGE_KEY);
  };

  return {
    messages,
    setCustomMessage,
    getCustomMessage,
    resetMessage,
    resetAllMessages
  };
}
