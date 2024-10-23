import { request } from 'umi';

export async function queryNumberPool(params: API.PageParams, options?: { [key: string]: any }) {
  return request<NumberPool.NumberPoolList>('/api/v1/pool', {
    method: 'GET',
    params,
    ...(options || {}),
  }).then(res => {
    return {
      total: res.total,
      data: res.list,
      success: true,
    }
  });
}

export async function updateNumberPool(body: NumberPool.NumberPoolListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/v1/pool', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function addNumberPool(body: NumberPool.NumberPoolListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/v1/pool', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function removeNumberPool(params: NumberPool.NumberPoolListItem, options?: { [key: string]: any }) {
  const { ID } = params;
  return request<Record<string, any>>(`/api/v1/pool/${ID}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}
