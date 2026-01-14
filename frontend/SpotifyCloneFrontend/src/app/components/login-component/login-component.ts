import { Component } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './login-component.html',
  styleUrl: './login-component.css',
})
export class LoginComponent {
  username = '';
  password = '';
  otpCode = '';

  loading = false;
  errorMessage: string | null = null;
  showPassword = false;

  year = new Date().getFullYear();

  step = 1;
  tempToken = '';

  constructor(private auth: AuthService, private router: Router) { }

  onSubmit(): void {
    this.errorMessage = null;

    if (this.step === 1) {
      this.handleLogin();
    } else {
      this.handleOTP();
    }
  }

  handleLogin(): void {
    const u = this.username.trim();
    const p = this.password;

    if (!u || !p) {
      this.errorMessage = 'Unesi username i password.';
      return;
    }

    this.loading = true;

    this.auth.login({ username: u, password: p }).subscribe({
      next: (res) => {
        this.tempToken = res.temp_token;
        this.step = 2; // Move to OTP step
        this.errorMessage = null;
      },
      error: () => {
        this.errorMessage = 'Login nije uspeo. Proveri kredencijale.';
        this.loading = false;
      },
      complete: () => (this.loading = false),
    });
  }

  handleOTP(): void {
    const code = this.otpCode.trim();
    if (!code || code.length !== 6) {
      this.errorMessage = 'Unesi 6-cifreni OTP kod.';
      return;
    }

    this.loading = true;

    this.auth.verifyOTP({ temp_token: this.tempToken, otp_code: code }).subscribe({
      next: () => this.router.navigate(['/home']),
      error: () => {
        this.errorMessage = 'Pogrešan ili istekao OTP kod.';
        this.loading = false;
      },
      complete: () => (this.loading = false),
    });
  }
}
