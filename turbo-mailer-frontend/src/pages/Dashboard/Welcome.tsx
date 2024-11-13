import {PageContainer, StatisticCard} from '@ant-design/pro-components';
import RcResizeObserver from 'rc-resize-observer';
import { useModel } from '@umijs/max';
import {Card, Col, Divider, Row, theme} from 'antd';
import React, {useEffect, useState} from 'react';
import {queryDashboardStats} from "@/pages/Dashboard/service";

const Welcome: React.FC = () => {
  const { token } = theme.useToken();
  const { initialState } = useModel('@@initialState');
  const [dashBoardStats, setDashBoardStats] = useState<Dashboard.Stats>({});
  const [responsive, setResponsive] = useState(false);

  const fetchDashboardStats = async () => {
    return await queryDashboardStats();
  };

  useEffect(() => {
    fetchDashboardStats().then((response) => setDashBoardStats(response));
  }, []);

  return (dashBoardStats.topPools !== undefined &&
    <PageContainer>
      <RcResizeObserver
        key="resize-observer"
        onResize={(offset) => {
          setResponsive(offset.width < 596);
        }}
      >
        <Row gutter={24}>
          <Col span={8}>
            <StatisticCard.Group direction={responsive ? 'column' : 'row'} style={{
              borderRadius: 8,
              marginBottom: 24,
            }}>
              <StatisticCard
                statistic={{
                  title: '号池数',
                  tip: '号池',
                  value: dashBoardStats.poolCount,
                }}
              />
              <Divider type={responsive ? 'horizontal' : 'vertical'} />
              <StatisticCard
                statistic={{
                  title: '热门号池',
                  value: dashBoardStats.topPools[0].name
                }}
              />
            </StatisticCard.Group>
          </Col>
          <Col span={16}>
            <StatisticCard.Group direction={responsive ? 'column' : 'row'} style={{
              borderRadius: 8,
              marginBottom: 24,
            }}>
              <StatisticCard
                statistic={{
                  title: '任务数',
                  tip: '帮助文字',
                  value: dashBoardStats.taskCount,
                }}
              />
              <Divider type={responsive ? 'horizontal' : 'vertical'} />
              <StatisticCard
                statistic={{
                  title: '待处理',
                  value: dashBoardStats.taskStateCount?.pending || 0,
                  status: 'default',
                }}
              />
              <StatisticCard
                statistic={{
                  title: '已调度',
                  value: dashBoardStats.taskStateCount?.dispatched || 0,
                  status: 'processing',
                }}
              />
              <StatisticCard
                statistic={{
                  title: '已完成',
                  value: dashBoardStats.taskStateCount?.finished || 0,
                  status: 'success',
                }}
              />
            </StatisticCard.Group>
          </Col>
        </Row>
      </RcResizeObserver>
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
              paddingBottom: 128,
              width: '65%',
            }}
          >
           Turbo Mailer 是一个高性能的营销邮件发送系统，其后端采用 Golang 构建，前端则使用 TypeScript/React。该系统通过 Kubernetes (k8s) 和 Helm 进行部署，为各种规模的企业提供了一个强大、经济高效且易于维护的解决方案。
          </p>
        </div>
      </Card>
    </PageContainer>
  );
};

export default Welcome;
