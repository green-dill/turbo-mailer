import React, {useEffect, useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import {addPoolSender, removePoolSender, queryPoolSender, updatePoolSender} from './service';
import {Button, message, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {
  ModalForm,
  ProDescriptions,
  ProFormText,
} from "@ant-design/pro-components";
import {useParams} from "react-router";
import {queryNumberPoolById} from "@/pages/NumberPool/List/service";


const PoolSenderList: React.FC = () => {
  /**
   * @en-US Pop-up window of new window
   * @zh-CN 新建窗口的弹窗
   *  */
  const [createModalOpen, handleModalOpen] = useState<boolean>(false);
  /**
   * @en-US The pop-up window of the distribution update window
   * @zh-CN 分布更新窗口的弹窗
   * */
  const [updateModalOpen, handleUpdateModalOpen] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<PoolSender.PoolSenderListItem>();

  const params  = useParams();

  /**
   * 添加发件人
   * @param fields 发件人
   */
  const handleAdd = async (fields: PoolSender.PoolSenderListItem) => {
    try {
      await addPoolSender(params.poolId, fields);
      message.success('添加成功');
      return true;
    } catch (error) {
      message.error('添加失败请重试！');
      return false;
    }
  };

  /**
   * 修改发件人
   * @param fields 发件人
   */
  const handleUpdate = async (fields: PoolSender.PoolSenderListItem) => {
    try {
      await updatePoolSender(params.poolId, fields);
      message.success('编辑成功');
      return true;
    } catch (error) {
      message.error('编辑失败请重试！');
      return false;
    }
  };

  /**
   *  删除发件人
   * @param fields 发件人
   */
  const handleRemove = async (fields: PoolSender.PoolSenderListItem) => {
    try {
      await removePoolSender(fields);
      message.success('删除成功');
      return true;
    } catch (error) {
      message.error('删除失败请重试！');
      return false;
    }
  };

  const columns: ProColumns<PoolSender.PoolSenderListItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      hideInForm: true,
      hideInTable: true,
    },
    {
      title: 'PoolID',
      dataIndex: 'pool_id',
      hideInSearch: true,
      hideInForm: true,
      hideInTable: true,
    },
    {
      title: '发件人名称',
      dataIndex: 'from_name',
      hideInSearch: true,
    },
    {
      title: '发件人地址',
      dataIndex: 'from_email',
      hideInSearch: true,
    },
    {
      title: '回复地址',
      dataIndex: 'reply_to',
      hideInSearch: true,
    },
    {
      title: '域名',
      dataIndex: 'domain',
      hideInSearch: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'date',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
      defaultSortOrder: 'descend',
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'date',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record: PoolSender.PoolSenderListItem) => [
        <Button
          type="link"
          size="small"
          style={{padding: 0}}
          key='edit'
          onClick={() => {
            handleUpdateModalOpen(true);
            setCurrentRow(record);
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
    <PageContainer
      content={
        <ProDescriptions<NumberPool.NumberPoolListItem>
          column={3}
          request={() => queryNumberPoolById({id: params.poolId})}
          columns={[
            {
              title: 'ID',
              dataIndex: 'id',
            },
            {
              title: '名称',
              dataIndex: 'name',
            },
            {
              title: '描述',
              dataIndex: 'description',
            },
            {
              title: '邮箱数量',
              dataIndex: 'sender_count',
            },
            {
              title: '创建时间',
              dataIndex: 'created_at',
              valueType: 'date',
            },
            {
              title: '更新时间',
              dataIndex: 'updated_at',
              valueType: 'date',
            },
          ]}
        />
      }
    >

      <ProTable<PoolSender.PoolSenderListItem, API.PageParams>
        headerTitle="发件人列表"
        actionRef={actionRef}
        rowKey="id"
        search={false}
        toolBarRender={() => [
          <Button type="primary" onClick={() => handleModalOpen(true)}>
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={() => queryPoolSender(params.poolId)}
        columns={columns}
      />
      <ModalForm
        title="添加发件人"
        width="400px"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        onFinish={async (value) => {
          const success = await handleAdd(value as PoolSender.PoolSenderListItem);
          if (success) {
            handleModalOpen(false);
            setCurrentRow(undefined);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="from_name"
          label="发件人名称"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="from_email"
          label="发件人地址"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="reply_to"
          label="回复地址"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="domain"
          label="域名"
        />
      </ModalForm>
      <ModalForm
        title="编辑发件人"
        width="400px"
        open={updateModalOpen}
        onOpenChange={handleUpdateModalOpen}
        initialValues={currentRow} // 设置初始值
        onFinish={async (value) => {
          const success = await handleUpdate(value as PoolSender.PoolSenderListItem);
          if (success) {
            handleUpdateModalOpen(false);
            setCurrentRow(undefined);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="from_name"
          label="发件人名称"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="from_email"
          label="发件人地址"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="reply_to"
          label="回复地址"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="domain"
          label="域名"
        />
      </ModalForm>
    </PageContainer>
  );
};

export default PoolSenderList;
