import DB from '../databases';
import { getUnixTimestamp } from '../utils';

class AIMessageService {
    public AIMessage = DB.AIMessage;
    
    public now:number = getUnixTimestamp();

    // 更新human_rating和remark字段
    public async updateHumanRatingAndRemark(Data: any): Promise<any> {
        const { id, human_rating, remark } = Data;
        const res: any = await this.AIMessage.update({ 
            human_rating, remark 
        }, { where: { id } });
        return res;
    }
}

export default AIMessageService;