import { NextFunction, Request, Response } from 'express';
import { ResponseMap, HttpCodeSuccess } from '../utils/const'
import { checkSignToken, TokenSecretKey, IsTokenExpired } from '../utils/util';

const { TokenExpired } = ResponseMap;
const whitePathList: string[] = [
    '/v1/login',
    '/v1/register',
    '/v1/webhook/merge'
]
// 自定义中间件，用于检查 JWT
const authenticateJwt = async (req: Request, res: Response, next: NextFunction) => {
    if(!whitePathList.includes(req.path)){
        const jwtToken = req.header('jwt_token') as string;
        if (!jwtToken) {
          res.status(HttpCodeSuccess).json({ ...TokenExpired});
          return;
        }
        try {
          const tokenInfo:any = await checkSignToken(jwtToken, TokenSecretKey)
          const { exp } = tokenInfo;
          // TODO: 这里只对Token的有效期进行了验证，更安全的做法是，读取数据库对比
          if(IsTokenExpired(exp)){
            res.status(HttpCodeSuccess).json({ ...TokenExpired});
            return
          }
          // TODO: 如果有效期的时间少于5小时，则更新有效期+5h
          req.headers["remoteUser"] = tokenInfo.user
          req.headers["userId"] = tokenInfo.id
        }catch(err:any){
          res.status(HttpCodeSuccess).json({ ...TokenExpired});
          return
        }
    }
    next();
};

export { authenticateJwt };