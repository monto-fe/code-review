import { NextFunction, Request, Response } from 'express';
import aiConfigManager from '../services/aiConfigManager';
import ResponseHandler from '../utils/responseHandler';
import { ResponseMap } from '../utils/const';

const { ParamsError } = ResponseMap;

class AIConfigController {

    // ai config 的增删改查
  public CreateAIConfig = async (req: Request, res: Response, next: NextFunction) => {
    const { name, api_url, api_key, model, is_active } = req.body;
    if (!name || !api_url || !api_key || !model) {
      return ResponseHandler.error(res, { message: 'name, api_url, api_key, model is required' }, ParamsError.message, ParamsError.ret_code);
    }
    const response = await aiConfigManager.addAIConfig({
        name, api_url, api_key, model, is_active
    });
    ResponseHandler.success(res, { data: response }, 'success');
  }

  public GetAIConfig = async (req: Request, res: Response, next: NextFunction) => {
    // const { language } = req.query as any;
    const result = await aiConfigManager.getConfigList();
    ResponseHandler.success(res, { data: result }, 'success');
  }

  public UpdateAIConfig = async (req: Request, res: Response, next: NextFunction) => {
    const { id, name, api_url, api_key, model, is_active } = req.body as any;
    const result = await aiConfigManager.updateAIConfig({ id, name, api_url, api_key, model, is_active });
    ResponseHandler.success(res, { data: result }, 'success');
  }
}

export default AIConfigController;