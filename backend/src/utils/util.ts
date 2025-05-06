/* eslint-disable */
const jwt = require('jsonwebtoken')
const fs = require('fs');
const dayjs = require('dayjs');
const { Buffer } = require('buffer');
/**
 * @method isEmpty
 * @param {String | Number | Object} value
 * @returns {Boolean} true & false
 * @description this value is Empty Check
 */
export const isEmpty = (value: string | number | object): boolean => {
  if (value === null) {
    return true;
  } else if (typeof value !== 'number' && value === '') {
    return true;
  } else if (typeof value === 'undefined' || value === undefined) {
    return true;
  } else if (value !== null && typeof value === 'object' && !Object.keys(value).length) {
    return true;
  } else {
    return false;
  }
};

export const pageSize = {
  offset: 0,
  limit: 10
}

export const TokenSecretKey = "acl"
export const TokenExpired = '8h'
export const RefreshTokenExpired = '3 days'
export const RetCodeMap = {
  Success: 0,
  NotLogIn: 10000
}

export const generateSignToken = (
  option:any, 
  secretKey: string, 
  expiresIn?: string | number | undefined
) => {
  return jwt.sign(
    option, 
    secretKey,
    {
      expiresIn: expiresIn || TokenExpired
    }
  )
}

export const checkSignToken = async (jwtToken: string, secretKey: string) => {
  return new Promise((resolve, reject) => {
    jwt.verify(
      jwtToken, 
      secretKey, 
      function(err:Error, value:any) {
        if(err){
          reject(err.message)
        }else{
          resolve(value)
        }
      }
    )
  })
}

export const IsTokenExpired = (token_expired_at: number) => {
  const now = dayjs().unix();
  // const hour5 = 5 * 60 * 60
  if(token_expired_at - now < 0){
    return true
  }
  return false
}

export function splitIntegerToObject(n: number, j: string[]): Record<string, number> {
  const length = j.length;
  const result: Record<string, number> = {};

  if (length === 0) {
    return result; // 如果 j 是空数组，直接返回空对象
  }

  const quotient = Math.floor(n / length);
  const remainder = n % length;

  for (const element of j) {
    result[element] = quotient;
  }

  // 分配余数
  for (let i = 0; i < remainder; i++) {
    result[j[i]]++;
  }

  return result;
}

export function base64ToBlob(base64String:string, mimeType:string) {
  const binaryData = Buffer.from(base64String, 'base64');
  return { data: binaryData, mimeType };
}

export function generateAudioFileKey(userId: number, questionId: number):string{
  // 格式化日期
  const CurrentDay = dayjs();
  const formattedDate = CurrentDay.format('YYYY_MM_DD');
  const timestamp = CurrentDay.unix();
  // 生成字符串
  const resultString = `${userId}_${timestamp}_${formattedDate}_${questionId}.mp3`;
  return resultString
}