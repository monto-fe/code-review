import React, { useState } from 'react';
import { Flex, Rate as AntRate, message } from 'antd';
import { updateRating } from './service';

const desc = ['terrible', 'bad', 'normal', 'good', 'wonderful'];

interface RateProps {
  id: number;
  initialValue?: number;
}

const Rate: React.FC<RateProps> = ({ id, initialValue = 0 }) => {
  const [value, setValue] = useState(initialValue);
  const [loading, setLoading] = useState(false);

  const handleChange = async (newValue: number) => {
    setLoading(true);
    try {
      await updateRating(id, newValue);
      setValue(newValue);
      message.success('评分更新成功');
    } catch (error) {
      message.error('评分更新失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Flex gap="middle" vertical>
      <AntRate 
        tooltips={desc} 
        onChange={handleChange} 
        value={value} 
        disabled={loading}
      />
      {value ? <span>{desc[value - 1]}</span> : null}
    </Flex>
  );
};

export default Rate;