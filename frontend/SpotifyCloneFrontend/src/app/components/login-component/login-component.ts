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
  step = 1; // 1: Creds, 2: OTP, 3: Magic Link Request

  username = '';
  password = '';
  otpCode = '';
  magicEmail = '';

  showPassword = false;
  loading = false;
  errorMessage: string | null = null;

  tempToken = ''; // from step 1

  year = new Date().getFullYear();

  constructor(private authService: AuthService, private router: Router) { }

  toggleMagicLink(): void {
    this.step = 3;
    this.errorMessage = null;
  }

  onSubmit(): void {
    this.errorMessage = null;

    if (this.step === 1) {
      this.handleLogin();
    } else if (this.step === 2) {
      this.handleOTP();
    } else if (this.step === 3) {
      this.handleMagicLinkRequest();
    }
  }

  handleMagicLinkRequest(): void {
    if (!this.magicEmail) return;
    this.loading = true;

    this.authService.requestMagicLink(this.magicEmail).subscribe({
      next: (res) => {
        this.loading = false;
        alert(res.message || 'Ako nalog postoji, link je poslat!');
        this.step = 1; // Return to login
      },
      error: () => {
        this.loading = false;
        this.errorMessage = 'Greška pri slanju zahteva.';
      }
    });
  }

  handleLogin(): void {
    if (!this.username || !this.password) return;

    this.loading = true;
    this.authService.login({ username: this.username, password: this.password }).subscribe({
      next: (res) => {
        this.tempToken = res.temp_token;
        this.step = 2; // Move to OTP step
        this.loading = false;
      },
      error: (err) => {
        this.errorMessage = err?.error?.error || 'Login failed';
        this.loading = false;
      }
    });
  }

  handleOTP(): void {
    if (!this.otpCode) return;
    this.loading = true;

    this.authService.verifyOTP({ temp_token: this.tempToken, otp_code: this.otpCode }).subscribe({
      next: (res) => {
        this.router.navigate(['/home']);
      },
      error: (err) => {
        this.errorMessage = err?.error?.error || 'Invalid OTP';
        this.loading = false;
      }
    });
  }
}
