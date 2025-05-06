// import request from '@/utils/request';

export async function weeknewWorks(): Promise<any> {
  return {
    data: {
      total: 100,
      num: 34,
      chart: {
        day: [1, 2, 3, 4, 5, 6, 67, 78],
        num: [66, 21, 10, 14, 50, 80, 100, 18],
      },
    },
  };
}
