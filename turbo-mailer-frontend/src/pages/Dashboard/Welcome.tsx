import {PageContainer, ProCard, Statistic} from '@ant-design/pro-components';
import { useModel } from '@umijs/max';
import {Card, Col, Row, Space, theme} from 'antd';
import React, {useEffect, useState} from 'react';
import {queryDashboardStats} from "@/pages/Dashboard/service";

const Welcome: React.FC = () => {
  const { token } = theme.useToken();
  const { initialState } = useModel('@@initialState');
  const [dashBoardStats, setDashBoardStats] = useState<Dashboard.Stats>({});

  const fetchDashboardStats = async () => {
    return await queryDashboardStats();
  };

  useEffect(() => {
    fetchDashboardStats().then((response) => setDashBoardStats(response));
  }, []);


  return (dashBoardStats.topPools !== undefined &&
    <PageContainer>
      <Card
        style={{
          borderRadius: 8,
          marginBottom: 24,
        }}
        styles={{
          body: {
            backgroundImage:
              initialState?.settings?.navTheme === 'realDark'
                ? 'background-image: linear-gradient(75deg, #1A1B1F 0%, #191C1F 100%)'
                : 'background-image: linear-gradient(75deg, #FBFDFF 0%, #F5F7FF 100%)',
          },
        }}
      >
        <div
          style={{
            backgroundPosition: '100% -30%',
            backgroundRepeat: 'no-repeat',
            backgroundSize: '274px auto',
            backgroundImage:
              "url('https://gw.alipayobjects.com/mdn/rms_a9745b/afts/img/A*BuFmQqsB2iAAAAAAAAAAAAAAARQnAQ')",
          }}
        >
          <div
            style={{
              fontSize: '20px',
              color: token.colorTextHeading,
            }}
          >
            欢迎使用 Turbo Mailer
          </div>
          <p
            style={{
              fontSize: '14px',
              color: token.colorTextSecondary,
              lineHeight: '22px',
              marginTop: 16,
              marginBottom: 32,
              width: '65%',
            }}
          >
           Turbo Mailer 是一个高性能的营销邮件发送系统，其后端采用 Golang 构建，前端则使用 TypeScript/React。该系统通过 Kubernetes (k8s) 和 Helm 进行部署，为各种规模的企业提供了一个强大、经济高效且易于维护的解决方案。
          </p>
        </div>
      </Card>

      <Row gutter={24}>
        <Col span={6}>
          <Card title="号池总数" bordered={false}>
            <Space style={{ marginBottom: 8 }}>{dashBoardStats.poolCount}</Space>
          </Card>
        </Col>

        <Col span={6}>
          <Card title="号池排行" bordered={false}>
            <Space direction="vertical">
              {dashBoardStats.topPools.map((pool, index) => (
                <Space key={index} style={{ marginBottom: 8 }}>
                  {pool.name} - {pool.senderCount}
                </Space>
              ))}
            </Space>
          </Card>
        </Col>

        <Col span={6}>
          <Card title="任务总数" bordered={false}>
            <Space style={{ marginBottom: 8 }}>{dashBoardStats.taskCount}</Space>
          </Card>
        </Col>

        <Col span={6}>
          <Card title="任务状态统计" bordered={false}>
            <Space direction="vertical" style={{width: '100%'}}>
              <Space key='pending' style={{marginBottom: 8}}>
                待处理 - {dashBoardStats.taskStateCount?.pending}
              </Space>
              <Space key='dispatched' style={{ marginBottom: 8 }}>
                已调度 - {dashBoardStats.taskStateCount?.dispatched}
              </Space>
              <Space key='finished' style={{ marginBottom: 8 }}>
                已完成 - {dashBoardStats.taskStateCount?.finished}
              </Space>
            </Space>
          </Card>
        </Col>
      </Row>
    </PageContainer>
  );
};

export default Welcome;
