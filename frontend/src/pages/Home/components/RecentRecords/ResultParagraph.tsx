import React, { useState } from 'react';
import { Typography } from 'antd';
import ReactMarkdown from 'react-markdown';

const { Paragraph } = Typography;

const ResultParagraph = ({ text }: { text: string }) => {
  const [expanded, setExpanded] = useState(false);
  return (
    <Paragraph
      ellipsis={{
        rows: 6,
        expandable: 'collapsible',
        expanded,
        onExpand: (_, info) => setExpanded(info.expanded),
      }}
    >
      <ReactMarkdown>{text}</ReactMarkdown>
    </Paragraph>
  );
};

export default ResultParagraph;
