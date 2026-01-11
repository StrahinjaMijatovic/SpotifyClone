export type LoginRequest = {
  username: string;
  password: string;
};

export type LoginUser = {
  id: string;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
};

export type LoginResponse = {
  token: string;
  user: LoginUser;
};
