export interface CurrentUser {
  id: number;
  user: string;
  email?: string;
  phone_number?: string;
  job?: string;
  avatar: string;
  roles: string[];
}

export const initialState: CurrentUser = {
  id: 0,
  user: '',
  avatar: '',
  roles: [],
};
