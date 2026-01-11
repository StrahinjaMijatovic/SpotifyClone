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

export type RegisterRequest = {
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  password: string;
  password_confirm: string;
};

export type RegisterResponse = {
  message: string;
  user_id: string;
};
