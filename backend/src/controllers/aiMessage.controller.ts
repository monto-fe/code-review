import { NextFunction, Request, Response } from 'express';
import AICheckService from '../services/aiCheck.service';
import GitlabService from '../services/gitlab.service';
import AIRuleService from '../services/aiRule.service';
import ResponseHandler from '../utils/responseHandler';
import { ResponseMap } from '../utils/const';

const { ParamsError } = ResponseMap;

class AIMessageController {
  public AICheckService = new AICheckService();
  public GitlabService = new GitlabService();
  public AIRuleService = new AIRuleService();

  public CreateCommonRule = async (req: Request, res: Response, next: NextFunction) => {
    const { name, language, rule, description } = req.body as any;
    if (!name || !language || !rule) {
      return ResponseHandler.error(res, { message: 'name, language, rule is required' }, ParamsError.message, ParamsError.ret_code);
    }
    const { rows, count } = await this.AIRuleService.createCommonRule({
        name, language, rule, description
    });
    ResponseHandler.success(res, { data: rows, total: count }, 'success');
  }

  public GetCommonRule = async (req: Request, res: Response, next: NextFunction) => {
    const { language } = req.query as any;
    const result = await this.AIRuleService.getCommonRule({language});
    ResponseHandler.success(res, { data: result }, 'success');
  }

  public UpdateCommonRule = async (req: Request, res: Response, next: NextFunction) => {
    const { id, name, language, rule, description } = req.body as any;
    const result = await this.AIRuleService.updateCommonRule({
        id, name, language, rule, description
    });
    ResponseHandler.success(res, { data: result }, 'success');
  }

  public GetCustomRule = async (req: Request, res: Response, next: NextFunction) => {
    const { project_id } = req.query as any;
    const { rows, count} = await this.AIRuleService.getCustomRuleByProjectId({ project_id });
    ResponseHandler.success(res, { data: rows || [], total: count }, 'success');
  }

  public CreateCustomRule = async (req: Request, res: Response, next: NextFunction) => {
    const { project_name, project_id, rule, status, operator } = req.body as any;
    const result = await this.AIRuleService.createCustomRule({ 
      project_name, project_id, rule, status, operator 
    });
    ResponseHandler.success(res, { data: result }, 'success');
  }

  public UpdateCustomRule = async (req: Request, res: Response, next: NextFunction) => {
    const { id, project_name, project_id, rule, status, operator } = req.body as any;
    const result = await this.AIRuleService.updateCustomRule({
      id, project_name, project_id, rule, status, operator 
    });
    ResponseHandler.success(res, { data: result }, 'success');
  }

  // 更新human_rating和remark字段
  public UpdateHumanRatingAndRemark = async (req: Request, res: Response, next: NextFunction) => {
    const { id, human_rating, remark } = req.body as any;
    if (!id || !human_rating) {
      return ResponseHandler.error(res, { message: 'id, human_rating is required' }, ParamsError.message, ParamsError.ret_code);
    }
    const result = await this.AIRuleService.updateHumanRatingAndRemark({
      id, human_rating, remark
    });
    ResponseHandler.success(res, { data: result }, 'success');
  }

}

export default AIMessageController;