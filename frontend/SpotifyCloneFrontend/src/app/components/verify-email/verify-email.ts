import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-verify-email',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="min-vh-100 d-flex align-items-center justify-content-center bg-black text-light">
      <div class="text-center">
        <div *ngIf="loading" class="mb-3">
          <div class="spinner-border text-success" role="status" style="width: 3rem; height: 3rem;"></div>
          <p class="mt-3">Verifikacija email-a...</p>
        </div>

        <div *ngIf="successMessage" class="alert alert-success">
          {{ successMessage }}
          <div class="mt-3">
             <button class="btn btn-success" (click)="goToLogin()">Idi na Login</button>
          </div>
        </div>

        <div *ngIf="errorMessage" class="alert alert-danger">
          {{ errorMessage }}
          <div class="mt-3">
             <button class="btn btn-outline-light" (click)="goToLogin()">Nazad na Login</button>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: []
})
export class VerifyEmailComponent implements OnInit {
  loading = true;
  successMessage: string | null = null;
  errorMessage: string | null = null;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private auth: AuthService
  ) { }

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      const token = params['token'];
      if (token) {
        this.verifyToken(token);
      } else {
        this.loading = false;
        this.errorMessage = 'Token nije pronađen.';
      }
    });
  }

  verifyToken(token: string): void {
    this.auth.verifyEmail(token).subscribe({
      next: (res) => {
        this.loading = false;
        this.successMessage = res.message || 'Email uspešno verifikovan! Sada se možeš ulogovati.';
      },
      error: (err) => {
        this.loading = false;
        this.errorMessage = err?.error?.error || 'Link je istekao ili je neispravan.';
      }
    });
  }

  goToLogin(): void {
    this.router.navigate(['/login']);
  }
}
