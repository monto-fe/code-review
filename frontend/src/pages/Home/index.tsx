import { useCallback, useEffect, useState } from 'react';
import { Card, Row, Col, Typography, Spin } from 'antd';
// import { BasicContext } from '@/store/context';
import { GuideStepsForm, UpdateLog, ProductAbility } from './components';
import MergeStats from './components/MergeStats';
import MergeProblemStats from './components/MergeStats/MergeProblemStats';
import AIDetectionChart from './components/AIDetectionChart';
import RecentRecords from './components/RecentRecords';
import { SearchOutlined, BugOutlined } from '@ant-design/icons';
import DocLinksCard from './components/DocLinksCard';
import AlertTips from './components/AlertTips';
import { queryList } from '@/pages/aicodecheck/AIModelManager/service';
import { queryList as getGitlabTokenList } from '@/pages/aicodecheck/GitlabToken/service';
import BilibiliVideoCard from './components/BilibiliVideoCard';

const { Title } = Typography;

function Dashboard() {
  const [AIModalList, setAIModalList] = useState<any>([]);
  const [GitlabToken, setGitlabToken] = useState<any>([]);
  const [AIModalLoading, setAIModalLoading] = useState<boolean>(false);
  const [GitlabTokenLoading, setGitlabTokenLoading] = useState<boolean>(false);

  const showConfigAlert = AIModalList.length === 0 || GitlabToken.length === 0;

  const getAIModalList = useCallback(() => {
    setAIModalLoading(true);
    queryList().then((res) => {
      const { data: { data } } = res;
      setAIModalList(data);
      setAIModalLoading(false);
    }).catch((err) => {
      console.log("err", err);
      setAIModalList([]);
      setAIModalLoading(false);
    });
  }, []);

  const getGitlabToken = useCallback(() => {
    setGitlabTokenLoading(true);
    getGitlabTokenList().then((res) => {
      const { data: { data} } = res;
      setGitlabToken(data);
      setGitlabTokenLoading(false);
    }).catch((err) => {
      console.log("err", err);
      setGitlabToken([]);
      setGitlabTokenLoading(false);
    });
  }, []);

  useEffect(() => {
    getAIModalList();
    getGitlabToken();
  }, [getAIModalList, getGitlabToken]);


  return (
    <Spin spinning={AIModalLoading || GitlabTokenLoading}>
    <div className="layout-main-conent" style={{ background: '#fafbfc', minHeight: '100vh', padding: 24 }}>
      <AlertTips AIConfig={AIModalList.length > 0} GitlabConfig={GitlabToken.length > 0} />
      
      <Row gutter={32}>
        {/* 主内容区 */}
        {!showConfigAlert && <Col xs={24} md={16}>
          {/* 统计卡片 */}
          <Row gutter={24} style={{ marginBottom: 24 }}>
            <Col span={12}>
              <Card>
                <Row gutter={32}>
                  <Col flex="32px">
                    <SearchOutlined style={{ fontSize: 32, color: '#1890ff' }} />
                  </Col>
                  <Col flex="auto">
                    <div style={{ fontSize: 16, color: '#666' }}>检测次数（近30天）</div>
                    <div style={{ fontSize: 36, fontWeight: 600, margin: '8px 0' }}>
                      <MergeStats />
                    </div>
                  </Col>
                </Row>
              </Card>
            </Col>
            <Col span={12}>
              <Card>
                <Row gutter={32}>
                  <Col flex="32px">
                    <BugOutlined style={{ fontSize: 32, color: '#fa541c' }} />
                  </Col>
                  <Col flex="auto">
                    <div style={{ fontSize: 16, color: '#666' }}>发现问题数（近30天）</div>
                    <div style={{ fontSize: 36, fontWeight: 600, margin: '8px 0' }}>
                      <MergeProblemStats />
                    </div>
                  </Col>
                </Row>
              </Card>
            </Col>
          </Row>
          {/* AI检测效果 */}
          <Card style={{ marginBottom: 24 }}>
            <Title level={5} style={{ marginBottom: 16 }}>AI检测效果（近30天）</Title>
            <AIDetectionChart />
          </Card>
          {/* 最近提交和审核记录 */}
            <RecentRecords />
        </Col>}
        {showConfigAlert && <Col xs={24} md={16}>
          <Card style={{ marginBottom: 24 }}>
            <GuideStepsForm AIConfig={AIModalList} GitlabConfig={GitlabToken} callback={(key: string) => {
              if (key === 'AIConfig') {
                getAIModalList();
              } else if (key === 'GitlabConfig') {
                getGitlabToken();
              }
            }} />
          </Card>
          <Card title="配置视频教程" style={{ marginBottom: 24 }}>
            <BilibiliVideoCard />
          </Card>
        </Col>}
        {/* 右侧辅助区 */}
        <Col xs={24} md={8}>
          <Card title="产品能力" style={{ marginBottom: 24 }}>
            <ProductAbility />
          </Card>
          <Card title="使用文档与操作说明" style={{ marginBottom: 24 }}>
            <DocLinksCard />
          </Card>
          <Card title="更新日志">
            <UpdateLog />
          </Card>
        </Col>
        </Row>
      </div>
    </Spin>
  );
}

export default Dashboard; 