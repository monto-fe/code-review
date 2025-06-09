import React, { useState } from 'react';
import { Drawer, Collapse, Button } from 'antd';

const { Panel } = Collapse;

interface PromptDrawerProps {
  visible: boolean;
  onClose: () => void;
  prompts: { title: string; content: string }[];
  onUse: (content: string) => void;
}

const PromptDrawer: React.FC<PromptDrawerProps> = ({ visible, onClose, prompts, onUse }) => {
  // 只展示前3个
  const displayPrompts = prompts.slice(0, 3);

  return (
    <Drawer
      title="提示词列表"
      placement="right"
      width={800}
      open={visible}
      onClose={onClose}
      destroyOnClose
    >
      <Collapse accordion>
        {displayPrompts.map((item, idx) => (
          <Panel header={item.title || `提示词${idx + 1}`} key={idx}>
            <div style={{ whiteSpace: 'pre-wrap', marginBottom: 16 }}>
              {item.content}
            </div>
            <Button type="primary" size="small" onClick={() => onUse(item.content)}>
              使用
            </Button>
          </Panel>
        ))}
      </Collapse>
    </Drawer>
  );
};

export default PromptDrawer;
