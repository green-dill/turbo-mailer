import React from 'react';
import {Button, Typography, Badge, Space} from '@arco-design/web-react';
import dayjs from 'dayjs';

const { Text } = Typography;
export const Status = ['禁用', '正常'];

export function getColumns(
  t: any,
  callback: (record: Record<string, any>, type: string) => Promise<void>
) {
  return [
    {
      title: t['searchTable.columns.id'],
      dataIndex: 'id',
      render: (value) => <Text copyable>{value}</Text>,
    },
    {
      title: t['searchTable.columns.name'],
      dataIndex: 'name',
    },
    {
      title: t['searchTable.columns.contentNum'],
      dataIndex: 'count',
      sorter: (a, b) => a.count - b.count,
      render(x) {
        return Number(x).toLocaleString();
      },
    },
    {
      title: t['searchTable.columns.createdAt'],
      dataIndex: 'createdAt',
      render: (x) => dayjs().subtract(x, 'days').format('YYYY-MM-DD HH:mm:ss'),
      sorter: (a, b) => b.createdTime - a.createdTime,
    },
    {
      title: t['searchTable.columns.status'],
      dataIndex: 'status',
      render: (x) => {
        if (x === 0) {
          return <Badge status="error" text={Status[x]}></Badge>;
        }
        return <Badge status="success" text={Status[x]}></Badge>;
      },
    },
    {
      title: t['searchTable.columns.operations'],
      dataIndex: 'operations',
      headerCellStyle: { paddingLeft: '15px' },
      render: (_, record) => (
        <Space>
          <Button
            type="text"
            size="small"
            onClick={() => callback(record, 'view')}
          >
            {t['searchTable.columns.operations.view']}
          </Button>
          <Button
            type="text"
            size="small"
            onClick={() => callback(record, 'remove')}
          >
            {t['searchTable.columns.operations.remove']}
          </Button>
        </Space>
      ),
    },
  ];
}
