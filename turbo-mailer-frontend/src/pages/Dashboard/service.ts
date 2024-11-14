import { request } from 'umi';

export async function queryDashboardStats() {
  return request<Dashboard.Stats>('/api/v1/dashboard/stats');
}
