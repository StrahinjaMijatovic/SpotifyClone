import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import {
  LoginRequest,
  LoginInitiateResponse,
  VerifyOTPRequest,
  VerifyOTPResponse,
  RegisterRequest,
  RegisterResponse,
  UpdateProfileRequest,
  UpdateProfileResponse,
} from '../models/auth.models';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly apiBase = '/api/v1';

  constructor(private http: HttpClient) { }

  // Step 1: Login (returns temp_token)
  login(payload: LoginRequest): Observable<LoginInitiateResponse> {
    return this.http.post<LoginInitiateResponse>(`${this.apiBase}/login`, payload);
  }

  // Step 2: Verify OTP (returns JWT)
  verifyOTP(payload: VerifyOTPRequest): Observable<VerifyOTPResponse> {
    return this.http.post<VerifyOTPResponse>(`${this.apiBase}/verify-otp`, payload).pipe(
      tap((res) => {
        localStorage.setItem('token', res.token);
        // Optionally store user info if needed
      })
    );
  }

  logout(): Observable<any> {
    // Call backend to invalidate token
    return this.http.post(`${this.apiBase}/logout`, {}).pipe(
      tap(() => {
        localStorage.removeItem('token');
      })
    );
  }

  getToken(): string | null {
    return localStorage.getItem('token');
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }

  register(payload: RegisterRequest): Observable<RegisterResponse> {
    return this.http.post<RegisterResponse>(`${this.apiBase}/register`, payload);
  }

  getProfile(): Observable<any> {
    return this.http.get(`${this.apiBase}/profile`);
  }

  updateProfile(payload: UpdateProfileRequest): Observable<UpdateProfileResponse> {
    return this.http.put<UpdateProfileResponse>(`${this.apiBase}/profile`, payload);
  }

  deleteAccount(): Observable<any> {
    return this.http.delete(`${this.apiBase}/profile`);
  }
}
