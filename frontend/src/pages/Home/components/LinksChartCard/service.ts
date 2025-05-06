// import request from '@/utils/request';

export async function annualnewLinks(): Promise<any> {
  return {
    data: {
      total: 100,
      num: 34,
      chart: {
        day: [5, 4, 32, 1],
        num: [1, 2, 3, 4, 5],
      },
    },
  };
}
