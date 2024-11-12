import { request } from 'umi';
import PoolSender from "@/pages/NumberPool/PoolSender/index";

export async function queryDashboardStats() {
  return request<Dashboard.Stats>('/api/v1/dashboard/stats');
}
