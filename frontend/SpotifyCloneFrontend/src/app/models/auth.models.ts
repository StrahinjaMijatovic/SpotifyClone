export type LoginRequest = {
  username: string;
  password: string;
};

export type LoginInitiateResponse = {
  message: string;
  temp_token: string;
};

export type VerifyOTPRequest = {
  temp_token: string;
  otp_code: string;
};

export type LoginUser = {
  id: string;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
};

export type VerifyOTPResponse = {
  token: string;
  user: LoginUser;
};

export type UpdateProfileRequest = {
  first_name: string;
  last_name: string;
};

export interface UpdateProfileResponse {
  message: string;
  user: LoginUser;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface ResetPasswordRequest {
  email: string;
}

export interface ResetPasswordConfirmRequest {
  token: string;
  new_password: string;
}

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

