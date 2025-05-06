import dayjs from 'dayjs';

// 时间格式
export const timeFormatType = {
  date_no_year: 'MM-DD',
  date: 'YYYY-MM-DD',
  month: 'YYYY-MM',
  time: 'YYYY-MM-DD HH:mm:ss',
  ch_date_no_year: 'MM月DD日',
  ch_date: 'YYYY年MM月DD日',
  date_no_split: 'YYYYMMDD',
  time_xg: 'YYYY/MM/DD HH:mm',
  month_no_split: 'YYYYMM',
  date_xg: 'YYYY/MM/DD',
  date_point: 'YYYY.MM.DD',
};

// 获取指定格式的时间
export function renderDate(time: number | Date | dayjs.Dayjs, type: string) {
  return time ? dayjs(time).format(type || timeFormatType.time) : '-';
}

// 获取指定格式的时间区间
export function renderDateRange(times: number[] | Date[] | dayjs.Dayjs[], type: string) {
  return `${renderDate(times[0], type)} ~ ${renderDate(times[1], type)}`;
}

// 区分秒级和毫秒级时间戳来格式化
export function renderDateFromTimestamp(timestamp: number, type: string) {
  return renderDate(`${timestamp}`.length === 10 ? timestamp * 1000 : timestamp, type);
}

// 获取时间区间，不直接用object是因为这里定义object之后，获取key值就不会获取到实时的时间区间
export const dateRange: any = {
  '2H': [dayjs().subtract(2, 'hours'), dayjs()],
  '1天': [dayjs().startOf('day'), dayjs().endOf('day')],
  '3天': [dayjs().subtract(2, 'day').startOf('day'), dayjs().endOf('day')],
  '7天': [dayjs().subtract(6, 'day').startOf('day'), dayjs().endOf('day')],
  '15天': [dayjs().subtract(14, 'day').startOf('day'), dayjs().endOf('day')],
  '30天': [dayjs().subtract(29, 'day').startOf('day'), dayjs().endOf('day')],
  '90天': [dayjs().subtract(89, 'day').startOf('day'), dayjs().endOf('day')],
  昨天: [dayjs().subtract(1, 'day').startOf('day'), dayjs().subtract(1, 'day').endOf('day')],
  当日: [dayjs().startOf('day'), dayjs().endOf('day')],
  上周: [
    dayjs().subtract(1, 'week').startOf('week').startOf('day'),
    dayjs().subtract(1, 'week').endOf('week').endOf('day'),
  ],
  本周: [dayjs().startOf('week').startOf('day'), dayjs().endOf('day')],
  一周: [dayjs().subtract(1, 'week').startOf('day'), dayjs().endOf('day')],
  上月: [
    dayjs().subtract(1, 'month').startOf('month').startOf('day'),
    dayjs().subtract(1, 'month').endOf('month').endOf('day'),
  ],
  上上月: [
    dayjs().subtract(2, 'month').startOf('month').startOf('day'),
    dayjs().subtract(2, 'month').endOf('month').endOf('day'),
  ],
  本月: [dayjs().startOf('month').startOf('day'), dayjs().endOf('day')],
  一月: [dayjs().subtract(1, 'month').startOf('day'), dayjs().endOf('day')],
  三月: [dayjs().subtract(3, 'month').startOf('day'), dayjs().endOf('day')],
  本年: [dayjs().startOf('year').startOf('day'), dayjs().endOf('day')],
  一年: [dayjs().subtract(12, 'month').startOf('day'), dayjs().endOf('day')],
  前3月: [
    dayjs().subtract(3, 'month').startOf('month').startOf('day'),
    dayjs().subtract(1, 'month').endOf('month').endOf('day'),
  ],
  前6月: [
    dayjs().subtract(6, 'month').startOf('month').startOf('day'),
    dayjs().subtract(1, 'month').endOf('month').endOf('day'),
  ],
  后7天: [dayjs().add(1, 'day').startOf('day'), dayjs().add(7, 'days').endOf('day')],
};

export const getDateRange = (dateKey: string) => dateRange[dateKey];

export const getCustomDateRange = (times: number[] | Date[] | dayjs.Dayjs[]) => {
  if (times && times.length === 2) {
    return [dayjs(times[0]).startOf('day'), dayjs(times[1]).endOf('day')];
  }
  return times;
};
