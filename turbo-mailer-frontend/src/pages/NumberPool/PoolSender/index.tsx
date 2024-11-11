import React, {useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import {addPoolSender, removePoolSender, queryPoolSender, updatePoolSender} from './service';
import {Button, message, Modal, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {ProDescriptions} from "@ant-design/pro-components";
import {useParams} from "react-router";
import {queryNumberPoolById} from "@/pages/NumberPool/List/service";
const PoolSenderList: React.FC = () => {
  const [modalOpen, handleModalOpen] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<PoolSender.PoolSenderListItem>();
  const params  = useParams();

  /**
   * 添加或修改号池
   * @param fields 号池
   */
  const handleSubmit = async (fields: NumberPool.NumberPoolListItem) => {
    try {
      fields.id ? await updatePoolSender(params.poolId, fields) : await addPoolSender(params.poolId, fields);
      message.success(`${fields.id ? '编辑' : '添加'}成功`);
      return true;
    } catch (error) {
      message.error(`${fields.id ? '编辑' : '添加'}失败请重试!`);
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
      fieldProps: {
        rules: [
          {
            required: true,
            message: '请输入',
          },
        ],
      },
    },
    {
      title: '发件人地址',
      dataIndex: 'from_email',
      hideInSearch: true,
      fieldProps: {
        rules: [
          {
            required: true,
            whitespace: true,
            message: '请填写电子邮箱',
          },
          {
            type: 'email',
            message: '电子邮箱格式不正确',
          },
        ],
        style: { width: '100%' },
      },
    },
    {
      title: '回复地址',
      dataIndex: 'reply_to',
      hideInSearch: true,
      fieldProps: {
        rules: [
          {
            type: 'email',
            message: '电子邮箱格式不正确',
          },
        ],
        style: { width: '100%' },
      },
    },
    {
      title: '域名',
      dataIndex: 'domain',
      hideInSearch: true,
      hideInTable: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
      defaultSortOrder: 'descend',
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'dateTime',
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
            handleModalOpen(true);
            setCurrentRow(record);
          }}
          hidden
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
      <Modal
        title={`${currentRow?.id ? '更新' : '新建'}发件人`}
        open={modalOpen}
        onCancel={() => {
          handleModalOpen(false);
          setCurrentRow({});
        }}
        footer={null}
        destroyOnClose // 确保弹窗关闭时子组件被销毁
      >
        <ProTable<PoolSender.PoolSenderListItem, PoolSender.PoolSenderListItem>
          onSubmit={async (fields) => {
            const values = {
              ...fields,
              id: currentRow?.id,
            };
            const success = await handleSubmit(values as PoolSender.PoolSenderListItem);
            if (success) {
              handleModalOpen(false);
              setCurrentRow(undefined);
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
          rowKey="id"
          type="form"
          columns={columns}
          form={{
            initialValues: currentRow,
          }}
        />
      </Modal>
    </PageContainer>
  );
};

export default PoolSenderList;
