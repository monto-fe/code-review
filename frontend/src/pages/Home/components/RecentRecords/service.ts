import request from '@/utils/request';

export interface RecordItem {
  id: string;
  type: 'commit' | 'review';
  title: string;
  author: string;
  status: 'pending' | 'approved' | 'rejected';
  createdAt: string;
  url: string;
}

export async function getRecentRecords(): Promise<any> {
  return request({
    url: '/ai/message',
    method: 'GET',
  });
} 