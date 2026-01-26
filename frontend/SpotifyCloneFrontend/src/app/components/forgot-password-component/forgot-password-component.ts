import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-forgot-password',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <div class="d-flex flex-column min-vh-100 justify-content-center align-items-center" style="background: linear-gradient(#1e1e1e 0%, #121212 100%);">
      
      <div class="d-flex align-items-center gap-2 mb-4">
         <div style="width: 40px; height: 40px; background: var(--brand-color); border-radius: 50%;"></div>
         <span class="fw-bold fs-3 text-white" style="letter-spacing: -0.05em;">Spotify Clone</span>
      </div>

      <div class="container" style="max-width: 450px;">
         <div class="card border-0 shadow-lg fade-in" style="background-color: #121212;">
            <div class="card-body p-5 text-center">
               
               <h4 class="fw-bold text-white mb-3">Forgot Password?</h4>
               <p class="text-subdued mb-4 text-start">Enter your email address and we'll send you a link to reset your password.</p>

               <div *ngIf="successMessage" class="alert alert-success mb-4 text-start" style="background: rgba(29, 185, 84, 0.1); border: 1px solid var(--text-positive); color: var(--text-positive);">
                 {{ successMessage }}
               </div>
               <div *ngIf="errorMessage" class="alert alert-danger mb-4 text-start" style="background: rgba(229, 9, 20, 0.1); border: 1px solid var(--text-negative); color: var(--text-negative);">
                 {{ errorMessage }}
               </div>

               <form (ngSubmit)="onSubmit()" *ngIf="!successMessage">
                 <div class="mb-4 text-start">
                   <label class="form-label text-white fw-bold small">Email address</label>
                   <input type="email" class="form-input" [(ngModel)]="email" name="email" required placeholder="name@domain.com">
                 </div>
                 
                 <button class="btn-primary w-100 py-3 rounded-pill fw-bold text-uppercase" 
                         type="submit" 
                         [disabled]="loading"
                         style="letter-spacing: 0.1em; font-size: 0.85rem;">
                   <span *ngIf="loading" class="spinner-border spinner-border-sm me-2"></span>
                   Send Link
                 </button>
               </form>

               <div class="mt-4 pt-3 border-top border-secondary" style="border-color: #333 !important;">
                  <a routerLink="/login" class="text-subdued small text-decoration-none fw-bold text-uppercase" style="letter-spacing: 0.1em;">Back to Login</a>
               </div>
            </div>
         </div>
      </div>
    </div>
  `,
  styles: []
})
export class ForgotPasswordComponent {
  email = '';
  loading = false;
  successMessage: string | null = null;
  errorMessage: string | null = null;

  constructor(private auth: AuthService) { }

  onSubmit(): void {
    if (!this.email) return;

    this.loading = true;
    this.errorMessage = null;

    this.auth.requestPasswordReset(this.email.trim()).subscribe({
      next: (res) => {
        // Backend always returns 200 OK even if email not found (security)
        this.successMessage = res.message || 'Ako email postoji, link je poslat.';
        this.loading = false;
      },
      error: () => {
        this.errorMessage = 'Greška pri slanju zahteva.';
        this.loading = false;
      }
    });
  }
}
