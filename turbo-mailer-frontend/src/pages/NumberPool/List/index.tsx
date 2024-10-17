import React, {useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { addNumberPool, removeNumberPool, queryNumberPool, updateNumberPool, } from './service';
import {Button, message, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {ModalForm, ProFormText, ProFormTextArea} from "@ant-design/pro-components";

const NumberPoolList: React.FC = () => {
  const [createModalOpen, handleModalOpen] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();

  /**
   * 添加或修改号池
   * @param fields 号池
   */
  const handleSubmit = async (fields: NumberPool.NumberPoolListItem) => {
    try {
      const response = fields.id ? await updateNumberPool(fields) : await addNumberPool(fields);
      if (response && response.code === 200) {
        message.success('编辑成功');
      } else {
        message.error('编辑失败请重试！');
      }
      return true;
    } catch (error) {
      message.error('添加失败请重试！');
      return false;
    }
  };

  /**
   *  删除号池
   * @param fields 号池
   */
  const handleRemove = async (fields: NumberPool.NumberPoolListItem) => {
    try {
      const response = await removeNumberPool(fields);
      if (response && response.code === 200) {
        message.success('删除成功');
      }
      return true;
    } catch (error) {
      message.error('删除失败请重试！');
      return false;
    }
  };

  const columns: ProColumns<NumberPool.NumberPoolListItem>[] = [
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '描述',
      dataIndex: 'description',
      hideInSearch: true,
    },
    {
      title: '邮箱数量',
      dataIndex: 'quantity',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueEnum: {
        0: {
          text: '正常',
          status: 'Success',
        },
        1: {
          text: '禁用',
          status: 'Error',
        },
      },
      filters: true,
      onFilter: true,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      valueType: 'date',
      hideInSearch: true,
      sorter: true,
      defaultSortOrder: 'descend',
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      valueType: 'date',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record: NumberPool.NumberPoolListItem) => [
        <Button
          type="link"
          size="small"
          style={{padding: 0}}
          key='edit'
          onClick={() => {
            handleModalOpen(true);
          }}
        >
          编辑
        </Button>,
        <Popconfirm
          key='delete'
          title="确定删除吗？"
          onConfirm={async () => {
            const success = await handleRemove(record);
            if (success) {
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
        >
          <Button style={{padding: 0}} type="link" size="small">删除</Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<NumberPool.NumberPoolListItem, API.PageParams>
        headerTitle="号池列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button type="primary" onClick={() => handleModalOpen(true)}>
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={queryNumberPool}
        columns={columns}
      />
      <ModalForm
        title='创建号池'
        width="400px"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        onFinish={async (value) => {
          const success = await handleSubmit(value as NumberPool.NumberPoolListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          label='号池名称'
          placeholder='请输入至少五个字符'
          rules={[
            {
              required: true,
              message: '号池名称不能为空',
            },
          ]}
          width="md"
          name="name"
        />
        <ProFormTextArea
          label='号池描述'
          width="md"
          name="description"
        />
      </ModalForm>
    </PageContainer>
  );
};

export default NumberPoolList;
