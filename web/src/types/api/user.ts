// web/src/types/api/user.ts

/**
 * @description 用户/组织类型
 * @description 0: 平台
 * @description 1: 学校
 * @description 2: 供应商
 * @description 3: 食堂
 * @description 4: 商户
 */
export const UserType = {
  Platform: 0,
  School: 1,
  Supplier: 2,
  Canteen: 3,
  Merchant: 4,
} as const;

export type UserType = (typeof UserType)[keyof typeof UserType];

export const userTypeMap = {
  [UserType.Platform]: '平台',
  [UserType.School]: '学校',
  [UserType.Supplier]: '供应商',
  [UserType.Canteen]: '食堂',
  [UserType.Merchant]: '商户',
};
