import { useCallback, useEffect, useState, useContext, memo } from 'react';
import { Card, Row, Col, Typography, Spin } from 'antd';
import { SearchOutlined, BugOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';

import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { queryList as getGitlabTokenList } from '@/pages/aicodecheck/GitlabToken/service';
import { queryList } from '@/pages/aicodecheck/AIModelManager/service';

import MergeStats from './components/MergeStats';
import { GuideStepsForm, UpdateLog, ProductAbility } from './components';
import MergeProblemStats from './components/MergeStats/MergeProblemStats';
import AIDetectionChart from './components/AIDetectionChart';
import RecentRecords from './components/RecentRecords';
import DocLinksCard from './components/DocLinksCard';
import AlertTips from './components/AlertTips';
import BilibiliVideoCard from './components/BilibiliVideoCard';

const { Title } = Typography;

function App() {
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [AIModalList, setAIModalList] = useState<any>([]);
  const [GitlabToken, setGitlabToken] = useState<any>([]);
  const [AIModalLoading, setAIModalLoading] = useState<boolean>(false);
  const [GitlabTokenLoading, setGitlabTokenLoading] = useState<boolean>(false);

  const showConfigAlert = AIModalList.length === 0 || GitlabToken.length === 0;

  const getAIModalList = useCallback(() => {
    setAIModalLoading(true);
    queryList()
      .then((res) => {
        const {
          data: { data },
        } = res;
        setAIModalList(data);
        setAIModalLoading(false);
      })
      .catch((err) => {
        console.log('err', err);
        setAIModalList([]);
        setAIModalLoading(false);
      });
  }, []);

  const getGitlabToken = useCallback(() => {
    setGitlabTokenLoading(true);
    getGitlabTokenList()
      .then((res) => {
        const {
          data: { data },
        } = res;
        setGitlabToken(data);
        setGitlabTokenLoading(false);
      })
      .catch((err) => {
        console.log('err', err);
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
      <div className='layout-main-conent'>
        {!AIModalLoading && !GitlabTokenLoading && (
          <AlertTips AIConfig={AIModalList.length > 0} GitlabConfig={GitlabToken.length > 0} />
        )}

        <Row gutter={32}>
          {!showConfigAlert && (
            <Col xs={24} md={16}>
              <Row gutter={24} style={{ marginBottom: 24 }}>
                <Col span={12}>
                  <Card>
                    <Row gutter={32}>
                      <Col flex='32px'>
                        <SearchOutlined />
                      </Col>
                      <Col flex='auto'>
                        <div style={{ fontSize: 16 }}>{t('page.home.text-detection-times')}</div>
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
                      <Col flex='32px'>
                        <BugOutlined />
                      </Col>
                      <Col flex='auto'>
                        <div>{t('page.home.text-problem-number')}</div>
                        <div>
                          <MergeProblemStats />
                        </div>
                      </Col>
                    </Row>
                  </Card>
                </Col>
              </Row>
              <Card>
                <Title level={5}>{t('page.home.text-ai-detection-effect')}</Title>
                <AIDetectionChart />
              </Card>
              <RecentRecords />
            </Col>
          )}
          {showConfigAlert && (
            <Col xs={24} md={16}>
              <Card>
                <GuideStepsForm
                  AIConfig={AIModalList}
                  GitlabConfig={GitlabToken}
                  callback={(key: string) => {
                    if (key === 'AIConfig') {
                      getAIModalList();
                    } else if (key === 'GitlabConfig') {
                      getGitlabToken();
                    }
                  }}
                />
              </Card>
              <Card title={t('page.home.text-video-tutorial')}>
                <BilibiliVideoCard />
              </Card>
            </Col>
          )}
          <Col xs={24} md={8}>
            <Card title={t('page.home.text-product-ability')}>
              <ProductAbility />
            </Card>
            <Card title={t('page.home.text-doc-links')}>
              <DocLinksCard />
            </Card>
            <Card title={t('page.home.text-update-log')}>
              <UpdateLog />
            </Card>
          </Col>
        </Row>
      </div>
    </Spin>
  );
}

export default memo(observer(App));
