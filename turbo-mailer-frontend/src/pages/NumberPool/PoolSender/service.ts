import { request } from 'umi';
import PoolSender from "@/pages/NumberPool/PoolSender/index";

export async function queryPoolSender(poolId?: string, params?: API.PageParams, options?: { [key: string]: any }) {
  return request<NumberPool.NumberPoolList>(`/api/v1/pool-senders/${poolId}`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

export async function updatePoolSender(poolId?: string, body?: PoolSender.PoolSenderListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>(`/api/v1/pool-senders/${body?.id}`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

export async function addPoolSender(poolId?: string, body?: PoolSender.PoolSenderListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>(`/api/v1/pool-senders/${poolId}`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

export async function removePoolSender(params: PoolSender.PoolSenderListItem, options?: { [key: string]: any }) {
  const { id } = params;
  return request<Record<string, any>>(`/api/v1/pool-senders/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}
