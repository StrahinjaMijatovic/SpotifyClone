import { Component, OnInit, AfterViewInit, ChangeDetectorRef } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { AuthService } from '../../services/auth.service';
import { RecaptchaService } from '../../services/recaptcha.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './login-component.html',
  styleUrl: './login-component.css',
})
export class LoginComponent implements OnInit, AfterViewInit {
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

  // reCAPTCHA
  recaptchaToken = '';
  recaptchaWidgetId: number | null = null;
  recaptchaLoaded = false;

  constructor(
    private authService: AuthService,
    private router: Router,
    private recaptchaService: RecaptchaService,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    // Load reCAPTCHA script
    this.recaptchaService.load().then(() => {
      this.recaptchaLoaded = true;
    }).catch(err => {
      console.error('Failed to load reCAPTCHA:', err);
    });
  }

  ngAfterViewInit(): void {
    // Render reCAPTCHA after view is initialized
    setTimeout(() => {
      this.renderRecaptcha();
    }, 500);
  }

  renderRecaptcha(): void {
    if (!this.recaptchaLoaded) {
      // Retry after a short delay if not loaded yet
      setTimeout(() => this.renderRecaptcha(), 500);
      return;
    }

    const widgetId = this.recaptchaService.render('recaptcha-login', (token: string) => {
      this.recaptchaToken = token;
    });

    if (widgetId !== null) {
      this.recaptchaWidgetId = widgetId;
    }
  }

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

    // Check reCAPTCHA
    if (!this.recaptchaToken) {
      this.errorMessage = 'Molimo potvrdite da niste robot (reCAPTCHA)';
      return;
    }

    this.loading = true; // Show loading indicator
    this.errorMessage = null;

    this.authService.login({
      username: this.username,
      password: this.password,
      recaptcha_token: this.recaptchaToken
    }).subscribe({
      next: (res) => {
        this.loading = false;
        console.log('Login success, switching to step 2');
        this.tempToken = res.temp_token;
        this.step = 2; // Only switch to OTP step on success
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.loading = false;
        this.step = 1;
        this.errorMessage = err?.error?.error || 'Login failed';
        // Reset reCAPTCHA on error
        this.recaptchaToken = '';
        if (this.recaptchaWidgetId !== null) {
          this.recaptchaService.reset(this.recaptchaWidgetId);
        }
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
