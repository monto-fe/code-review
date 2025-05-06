// import { TableListQueryParams } from './data.d';

export async function hotSearchQueryList(): Promise<any> {
  return {
    data: {
      data: [
        {
          name: 'hot',
          hit: 123,
          id: 1,
        },
      ],
    },
  };
}

export async function hotTagsQueryList(): Promise<any> {
  return {
    data: {
      data: [
        {
          name: 'hot',
          hit: 456,
          id: 1,
        },
      ],
      total: 1,
    },
  };
}

export async function articleHitQueryList(): Promise<any> {
  return {
    data: {
      data: [
        {
          title: 'Vue 快速入门',
          hit: '999+',
          id: 1,
        },
      ],
      total: 1,
    },
  };
}

export async function worksHitQueryList(): Promise<any> {
  return {
    data: {
      data: [
        {
          title: 'colorful bookmark 谷歌书签',
          hit: '999+',
          id: 1,
        },
        {
          title: '时间戳转换 谷歌书签',
          hit: '999+',
          id: 2,
        },
      ],
      total: 1,
    },
  };
}
