declare module 'exceljs' {
  export default class ExcelJS {
    static Workbook: new () => Workbook;
  }

  export class Workbook {
    addWorksheet(name: string): Worksheet;
    xlsx: {
      writeBuffer(): Promise<Buffer>;
    };
  }

  export class Worksheet {
    addRow(values: any[] | any): Row;
    columns: Column[];
    rowCount?: number;
  }

  export class Row {
    values: any[];
  }

  export class Column {
    values?: any[];
    width?: number;
  }
} 