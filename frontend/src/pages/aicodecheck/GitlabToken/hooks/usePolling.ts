import { useEffect, useRef, useCallback } from 'react';
import { queryList } from '../service';

interface UsePollingOptions {
  tokenId: number;
  onSuccess?: () => void;
  onError?: () => void;
  onTimeout?: () => void;
  interval?: number; // 轮询间隔，默认3秒
  timeout?: number; // 超时时间，默认60秒
}

export const usePolling = ({
  tokenId,
  onSuccess,
  onError,
  onTimeout,
  interval = 3000,
  timeout = 60000
}: UsePollingOptions) => {
  const intervalRef = useRef<number>();
  const timeoutRef = useRef<number>();
  const startTimeRef = useRef<number>(Date.now());

  const stopPolling = useCallback(() => {
    if (intervalRef.current) {
      clearTimeout(intervalRef.current);
      intervalRef.current = undefined;
    }
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = undefined;
    }
  }, []);

  const startPolling = useCallback(() => {
    // 清除之前的定时器
    stopPolling();
    
    // 记录开始时间
    startTimeRef.current = Date.now();

    // 设置超时定时器
    timeoutRef.current = setTimeout(() => {
      stopPolling();
      onTimeout?.();
    }, timeout);

    // 开始轮询
    const poll = async () => {
      try {
        const response = await queryList();
        const token = response.data?.find((item: any) => item.id === tokenId);
        
        if (token) {
          // 检查同步状态：1-失败，2-同步中，3-成功
          if (token.project_ids_synced !== 2) {
            // 同步完成或失败，停止轮询
            stopPolling();
            if (token.project_ids_synced === 3) {
              onSuccess?.();
            } else if (token.project_ids_synced === 1) {
              onError?.();
            }
            return;
          }
        } else {
          // 如果找不到对应的 token，可能是刚创建还没同步到列表
          // 继续轮询
        }
        
        // 继续轮询
        intervalRef.current = setTimeout(poll, interval);
      } catch (error) {
        console.error('轮询失败:', error);
        // 出错时继续轮询，不停止
        intervalRef.current = setTimeout(poll, interval);
      }
    };

    // 立即执行一次
    poll();
  }, [tokenId, interval, timeout, stopPolling, onSuccess, onError, onTimeout]);

  // 组件卸载时清理定时器
  useEffect(() => {
    return () => {
      stopPolling();
    };
  }, [stopPolling]);

  return {
    startPolling,
    stopPolling,
    isPolling: !!intervalRef.current
  };
}; 