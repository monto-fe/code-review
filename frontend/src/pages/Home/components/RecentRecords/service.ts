import request from '@/utils/request';
import dayjs from 'dayjs';

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
    url: '/ai/message/list',
    method: 'post',
    data: {
      current: 1,
      page_size: 1000,
      start_date: dayjs().subtract(30, 'day').startOf('day').unix(),
      end_date: dayjs().endOf('day').unix()
    }
  });
} 