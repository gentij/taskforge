export function buildPaginationMeta(params: {
  page: number;
  pageSize: number;
  total: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}) {
  const totalPages =
    params.total === 0 ? 0 : Math.ceil(params.total / params.pageSize);
  return {
    page: params.page,
    pageSize: params.pageSize,
    total: params.total,
    totalPages,
    hasNext: totalPages > 0 && params.page < totalPages,
    hasPrev: totalPages > 0 && params.page > 1,
    sortBy: params.sortBy,
    sortOrder: params.sortOrder,
  };
}
