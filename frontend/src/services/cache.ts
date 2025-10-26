// Cache service for offline mode
interface CacheEntry<T> {
  data: T;
  timestamp: number;
  ttl?: number; // Time to live in milliseconds
}

class CacheService {
  private readonly CACHE_PREFIX = 'dota_gsi_cache_';
  
  // Save to cache
  set<T>(key: string, data: T, ttl?: number): void {
    const entry: CacheEntry<T> = {
      data,
      timestamp: Date.now(),
      ttl
    };
    
    try {
      localStorage.setItem(
        this.CACHE_PREFIX + key,
        JSON.stringify(entry)
      );
    } catch (error) {
      console.error('Failed to save to cache:', error);
    }
  }
  
  // Get from cache
  get<T>(key: string): T | null {
    try {
      const item = localStorage.getItem(this.CACHE_PREFIX + key);
      if (!item) return null;
      
      const entry: CacheEntry<T> = JSON.parse(item);
      
      // Check if expired
      if (entry.ttl) {
        const age = Date.now() - entry.timestamp;
        if (age > entry.ttl) {
          this.remove(key);
          return null;
        }
      }
      
      return entry.data;
    } catch (error) {
      console.error('Failed to get from cache:', error);
      return null;
    }
  }
  
  // Remove from cache
  remove(key: string): void {
    localStorage.removeItem(this.CACHE_PREFIX + key);
  }
  
  // Clear all cache
  clear(): void {
    const keys = Object.keys(localStorage);
    keys.forEach(key => {
      if (key.startsWith(this.CACHE_PREFIX)) {
        localStorage.removeItem(key);
      }
    });
  }
  
  // Check if cache exists and is valid
  has(key: string): boolean {
    return this.get(key) !== null;
  }
  
  // Get or set pattern
  async getOrSet<T>(
    key: string,
    fetcher: () => Promise<T>,
    ttl?: number
  ): Promise<T> {
    // Try cache first
    const cached = this.get<T>(key);
    if (cached !== null) {
      return cached;
    }
    
    // Fetch and cache
    try {
      const data = await fetcher();
      this.set(key, data, ttl);
      return data;
    } catch (error) {
      // If fetch fails, try to return stale cache
      const stale = this.get<T>(key);
      if (stale !== null) {
        console.warn('Using stale cache due to fetch error');
        return stale;
      }
      throw error;
    }
  }
}

export const cache = new CacheService();

// Hook for cached API calls
import { useState, useEffect } from 'react';

export function useCachedData<T>(
  key: string,
  fetcher: () => Promise<T>,
  ttl: number = 5 * 60 * 1000 // 5 minutes default
) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  
  useEffect(() => {
    let mounted = true;
    
    const loadData = async () => {
      try {
        setLoading(true);
        const result = await cache.getOrSet(key, fetcher, ttl);
        if (mounted) {
          setData(result);
          setError(null);
        }
      } catch (err) {
        if (mounted) {
          setError(err as Error);
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };
    
    loadData();
    
    return () => {
      mounted = false;
    };
  }, [key]);
  
  const refresh = async () => {
    cache.remove(key);
    setLoading(true);
    try {
      const result = await fetcher();
      cache.set(key, result, ttl);
      setData(result);
      setError(null);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  };
  
  return { data, loading, error, refresh };
}
