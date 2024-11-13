import React, {useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { addNumberPool, removeNumberPool, queryNumberPool, updateNumberPool, } from './service';
import {Button, message, Modal, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {Link} from "@@/exports";

const NumberPoolList: React.FC = () => {
  const [modalOpen, handleModalOpen] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<NumberPool.NumberPoolListItem>();

  /**
   * 添加或修改号池
   * @param fields 号池
   */
  const handleSubmit = async (fields: NumberPool.NumberPoolListItem) => {
    try {
      fields.id ? await updateNumberPool(fields) : await addNumberPool(fields);
      message.success(`${fields.id ? '编辑' : '添加'}成功`);
      return true;
    } catch (error) {
      message.error(`${fields.id ? '编辑' : '添加'}失败请重试!`);
      return false;
    }
  };

  /**
   *  删除号池
   * @param fields 号池
   */
  const handleRemove = async (fields: NumberPool.NumberPoolListItem) => {
    try {
      await removeNumberPool(fields);
      message.success('删除成功');
      return true;
    } catch (error) {
      message.error('删除失败请重试！');
      return false;
    }
  };

  const columns: ProColumns<NumberPool.NumberPoolListItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      hideInForm: true,
      hideInTable: true,
    },
    {
      title: '名称',
      dataIndex: 'name',
      formItemProps: {
        rules: [
          {
            required: true,
            message: '请填写',
          },
        ],
      },
    },
    {
      title: '描述',
      dataIndex: 'description',
      valueType: 'textarea',
      hideInSearch: true,
    },
    {
      title: '邮箱数量',
      dataIndex: 'sender_count',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
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
      render: (_, record: NumberPool.NumberPoolListItem) => [
        <Link key="setting" to={`/number-pool/list/${record.id}`}>
          配置
        </Link>,
        <Button
          type="link"
          size="small"
          style={{padding: 0}}
          key='edit'
          onClick={() => {
            handleModalOpen(true);
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
        columnsState={{
          persistenceKey: 'number-pool-list',
          persistenceType: 'localStorage',
          defaultValue: {
            option: {fixed: 'right', disable: true},
          },
        }}
        columns={columns}
      />
      <Modal
        title={`${currentRow?.id ? '更新' : '新建'}号池`}
        open={modalOpen}
        onCancel={() => {
          handleModalOpen(false);
          setCurrentRow({});
        }}
        footer={null}
        destroyOnClose // 确保弹窗关闭时子组件被销毁
      >
        <ProTable<NumberPool.NumberPoolListItem, NumberPool.NumberPoolListItem>
          onSubmit={async (fields) => {
            const values = {
              ...fields,
              id: currentRow?.id,
            };
            const success = await handleSubmit(values as NumberPool.NumberPoolListItem);
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

export default NumberPoolList;
