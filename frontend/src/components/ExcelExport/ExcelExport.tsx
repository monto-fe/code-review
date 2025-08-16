import React, { useRef, useState } from 'react';
import { Button, message, Modal, Form, Select, Checkbox, Space, Input, Tooltip } from 'antd';
import { DownloadOutlined, SettingOutlined } from '@ant-design/icons';
import ExcelJS from 'exceljs';
import type { ColumnsType } from 'antd/lib/table';
import { processTableData } from './utils';
import type { ExcelExportConfig, ExcelExportProps } from './types';

const { Option } = Select;



function ExcelExport<T = any>({
  data,
  columns,
  config = {},
  buttonText = '导出Excel',
  buttonType = 'default',
  buttonSize = 'middle',
  showSettings = false,
  disabled = false,
  loading = false,
  onExport,
  onBeforeExport
}: ExcelExportProps<T>) {
  const [settingsVisible, setSettingsVisible] = useState(false);
  const [exportLoading, setExportLoading] = useState(false);
  const [form] = Form.useForm();
  
  const defaultConfig: ExcelExportConfig = {
    filename: 'export_data',
    sheetName: 'Sheet1',
    includeHeaders: true,
    autoWidth: true,
    selectedColumns: columns.map(col => col.key as string).filter(Boolean),
    ...config
  };

  const handleExport = async (exportConfig?: ExcelExportConfig) => {
    try {
      setExportLoading(true);
      
      const finalConfig = { ...defaultConfig, ...exportConfig };
      
      // 处理数据
      let exportData = data;
      if (onBeforeExport) {
        exportData = await onBeforeExport(data, finalConfig);
      }
      
      // 使用工具函数处理数据
      const processedData = processTableData(
        exportData, 
        columns, 
        finalConfig.selectedColumns || []
      );
      console.log("processedData", processedData);
      
      // 转换为ExcelJS需要的格式
      const excelData = processedData.map(item => Object.values(item));
      
      // 使用ExcelJS导出Excel
      const workbook = new ExcelJS.Workbook();
      const worksheet = workbook.addWorksheet(finalConfig.sheetName || 'Sheet1');
      
      // 获取表头
      const headers = finalConfig.includeHeaders && processedData.length > 0 
        ? Object.keys(processedData[0]) 
        : [];
      
      console.log('headers:', headers);
      
      // 直接添加表头行
      if (finalConfig.includeHeaders && headers.length > 0) {
        worksheet.addRow(headers);
      }
      
      // 直接添加数据行（数组格式）
      excelData.forEach((row, index) => {
        console.log(`row ${index}:`, row);
        worksheet.addRow(row);
      });
      
      // 自动调整列宽
      if (finalConfig.autoWidth && headers.length > 0) {
        worksheet.columns = headers.map((header, index) => ({
          header: header,
          key: header,
          width: 20
        }));
      }
      
      // 生成并下载文件
      const buffer = await workbook.xlsx.writeBuffer();
      
      const blob = new Blob([buffer], { 
        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' 
      });
      
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${finalConfig.filename}.xlsx`;
      link.click();
      URL.revokeObjectURL(url);
      
      message.success('导出成功');
      onExport?.(finalConfig);
      
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败，请重试');
    } finally {
      setExportLoading(false);
    }
  };

  const handleSettingsExport = async () => {
    try {
      const values = await form.validateFields();
      const exportConfig: ExcelExportConfig = {
        filename: values.filename,
        sheetName: values.sheetName,
        includeHeaders: values.includeHeaders,
        autoWidth: values.autoWidth,
        selectedColumns: values.selectedColumns
      };
      
      await handleExport(exportConfig);
      setSettingsVisible(false);
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  const openSettings = () => {
    form.setFieldsValue({
      filename: defaultConfig.filename,
      sheetName: defaultConfig.sheetName,
      includeHeaders: defaultConfig.includeHeaders,
      autoWidth: defaultConfig.autoWidth,
      selectedColumns: defaultConfig.selectedColumns
    });
    setSettingsVisible(true);
  };

  return (
    <>
      <Tooltip title={buttonText}>
        <DownloadOutlined 
          onClick={() => handleExport()}
          style={{ 
            cursor: disabled || exportLoading ? 'not-allowed' : 'pointer',
            color: disabled || exportLoading ? '#d9d9d9' : undefined,
            fontSize: '16px'
          }}
        />
      </Tooltip>
      
      {showSettings && (
        <Button
          type="default"
          size={buttonSize}
          icon={<SettingOutlined />}
          onClick={openSettings}
          disabled={disabled}
        >
          导出设置
        </Button>
      )}

      <Modal
        title="导出设置"
        open={settingsVisible}
        onOk={handleSettingsExport}
        onCancel={() => setSettingsVisible(false)}
        confirmLoading={exportLoading}
        width={500}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={defaultConfig}
        >
          <Form.Item
            label="文件名"
            name="filename"
            rules={[{ required: true, message: '请输入文件名' }]}
          >
            <Select
              placeholder="请选择或输入文件名"
              allowClear
              showSearch
              mode="tags"
              options={[
                { label: '评论列表', value: 'comment_list' },
                { label: '检测结果', value: 'detection_results' },
                { label: '数据导出', value: 'data_export' }
              ]}
            />
          </Form.Item>
          
          <Form.Item
            label="工作表名称"
            name="sheetName"
            rules={[{ required: true, message: '请输入工作表名称' }]}
          >
            <Select
              placeholder="请选择或输入工作表名称"
              allowClear
              showSearch
              mode="tags"
              options={[
                { label: 'Sheet1', value: 'Sheet1' },
                { label: '数据', value: '数据' },
                { label: '结果', value: '结果' }
              ]}
            />
          </Form.Item>
          
          <Form.Item
            label="导出列"
            name="selectedColumns"
            rules={[{ required: true, message: '请选择要导出的列' }]}
          >
            <Select
              mode="multiple"
              placeholder="请选择要导出的列"
              options={columns
                .filter(col => col.key && col.title)
                .map(col => ({
                  label: col.title as string,
                  value: col.key as string
                }))}
            />
          </Form.Item>
          
          <Form.Item name="includeHeaders" valuePropName="checked">
            <Checkbox>包含表头</Checkbox>
          </Form.Item>
          
          <Form.Item name="autoWidth" valuePropName="checked">
            <Checkbox>自动调整列宽</Checkbox>
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
}

export default ExcelExport; 