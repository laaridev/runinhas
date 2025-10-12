import { create } from 'zustand';
import { v4 as uuidv4 } from 'uuid';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
  id: string;
  title?: string;
  message: string;
  type: ToastType;
  duration?: number;
}

interface ToastStore {
  toasts: Toast[];
  addToast: (toast: Omit<Toast, 'id'>) => void;
  removeToast: (id: string) => void;
  clearToasts: () => void;
}

export const useToastStore = create<ToastStore>((set) => ({
  toasts: [],
  
  addToast: (toast) => {
    const id = uuidv4();
    const newToast = { ...toast, id };
    
    set((state) => ({
      toasts: [...state.toasts, newToast]
    }));
    
    // Auto-remove after duration
    const duration = toast.duration || 3000;
    if (duration > 0) {
      setTimeout(() => {
        set((state) => ({
          toasts: state.toasts.filter(t => t.id !== id)
        }));
      }, duration);
    }
  },
  
  removeToast: (id) => {
    set((state) => ({
      toasts: state.toasts.filter(t => t.id !== id)
    }));
  },
  
  clearToasts: () => {
    set({ toasts: [] });
  }
}));

// Hook conveniente
export function useToast() {
  const { addToast, removeToast, clearToasts } = useToastStore();
  
  return {
    success: (message: string, title?: string) => 
      addToast({ message, title, type: 'success' }),
    
    error: (message: string, title?: string) => 
      addToast({ message, title, type: 'error' }),
    
    info: (message: string, title?: string) => 
      addToast({ message, title, type: 'info' }),
    
    warning: (message: string, title?: string) => 
      addToast({ message, title, type: 'warning' }),
    
    dismiss: removeToast,
    dismissAll: clearToasts
  };
}
