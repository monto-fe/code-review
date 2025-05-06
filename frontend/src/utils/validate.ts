/**
 * 判断是否是外链
 * @param {string} path
 * @returns {Boolean}
 * @author duheng1992
 */
export const isExternal = (path: string): boolean => /^(https?:|mailto:|tel:)/.test(path);
