import { request } from 'umi';

export async function queryNumberPool(params: API.PageParams, options?: { [key: string]: any }) {
  return request<NumberPool.NumberPoolList>('/number-pool', {
    method: 'GET',
    params: params,
    ...(options || {}),
  });
}

export async function updateNumberPool(body: NumberPool.NumberPoolListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/number-pool', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function addNumberPool(body: NumberPool.NumberPoolListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/number-pool', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function removeNumberPool(params: NumberPool.NumberPoolListItem, options?: { [key: string]: any }) {
  const { id } = params;
  return request<Record<string, any>>(`/number-pool/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}
