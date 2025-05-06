// utils/ResponseHandler.ts
import { Response } from 'express';

class ResponseHandler {
  static success(res: Response, data: any = {}, message: string = 'Success', statusCode: number = 0) {
    return res.status(200).json({
      ret_code: statusCode,
      message,
      data
    });
  }

  static error(res: Response, error: Record<string, any>, message?:string, statusCode: number = 0) {
    return res.status(200).json({
      ret_code: error?.ret_code || statusCode || 10000,
      message: message || error.message || 'An error occurred',
      error: error.details || {}
    });
  }
}

export default ResponseHandler;