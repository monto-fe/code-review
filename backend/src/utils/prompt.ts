export const systemPrompt = `
    你是一个代码审核专家，你需要根据用户的要求，对代码进行评审，然后输出简洁的符合标准的评论。
    。 `
export const generatePrompt = ({
    rule,
    title,
    description,
    author_name,
    web_url,
    diff,
}:{
    rule: string,
    title: string,
    description: string,
    author_name: string,
    web_url: string,
    diff: any[],
}):string => {
    const prompt = `
    请检查以下代码差异（diff），确保其符合以下原则：
    ${rule}
    输出要求：
    仅输出以下两部分内容，如无内容可省略：
        1. 不符合检查原则的地方（使用 Markdown 表格格式）：
          - 不符合的代码行号 | 违反的原则 | 修改建议
        2. 疑似 Bug 的地方（基于描述信息和逻辑判断分析，如无可省略）。
    代码信息：
        合并请求标题：${title}
        描述: ${description || '无描述'}
        提交人: ${author_name}
        查看合并请求: ${web_url}
        代码差异：${diff.map((change:any) => `文件: ${change.new_path}\n差异:\n${change.diff}`).join('\n')}
        请按照上述要求进行检查并输出结果。
    `
    return prompt;
}