import React, { useState } from 'react';
import { Button, Input } from 'antd';

const { TextArea } = Input;

interface EditableFeedbackProps {
  value?: string;
  onChange?: (val: string) => void;
  placeholder?: string;
}

const EditableFeedback: React.FC<EditableFeedbackProps> = ({ value = '', onChange, placeholder = '请输入您对AI建议的反馈' }) => {
  const [editing, setEditing] = useState(false);
  const [internalValue, setInternalValue] = useState(value);
  const [tempValue, setTempValue] = useState(value);

  // 外部 value 变化时同步 internalValue
  React.useEffect(() => {
    setInternalValue(value);
    setTempValue(value);
  }, [value]);

  const handleEdit = () => {
    setTempValue(internalValue);
    setEditing(true);
  };

  const handleSave = () => {
    setInternalValue(tempValue);
    setEditing(false);
    onChange && onChange(tempValue);
  };

  const handleCancel = () => {
    setTempValue(internalValue);
    setEditing(false);
  };

  if (editing) {
    return (
      <div>
        <TextArea
          value={tempValue}
          onChange={e => setTempValue(e.target.value)}
          placeholder={placeholder}
          autoSize={{ minRows: 6, maxRows: 10 }}
          onPressEnter={e => {
            if (!e.shiftKey && !e.ctrlKey && !e.altKey) {
              e.preventDefault();
              handleSave();
            }
          }}
        />
        <div style={{ marginTop: 8 }}>
          <Button type="primary" size="small" onClick={handleSave} style={{ marginRight: 8 }}>
            保存
          </Button>
          <Button size="small" onClick={handleCancel}>
            取消
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div style={{ minHeight: 48, whiteSpace: 'pre-wrap', color: internalValue ? undefined : '#aaa' }}>
        {internalValue || placeholder}
      </div>
      <Button type="link" size="small" onClick={handleEdit} style={{ padding: 0 }}>
        编辑
      </Button>
    </div>
  );
};

export default EditableFeedback;
