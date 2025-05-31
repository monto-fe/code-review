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

export async function getRecentRecords(): Promise<{ data: RecordItem[] }> {
  return request({
    url: '/api/records/recent',
    method: 'GET',
  });
} 