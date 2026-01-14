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
    <div class="min-vh-100 d-flex align-items-center bg-black">
      <div class="container">
        <div class="row justify-content-center">
          <div class="col-12 col-md-6 col-lg-4">
            <div class="card bg-dark border-secondary shadow-lg">
              <div class="card-body p-4 p-md-5">
                <div class="text-center mb-4">
                  <h3 class="fw-bold text-light">Zaboravljena lozinka?</h3>
                  <p class="text-secondary small">Unesi email i poslaćemo ti link za resetovanje.</p>
                </div>

                <div *ngIf="successMessage" class="alert alert-success">
                  {{ successMessage }}
                </div>
                <div *ngIf="errorMessage" class="alert alert-danger">
                  {{ errorMessage }}
                </div>

                <form (ngSubmit)="onSubmit()" *ngIf="!successMessage">
                  <div class="mb-3">
                    <label class="form-label text-secondary">Email</label>
                    <input type="email" class="form-control" [(ngModel)]="email" name="email" required placeholder="tvoj@email.com">
                  </div>
                  <button class="btn btn-success w-100 py-2" type="submit" [disabled]="loading">
                    <span *ngIf="loading" class="spinner-border spinner-border-sm me-2"></span>
                    Pošalji link
                  </button>
                </form>

                 <div class="text-center mt-4">
                  <a routerLink="/login" class="text-decoration-none text-secondary">Nazad na login</a>
                </div>
              </div>
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
