export function pageCompute(current: number, pageSize: number) {
  if (isNaN(current)) {
    current = 1;
  }
  if (isNaN(pageSize)) {
    pageSize = 10;
  }
  if (pageSize > 100000) {
    pageSize = 100000;
  }
  return {
    offset: (Number(current) - 1) * Number(pageSize),
    limit: Number(pageSize),
  }
}